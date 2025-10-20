package actions

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/ghostcms/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listPostsActionProps struct {
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
	Status  string `json:"status,omitempty"`
	Filter  string `json:"filter,omitempty"`
	Order   string `json:"order"`
	Include string `json:"include,omitempty"`
	Fields  string `json:"fields,omitempty"`
	Formats string `json:"formats,omitempty"`
}

type ListPostsAction struct{}

func (a *ListPostsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_posts",
		DisplayName:   "List Posts",
		Description:   "Lists all posts with filtering and pagination options",
		Type:          core.ActionTypeAction,
		Documentation: listPostsDocs,
		SampleOutput: map[string]any{
			"posts": []map[string]any{
				{
					"id":           "63f8d9a0e4b0d7001c3e3b1a",
					"uuid":         "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
					"title":        "Sample Blog Post",
					"slug":         "sample-blog-post",
					"html":         "<p>This is a sample blog post content.</p>",
					"status":       "published",
					"featured":     false,
					"visibility":   "public",
					"created_at":   "2024-02-24T10:30:00.000Z",
					"updated_at":   "2024-02-24T10:30:00.000Z",
					"published_at": "2024-02-24T10:30:00.000Z",
					"url":          "https://myblog.ghost.io/sample-blog-post/",
				},
			},
			"meta": map[string]any{
				"pagination": map[string]any{
					"page":  1,
					"limit": 15,
					"pages": 1,
					"total": 1,
					"next":  nil,
					"prev":  nil,
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListPostsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_posts", "List Posts")

	form.NumberField("limit", "Limit").
		Required(false).
		HelpText("Number of posts to retrieve (1-100)").
		DefaultValue(15)

	form.NumberField("page", "Page").
		Required(false).
		HelpText("Page number for pagination").
		DefaultValue(1)

	form.SelectField("status", "Status Filter").
		Required(false).
		HelpText("Filter posts by status").
		AddOptions([]*smartform.Option{
			{Value: "all", Label: "All"},
			{Value: "published", Label: "Published"},
			{Value: "draft", Label: "Draft"},
			{Value: "scheduled", Label: "Scheduled"},
		}...).
		DefaultValue("all")

	form.TextField("filter", "Custom Filter").
		Required(false).
		HelpText("NQL filter query (e.g., 'featured:true', 'tag:news')")

	form.SelectField("order", "Order").
		Required(false).
		HelpText("Sort order for posts").
		AddOptions([]*smartform.Option{
			{Value: "published_at desc", Label: "Published Descending"},
			{Value: "published_at asc", Label: "Published Ascending"},
			{Value: "created_at desc", Label: "Created Descending"},
			{Value: "created_at asc", Label: "Created Ascending"},
			{Value: "updated_at desc", Label: "Updated Descending"},
			{Value: "updated_at asc", Label: "Updated Ascending"},
		}...).
		DefaultValue("published_at desc")

	form.TextField("include", "Include").
		Required(false).
		HelpText("Comma-separated list of relations to include (e.g., 'authors,tags')")

	form.TextField("fields", "Fields").
		Required(false).
		HelpText("Comma-separated list of fields to include in response")

	form.TextField("formats", "Formats").
		Required(false).
		HelpText("Comma-separated list of formats to include (html, plaintext)").
		DefaultValue("html")

	schema := form.Build()
	return schema
}

func (a *ListPostsAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *ListPostsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listPostsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client, err := shared.NewGhostClient(authCtx.Extra)
	if err != nil {
		return nil, err
	}

	// Set defaults if not provided
	if input.Limit == 0 {
		input.Limit = 15
	}
	if input.Page == 0 {
		input.Page = 1
	}

	// Build query parameters
	params := []string{
		fmt.Sprintf("limit=%d", input.Limit),
		fmt.Sprintf("page=%d", input.Page),
	}

	if input.Status != "" && input.Status != "all" {
		params = append(params, fmt.Sprintf("filter=status:%s", input.Status))
	} else if input.Filter != "" {
		params = append(params, fmt.Sprintf("filter=%s", input.Filter))
	}

	if input.Order != "" {
		params = append(params, fmt.Sprintf("order=%s", input.Order))
	}

	if input.Include != "" {
		params = append(params, fmt.Sprintf("include=%s", input.Include))
	}

	if input.Fields != "" {
		params = append(params, fmt.Sprintf("fields=%s", input.Fields))
	}

	if input.Formats != "" {
		params = append(params, fmt.Sprintf("formats=%s", input.Formats))
	}

	endpoint := fmt.Sprintf("/posts/?%s", strings.Join(params, "&"))

	response, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewListPostsAction() sdk.Action {
	return &ListPostsAction{}
}
