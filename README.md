## Simple tools for syncing a specific folder to gdrive folder
You need to provide a credentials.json for Google Drive API. This code was generated using ChatGPT 4o, with some revision and feedback from me @purwaren.
### Steps to Get credentials.json:
1. Go to the Google Cloud Console.
2. Create a new project or select an existing project.
3. Enable the Google Drive API for the project.
4. Create OAuth 2.0 credentials (Client ID and Secret).
5. Download the credentials.json file.
6. Placing the credentials.json file:
7. Place the credentials.json file in the same directory as your Go application or executable.
8. Ensure the file is named credentials.json.

### How to get folderId on google drive
1. Go to Google Drive.
2. Navigate to the Target Folder:
3. Locate and open the folder you want to use as the target for your sync.
4. Get the Folder ID:
Look at the URL in your browser’s address bar. It will look something like this: https://drive.google.com/drive/u/0/folders/XXXXXXXXXXXXXXXXXXXXXXXXXXX.
The XXXXXXXXXXXXXXXXXXXXXXXXXXX part is the folder ID. Copy this ID.
