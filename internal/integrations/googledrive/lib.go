package googledrive

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/googledrive/actions"
	"github.com/wakflo/extensions/internal/integrations/googledrive/shared"
	"github.com/wakflo/extensions/internal/integrations/googledrive/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewGoogleDrive(), Flow, ReadME)

type GoogleDrive struct{}

func (n *GoogleDrive) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedGoogleDriveAuth,
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
