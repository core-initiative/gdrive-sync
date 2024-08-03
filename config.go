package main

import (
    "os"
    "io"
    "github.com/qiangxue/go-env"
    "gopkg.in/yaml.v2"
    "drive-sync/log"
)

type Config struct {
    Folder        string `yaml:"folder" env:"FOLDER"`
    Schedule      string `yaml:"schedule" env:"SCHEDULE"`
    TargetFolderID string `yaml:"targetFolderID" env:"TARGET_FOLDER_ID"`
}

func LoadConfig(filename string, logger log.Logger) (*Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    // load from environment variables prefixed with "APP_"
	if err = env.New("APP_", logger.Infof).Load(&config); err != nil {
		return nil, err
	}

    return &config, nil
}
