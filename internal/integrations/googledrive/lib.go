package googledrive

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googledrive/actions"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/extensions/internal/integrations/googledrive/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleDrive())

type GoogleDrive struct{}

func (n *GoogleDrive) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *GoogleDrive) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedGoogleDriveAuth,
	}
}

func (n *GoogleDrive) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewFolderTrigger(),
		triggers.NewNewFileTrigger(),
	}
}

func (n *GoogleDrive) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewUploadFileAction(),
		actions.NewListFoldersAction(),
		actions.NewListFilesAction(),
		actions.NewReadFileContentAction(),
		actions.NewGetFileAction(),
		actions.NewDuplicateFileAction(),
		actions.NewCreateFolderAction(),
		actions.NewCreateFileAction(),
	}
}

func NewGoogleDrive() sdk.Integration {
	return &GoogleDrive{}
}
