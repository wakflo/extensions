## Move File Action

This action moves a file from one location to another within your Dropbox account.

### Input Parameters

- **From Path** (required): The current path of the file (e.g. /folder1/oldfile.txt).
- **To Path** (required): The new path for the file (e.g. /folder2/newfile.txt).
- **Auto Rename** (optional): If there's a conflict, have the Dropbox server try to autorename the file to avoid conflict.
- **Allow Ownership Transfer** (optional): Allows move by owner even if it would result in an ownership transfer.

### Output

The action returns metadata about the moved file from the Dropbox API.

### Example Usage

Use this action when you need to reorganize files within your Dropbox storage.