## Move Folder Action

This action moves a folder and all its contents from one location to another within your Dropbox account.

### Input Parameters

- **From Path** (required): The current path of the folder (e.g. /folder1/sourceFolder).
- **To Path** (required): The new path for the folder (e.g. /folder2/destinationFolder).
- **Auto Rename** (optional): If there's a conflict, have the Dropbox server try to autorename the folder to avoid conflict.
- **Allow Ownership Transfer** (optional): Allows move by owner even if it would result in an ownership transfer.

### Output

The action returns metadata about the moved folder from the Dropbox API.

### Example Usage

Use this action when you need to reorganize your folder structure within your Dropbox storage.