package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zendeskapp/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getGroupsActionProps struct{}

type GetGroupsAction struct{}

func (a *GetGroupsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_groups",
		DisplayName:   "Get Groups",
		Description:   "Get Groups: Retrieves all the groups from your Zendesk account, providing access to team structures and organization information.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getGroupsDocs,
		SampleOutput: map[string]any{
			"groups": []map[string]any{
				{
					"id":   123456,
					"name": "Support Team",
				},
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetGroupsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_groups", "Get Groups")

	schema := form.Build()
	return schema
}

func (a *GetGroupsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	// Extract authentication details
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Check if required auth details are present
	email, ok := authCtx.Extra["email"]
	if !ok || email == "" {
		return nil, errors.New("missing zendesk email")
	}

	apiToken, ok := authCtx.Extra["api-token"]
	if !ok || apiToken == "" {
		return nil, errors.New("missing zendesk api_token")
	}

	subdomain, ok := authCtx.Extra["subdomain"]
	if !ok || subdomain == "" {
		return nil, errors.New("missing zendesk subdomain")
	}

	// Construct the URL
	fullURL := fmt.Sprintf("https://%s.zendesk.com/api/v2/groups.json", subdomain)

	// Make the request
	response, err := shared.ZendeskRequest(http.MethodGet, fullURL, email, apiToken, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching groups: %v", err)
	}

	return response, nil
}

func (a *GetGroupsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetGroupsAction() sdk.Action {
	return &GetGroupsAction{}
}
