## Copy File Action

This action copies a file from one location to another within Dropbox.

### Input Parameters

- **From Path** (required): The source path of the file (e.g. /folder1/sourcefile.txt).
- **To Path** (required): The destination path for the copied file (e.g. /folder2/destinationfile.txt).
- **Auto Rename** (optional): If there's a conflict, have the Dropbox server try to autorename the file to avoid conflict.
- **Allow Ownership Transfer** (optional): Allows copy by owner even if it would result in an ownership transfer.

### Output

The action returns the metadata of the copied file from the Dropbox API.

### Example Usage

Use this action when you need to duplicate a file within your Dropbox storage.