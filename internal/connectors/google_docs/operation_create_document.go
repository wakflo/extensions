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

package googledocs

import (
	"context"
	"errors"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createDocumentOperationProps struct {
	Name string `json:"name"`
	// Body any `json:"body"`
}

type CreateDocumentOperation struct {
	options *sdk.OperationInfo
}

func NewCreateDocumentOperation() *CreateDocumentOperation {
	return &CreateDocumentOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Document",
			Description: "Search for document by name",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"name": autoform.NewShortTextField().
					SetDisplayName("Document Name").
					SetDescription("The name of the document.").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateDocumentOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[createDocumentOperationProps](ctx)
	docService, err := docs.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	document, err := docService.Documents.Create(&docs.Document{
		Title: input.Name,
	}).
		Do()
	return document, err
}

func (c *CreateDocumentOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateDocumentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
