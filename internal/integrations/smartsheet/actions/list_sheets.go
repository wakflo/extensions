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
	"github.com/wakflo/extensions/internal/integrations/smartsheet/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listSheetProps struct {
	Name string `json:"name"`
}

type ListSheetsAction struct{}

func (c *ListSheetsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c ListSheetsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c ListSheetsAction) Name() string {
	return "List Sheets"
}

func (c ListSheetsAction) Description() string {
	return "list all sheets"
}

func (c ListSheetsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listSheetDocs,
	}
}

func (c ListSheetsAction) Icon() *string {
	return nil
}

func (c ListSheetsAction) SampleData() sdkcore.JSON {
	return nil
}

func (c ListSheetsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (c ListSheetsAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c ListSheetsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[listSheetProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/2.0/sheets"

	sheets, err := shared.GetSmartSheetClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return sheets, nil
}

func NewListSheetAction() sdk.Action {
	return &ListSheetsAction{}
}
