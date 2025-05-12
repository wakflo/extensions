package shared

import (
	"github.com/juicycleff/smartform/v1"
)

var (
	form = smartform.NewAuthForm("aftership-auth", "Aftership Auth Screen", smartform.AuthStrategyOAuth2)
	_    = form.
		APIField("api-key", "API Key").
		Placeholder("Enter API Key").
		Required(true).
		HelpText("API Application Key").
		Build()

	AfterShipSharedAuth = form.Build()
)

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
// }
