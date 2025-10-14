package actions

import _ "embed"

//go:embed list_videos.md
var listVideosDocs string

//go:embed get_video.md
var getVideoDocs string

//go:embed update_video.go
var updateVideoDocs string

//go:embed upload_video.go
var uploadVideoDocs string

//go:embed: download_caption.go
var downloadCaptionDocs string
