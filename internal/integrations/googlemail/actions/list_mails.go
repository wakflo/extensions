package actions

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type listMailsActionProps struct {
	Label     string `json:"label"`
	PageSize  string `json:"pageSize,omitempty"`
	PageToken string `json:"pageToken,omitempty"`
}

type EmailSummary struct {
	ID       string `json:"id"`
	ThreadID string `json:"threadId"`
	Subject  string `json:"subject"`
}

type ListMailsResponse struct {
	Emails             []EmailSummary `json:"emails"`
	NextPageToken      string         `json:"nextPageToken,omitempty"`
	ResultSizeEstimate int64          `json:"resultSizeEstimate"`
	HasMore            bool           `json:"hasMore"`
	PageSize           int64          `json:"pageSize"`
}

type ListMailsAction struct{}

// Metadata returns metadata about the action
func (a *ListMailsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_mails",
		DisplayName:   "List Mails",
		Description:   "Retrieve a paginated list of emails with ID, ThreadID, and Subject from your Gmail account.",
		Type:          core.ActionTypeAction,
		Documentation: listMailsDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"emails": []map[string]any{
				{
					"id":       "18abc123def",
					"threadId": "18abc123def",
					"subject":  "Meeting Tomorrow",
				},
			},
			"nextPageToken":      "CAIQABiAwtjX7",
			"resultSizeEstimate": 1500,
			"hasMore":            true,
			"pageSize":           50,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListMailsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_mails", "List Mails")

	form.TextField("label", "Label").
		Placeholder("inbox").
		HelpText("The mail label to read from (e.g, inbox, sent, drafts, spam, trash)").
		Required(true)

	form.TextField("pageSize", "Page Size").
		Placeholder("50").
		HelpText("Number of emails per page (default: 50, max: 200 for subject-only fetch)").
		Required(false)

	form.TextField("pageToken", "Page Token").
		Placeholder("").
		HelpText("Token for fetching the next page of results").
		Required(false)

	schema := form.Build()

	return schema
}

func (a *ListMailsAction) Auth() *core.AuthMetadata {
	return nil
}

func fetchEmailsWithSubjects(gmailService *gmail.Service, messages []*gmail.Message) []EmailSummary {
	emails := make([]EmailSummary, 0, len(messages))

	batchSize := 25
	for i := 0; i < len(messages); i += batchSize {
		end := i + batchSize
		if end > len(messages) {
			end = len(messages)
		}

		type result struct {
			index int
			email EmailSummary
			err   error
		}
		results := make(chan result, end-i)

		for j := i; j < end; j++ {
			go func(idx int, msg *gmail.Message) {
				message, err := gmailService.Users.Messages.Get("me", msg.Id).
					Format("metadata").
					MetadataHeaders("Subject").
					Fields("id,threadId,payload/headers").
					Do()
				if err != nil {
					results <- result{index: idx, err: err}
					return
				}

				// Extract subject from headers
				var subject string
				if message.Payload != nil {
					for _, header := range message.Payload.Headers {
						if header.Name == "Subject" {
							subject = header.Value
							break
						}
					}
				}

				results <- result{
					index: idx,
					email: EmailSummary{
						ID:       message.Id,
						ThreadID: message.ThreadId,
						Subject:  subject,
					},
				}
			}(j, messages[j])
		}

		// Collect results
		batchEmails := make([]EmailSummary, end-i)
		for j := 0; j < end-i; j++ {
			res := <-results
			if res.err == nil {
				batchEmails[res.index-i] = res.email
			}
		}

		// Add non-empty results to final list
		for _, email := range batchEmails {
			if email.ID != "" {
				emails = append(emails, email)
			}
		}
	}

	return emails
}

// Perform executes the action with the given context and input
func (a *ListMailsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listMailsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.Label == "" {
		return nil, errors.New("label name is required")
	}

	pageSize := parsePageSize(input.PageSize)

	// Build query
	query := fmt.Sprintf("in:%s", input.Label)

	// Create list call with pagination
	listCall := gmailService.Users.Messages.List("me").
		Q(query).
		MaxResults(pageSize).
		Fields("messages(id,threadId),nextPageToken,resultSizeEstimate")

	if input.PageToken != "" {
		listCall = listCall.PageToken(input.PageToken)
	}

	listResponse, err := listCall.Do()
	if err != nil {
		return nil, err
	}

	// Fetch emails with just the subject header
	emails := fetchEmailsWithSubjects(gmailService, listResponse.Messages)

	// Build response
	response := ListMailsResponse{
		Emails:             emails,
		NextPageToken:      listResponse.NextPageToken,
		ResultSizeEstimate: listResponse.ResultSizeEstimate,
		HasMore:            listResponse.NextPageToken != "",
		PageSize:           pageSize,
	}

	fmt.Println(reflect.TypeOf(input.PageSize))

	return response, nil
}

func parsePageSize(pageSizeStr string) int64 {
	// Default value if empty or invalid
	const defaultPageSize int64 = 50
	const maxPageSize int64 = 200

	if pageSizeStr == "" {
		return defaultPageSize
	}

	// Try to parse the string
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		// If parsing fails, return default
		return defaultPageSize
	}

	// Apply bounds
	if pageSize <= 0 {
		return defaultPageSize
	}
	if pageSize > maxPageSize {
		return maxPageSize
	}

	return pageSize
}

func NewListMailsAction() sdk.Action {
	return &ListMailsAction{}
}
