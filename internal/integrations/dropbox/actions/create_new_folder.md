## Create Folder Action

This action creates a new folder at the specified path in Dropbox.

### Input Parameters

- **Path** (required): The path where the new folder should be created (e.g., /Homework/math).
- **Auto Rename** (optional): If there's a conflict, have the Dropbox server try to autorename the folder to avoid conflict.

### Output

The action returns the metadata of the newly created folder from the Dropbox API.

### Example Usage

Use this action when you need to create a new organizational structure within your Dropbox storage.