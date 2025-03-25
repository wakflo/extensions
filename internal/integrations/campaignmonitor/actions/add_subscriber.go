package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/campaignmonitor/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type customField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type addSubscriberActionProps struct {
	ListID                string        `json:"listId"`
	Email                 string        `json:"email"`
	Name                  string        `json:"name"`
	CustomFields          []customField `json:"customFields"`
	ConsentToTrack        string        `json:"consentToTrack"`
	ConsentToSendSMS      string        `json:"consentToSendSMS"`
	Resubscribe           bool          `json:"resubscribe"`
	RestartAutoresponders bool          `json:"restartAutoresponders"`
}

type AddSubscriberAction struct{}

func (a *AddSubscriberAction) Name() string {
	return "Add Subscriber"
}

func (a *AddSubscriberAction) Description() string {
	return "Add a new subscriber to a specific list."
}

func (a *AddSubscriberAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddSubscriberAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addSubscriberDocs,
	}
}

func (a *AddSubscriberAction) Icon() *string {
	icon := "mdi:account-plus"
	return &icon
}

func (a *AddSubscriberAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"listId": shared.GetCreateSendSubscriberListsInput(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("The email address of the subscriber.").
			SetRequired(true).
			Build(),
		"name": autoform.NewShortTextField().
			SetDisplayName("Name").
			SetDescription("The name of the subscriber.").
			SetRequired(false).
			Build(),
		"customFields": autoform.NewArrayField().
			SetDisplayName("Custom Fields").
			SetDescription("Custom fields to add for the subscriber.").
			SetRequired(false).
			SetItems(
				autoform.NewObjectField().
					SetProperties(map[string]*sdkcore.AutoFormSchema{
						"key": autoform.NewShortTextField().
							SetDisplayName("Key").
							SetDescription("The key/name of the custom field.").
							SetRequired(true).
							Build(),
						"value": autoform.NewShortTextField().
							SetDisplayName("Value").
							SetDescription("The value of the custom field.").
							SetRequired(true).
							Build(),
					}).
					Build(),
			).
			Build(),
		"consentToTrack": autoform.NewSelectField().
			SetDisplayName("Consent To Track").
			SetDescription("Whether the subscriber has given consent to track their opens and clicks.").
			SetRequired(true).
			SetOptions([]*sdkcore.AutoFormSchema{
				autoform.NewShortTextField().SetDisplayName("Yes").SetDefaultValue("Yes").Build(),
				autoform.NewShortTextField().SetDisplayName("No").SetDefaultValue("No").Build(),
				autoform.NewShortTextField().SetDisplayName("Unchanged").SetDefaultValue("Unchanged").Build(),
			}).
			SetDefaultValue("Yes").
			Build(),
		"consentToSendSMS": autoform.NewSelectField().
			SetDisplayName("Consent To Send SMS").
			SetDescription("Whether the subscriber has given consent to receive SMS messages.").
			SetRequired(false).
			SetOptions([]*sdkcore.AutoFormSchema{
				autoform.NewShortTextField().SetDisplayName("Yes").SetDefaultValue("Yes").Build(),
				autoform.NewShortTextField().SetDisplayName("No").SetDefaultValue("No").Build(),
				autoform.NewShortTextField().SetDisplayName("Unchanged").SetDefaultValue("Unchanged").Build(),
			}).
			Build(),
		"resubscribe": autoform.NewBooleanField().
			SetDisplayName("Resubscribe").
			SetDescription("Whether to resubscribe the subscriber if they have previously unsubscribed.").
			SetRequired(false).
			SetDefaultValue(false).
			Build(),
		"restartAutoresponders": autoform.NewBooleanField().
			SetDisplayName("Restart Autoresponders").
			SetDescription("Whether to restart autoresponders for the subscriber.").
			SetRequired(false).
			SetDefaultValue(false).
			Build(),
	}
}

func (a *AddSubscriberAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	// Get the input parameters
	input, err := sdk.InputToTypeSafely[addSubscriberActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.ListID == "" {
		return nil, fmt.Errorf("list ID is required")
	}

	if input.Email == "" {
		return nil, fmt.Errorf("email address is required")
	}

	// Build the endpoint - correct endpoint format for Campaign Monitor
	endpoint := fmt.Sprintf("subscribers/%s.json", input.ListID)

	// Prepare the request body
	body := map[string]interface{}{
		"EmailAddress":          input.Email,
		"Resubscribe":           input.Resubscribe,
		"RestartAutoresponders": input.RestartAutoresponders,
		"ConsentToTrack":        input.ConsentToTrack,
	}

	// Add optional fields if provided
	if input.Name != "" {
		body["Name"] = input.Name
	}

	if input.ConsentToSendSMS != "" {
		body["ConsentToSendSMS"] = input.ConsentToSendSMS
	}

	// Add custom fields if provided
	if len(input.CustomFields) > 0 {
		customFields := make([]map[string]string, 0, len(input.CustomFields))
		for _, cf := range input.CustomFields {
			customFields = append(customFields, map[string]string{
				"Key":   cf.Key,
				"Value": cf.Value,
			})
		}
		body["CustomFields"] = customFields
	}

	// Call Campaign Monitor API with correct parameters
	response, err := shared.GetCampaignMonitorClient(
		ctx.Auth.Extra["api-key"],
		ctx.Auth.Extra["client-id"],
		endpoint,
		http.MethodPost,
		body)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": "Subscriber added successfully",
		"result":  response,
	}, nil
}

func (a *AddSubscriberAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddSubscriberAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
		"EmailAddress": "subscriber@example.com",
		"Name":         "John Smith",
		"Status":       "Active",
		"Message":      "Subscriber was added successfully",
	}
}

func (a *AddSubscriberAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddSubscriberAction() sdk.Action {
	return &AddSubscriberAction{}
}
