package actions

import _ "embed"

//go:embed create_file.md
var createFileDocs string

//go:embed create_folder.md
var createFolderDocs string

//go:embed duplicate_file.md
var duplicateFileDocs string

//go:embed get_file.md
var getFileDocs string

//go:embed list_files.md
var listFilesDocs string

//go:embed list_folders.md
var listFoldersDocs string

//go:embed read_file_content.md
var readFileContentDocs string

//go:embed upload_file.md
var uploadFileDocs string
