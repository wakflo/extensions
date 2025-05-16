package shared

import (
	"fmt"

	"github.com/aftership/tracking-sdk-go/v5"
	"github.com/aftership/tracking-sdk-go/v5/model"
	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("aftership-auth", "Aftership Auth Screen", smartform.AuthStrategyCustom)
	_    = form.
		TextField("api-key", "API Key").
		Placeholder("Enter API Key").
		Required(true).
		HelpText("API Application Key").
		Build()

	AfterShipSharedAuth = form.Build()
)

func RegisterSlugProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getSlugs := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		afterShipSdk, err := tracking.New(tracking.WithApiKey(authCtx.Extra["api-key"]))
		if err != nil {
			fmt.Println("Error initializing AfterShip SDK:", err)
			return nil, err
		}

		result, err := afterShipSdk.Courier.GetAllCouriers().Execute()
		if err != nil {
			fmt.Println("Error fetching couriers:", err)
			return nil, err
		}

		slugss := result
		total := slugss.Couriers

		mappedSlugs := arrutil.Map[model.Courier, map[string]any](total, func(input model.Courier) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.Slug,
				"name": input.Name,
			}, true
		})

		return ctx.Respond(mappedSlugs, len(mappedSlugs))
	}

	return form.SelectField("slug", "Couriers").
		Placeholder("Select a courier").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSlugs)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a courier")
}
