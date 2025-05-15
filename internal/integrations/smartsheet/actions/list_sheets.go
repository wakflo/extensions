// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/smartsheet/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listSheetProps struct{}

type ListSheetsAction struct{}

func (c *ListSheetsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_sheets",
		DisplayName:   "List Sheets",
		Description:   "list all sheets",
		Type:          core.ActionTypeAction,
		Documentation: listSheetDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

func (c *ListSheetsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_sheets", "List Sheets")

	schema := form.Build()
	return schema
}

func (c *ListSheetsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	_, err := sdk.InputToTypeSafely[listSheetProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := "/2.0/sheets"

	sheets, err := shared.GetSmartSheetClient(authCtx.Token.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return sheets, nil
}

func (c *ListSheetsAction) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Inherit: true,
	}
}

func NewListSheetAction() sdk.Action {
	return &ListSheetsAction{}
}
