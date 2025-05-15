package actions

import _ "embed"

//go:embed copy_file.md
var copyFileDocs string

//go:embed create_new_folder.md
var createFolderDocs string

//go:embed copy_folder.md
var copyFolderDocs string

//go:embed move_file.md
var moveFileDocs string

//go:embed move_folder.md
var moveFolderDocs string

//go:embed delete_file.md
var deleteFileDocs string

//go:embed delete_folder.md
var deleteFolderDocs string

//go:embed list_folder_content.md
var listFolderDocs string

//go:embed get_file_link.md
var getFileLinkDocs string
