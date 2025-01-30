package actions

import (
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type CreateFileAction struct {
}

func (c *CreateFileAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c *CreateFileAction) Name() string {
	return "Create File"
}

func (c *CreateFileAction) Description() string {
	return "Create a new file in Google Drive"
}

func (c *CreateFileAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &createFileDocs,
	}
}

func (c *CreateFileAction) Icon() *string {
	return nil
}

func (c *CreateFileAction) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (c *CreateFileAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (c *CreateFileAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c *CreateFileAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	//TODO implement me
	return nil, nil
}

func NewCreateFileAction() integration.Action {
	return &CreateFileAction{}
}
