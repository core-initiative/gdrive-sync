package main

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/drive/v3"
    "google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {
    tokFile := "token.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(tokFile, tok)
    }
    return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

func saveToken(path string, token *oauth2.Token) {
    f, err := os.Create(path)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        log.Fatalf("Unable to read authorization code: %v", err)
    }

    tok, err := config.Exchange(context.TODO(), authCode)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web: %v", err)
    }
    return tok
}

func SyncFolder(folderPath string, targetFolderID string) {
    b, err := os.ReadFile("credentials.json")
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

    config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
    if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
    }

    client := getClient(config)
    srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
    if err != nil {
        log.Fatalf("Unable to retrieve Drive client: %v", err)
    }

    driveFiles, err := listDriveFiles(srv, targetFolderID)
    if err != nil {
        log.Fatalf("Unable to list files on Drive: %v", err)
    }

    err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() {
            localFileHash, err := calculateFileHash(path)
            if err != nil {
                log.Printf("Failed to calculate hash for file %s: %v", path, err)
                return nil
            }

            fileID, exists := fileExistsOnDrive(info.Name(), driveFiles)
            if exists {
                driveFileHash, err := getDriveFileHash(srv, fileID)
                if err != nil {
                    log.Printf("Failed to get hash for drive file %s: %v", info.Name(), err)
                    return nil
                }

                if localFileHash != driveFileHash {
                    if err := updateFile(srv, fileID, path); err != nil {
                        log.Printf("Failed to update file %s: %v", path, err)
                    }
                } else {
                    fmt.Printf("File is up to date on Drive: %s\\n", info.Name())
                }
            } else {
                if err := uploadFile(srv, path, targetFolderID); err != nil {
                    log.Printf("Failed to upload file %s: %v", path, err)
                }
            }
        }
        return nil
    })

    if err != nil {
        log.Fatalf("Unable to walk folder: %v", err)
    }
}

func listDriveFiles(srv *drive.Service, folderID string) (map[string]string, error) {
    files := make(map[string]string)
    pageToken := ""
    for {
        request := srv.Files.List().Q("'" + folderID + "' in parents").Fields("nextPageToken, files(id, name)").PageToken(pageToken)
        result, err := request.Do()
        if err != nil {
            return nil, err
        }
        for _, file := range result.Files {
            files[file.Name] = file.Id
        }
        pageToken = result.NextPageToken
        if pageToken == "" {
            break
        }
    }
    return files, nil
}

func fileExistsOnDrive(fileName string, driveFiles map[string]string) (string, bool) {
    fileID, exists := driveFiles[fileName]
    return fileID, exists
}

func getDriveFileHash(srv *drive.Service, fileID string) (string, error) {
    file, err := srv.Files.Get(fileID).Fields("md5Checksum").Do()
    if err != nil {
        return "", err
    }
    return file.Md5Checksum, nil
}

func calculateFileHash(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hasher := sha256.New()
    if _, err := io.Copy(hasher, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(hasher.Sum(nil)), nil
}

func uploadFile(srv *drive.Service, filePath, folderID string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        return err
    }

    driveFile := &drive.File{
        Name: fileInfo.Name(),
        Parents: []string{folderID},
    }

    _, err = srv.Files.Create(driveFile).Media(file).Do()
    if err != nil {
        return err
    }

    fmt.Printf("Uploaded file: %s\n", filePath)
    return nil
}

func updateFile(srv *drive.Service, fileID, filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = srv.Files.Update(fileID, nil).Media(file).Do()
    if err != nil {
        return err
    }

    fmt.Printf("Updated file: %s\n", filePath)
    return nil
}