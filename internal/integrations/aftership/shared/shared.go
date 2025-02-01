package shared

import (
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var AfterShipSharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().SetDisplayName("API Key").
			SetDescription("API Application Key").
			SetRequired(true).
			Build(),
	}).
	Build()

	//  func getSlugInput() *sdkcore.AutoFormSchema {
	//	getSlugs := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
	//		afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	//		if err != nil {
	//			fmt.Println("Error initializing AfterShip SDK:", err)
	//			return nil, err
	//		}
	//
	//		result, err := afterShipSdk.Courier.GetAllCouriers().Execute()
	//		if err != nil {
	//			fmt.Println("Error fetching couriers:", err)
	//			return nil, err
	//		}
	//
	//		slugss := result
	//
	//		total := slugss.Couriers
	//
	//		mappedSlugs := arrutil.Map[model.Courier, map[string]any](total, func(input model.Courier) (target map[string]any, find bool) {
	//			return map[string]any{
	//				"id":   input.Slug,
	//				"name": input.Name,
	//			}, true
	//		})
	//
	//		return mappedSlugs, nil
	//	}
	//
	//	return autoform.NewDynamicField(sdkcore.String).
	//		SetDisplayName("Couriers").
	//		SetDescription("Select a courier").
	//		SetDynamicOptions(&getSlugs).
	//		SetRequired(true).Build()
	//}
