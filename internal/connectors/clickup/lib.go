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

package clickup

import (
	sdk "github.com/wakflo/go-sdk/connector"
)

func NewConnector() (*sdk.ConnectorPlugin, error) {
	return sdk.CreateConnector(&sdk.CreateConnectorArgs{
		Name:        "Clickup",
		Description: "All-in-one productivity platform",
		Logo:        "simple-icons:clickup",
		Version:     "0.0.1",
		Group:       sdk.ConnectorGroupApps,
		Authors:     []string{"Wakflo <integrations@wakflo.com>"},
		Triggers:    []sdk.ITrigger{},
		Operations: []sdk.IOperation{
			NewGetSpaceOperation(),
			NewGetSpacesOperation(),
			NewGetTasksOperation(),
			NewGetTaskOperation(),
			NewCreateTaskOperation(),
			NewCreateSpaceOperation(),
			NewSearchTaskOperation(),
			NewGetListOperation(),
			NewUpdateTaskOperation(),
			NewDeleteTaskOperation(),
			NewCreateFolderlessListOperation(),
			NewGetFolderlesslistOperation(),
			NewCreateFolderOperation(),
			NewGetFolderOperation(),
			NewGetFoldersOperation(),
			NewUpdateSpaceOperation(),
		},
	})
}
