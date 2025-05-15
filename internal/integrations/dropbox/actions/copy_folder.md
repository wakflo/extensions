## Copy Folder Action

This action copies a folder and its contents from one location to another within Dropbox.

### Input Parameters

- **From Path** (required): The source path of the folder (e.g. /folder1/sourceFolder).
- **To Path** (required): The destination path for the copied folder (e.g. /folder2/destinationFolder).
- **Auto Rename** (optional): If there's a conflict, have the Dropbox server try to autorename the folder to avoid conflict.
- **Allow Ownership Transfer** (optional): Allows copy by owner even if it would result in an ownership transfer.

### Output

The action returns the metadata of the copied folder from the Dropbox API.

### Example Usage

Use this action when you need to duplicate a folder and all its contents within your Dropbox storage.