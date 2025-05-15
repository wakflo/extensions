## List Folder Content Action

This action retrieves a list of all files and folders within a specified Dropbox folder.

### Input Parameters

- **From Path** (required): The path of the folder to be listed (e.g. /folder1). Use an empty string for the root folder.
- **Limit** (optional): The maximum number of results to return (between 1 and 2000). Default is 2000 if not specified.
- **Recursive** (optional): If set to true, the list folder operation will be applied recursively to all subfolders and the response will contain contents of all subfolders.

### Output

The action returns a JSON object containing metadata about the files and folders within the specified folder.

### Example Usage

Use this action when you need to inventory the contents of a folder in your Dropbox account.