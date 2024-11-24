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

package trackingmore

import (
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createTrackingsProps struct {
	TrackingNumber string                   `json:"tracking_number"`
	CourierCode    string                   `json:"courier_code"`
	Body           []map[string]interface{} `json:"body"`
}

type CreateTrackingsOperation struct {
	options *sdk.OperationInfo
}

func NewCreateTrackingsOperation() *CreateTrackingsOperation {
	return &CreateTrackingsOperation{
		options: &sdk.OperationInfo{
			Name:        "Create batch tracking",
			Description: "create a batch tracking",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tracking_number": autoform.NewShortTextField().
					SetDisplayName("Tracking Number").
					SetDescription("Tracking number of a package.").
					SetRequired(true).
					Build(),
				"body": autoform.NewArrayField().
					SetDisplayName("Courier Code").
					SetDescription("Courier code").
					SetItems(autoform.NewObjectField().SetDisplayName("").SetProperties(
						map[string]*sdkcore.AutoFormSchema{
							"courier_code": autoform.NewSelectField().
								SetDisplayName("Courier Code").
								SetDescription("Courier code").
								SetOptions(courierCodes).
								SetRequired(true).
								Build(),
							"tracking_number": autoform.NewShortTextField().
								SetDisplayName("Tracking Number").
								SetDescription("Tracking number of a package.").
								SetRequired(true).
								Build(),
						}).Build()).
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

func (c *CreateTrackingsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[createTrackingsProps](ctx)

	applicationKey := ctx.Auth.Extra["key"]

	//  payload := []map[string]interface{}{
	//	{
	//		"courier_code":    "usps",
	//		"tracking_number": "9269990312443844954410",
	//	},
	//	{
	//		"courier_code":    "fedex",
	//		"tracking_number": "608285157867",
	//	},
	//}

	response, err := createBatchTracking(applicationKey, input.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *CreateTrackingsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateTrackingsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
