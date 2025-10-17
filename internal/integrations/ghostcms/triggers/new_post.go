package triggers

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/ghostcms/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

//go:embed new_post.md
var newPostDocs string

type newPostTriggerProps struct {
	Status  string `json:"status"`
	Include string `json:"include"`
}

type NewPostTrigger struct{}

func (t *NewPostTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_post",
		DisplayName:   "New Post",
		Description:   "Triggers when a new post is created in your Ghost publication",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newPostDocs,
		SampleOutput: map[string]any{
			"id":           "63f8d9a0e4b0d7001c3e3b1a",
			"uuid":         "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			"title":        "New Blog Post",
			"slug":         "new-blog-post",
			"html":         "<p>This is the content of the new blog post.</p>",
			"status":       "published",
			"featured":     false,
			"visibility":   "public",
			"created_at":   "2024-02-24T10:30:00.000Z",
			"updated_at":   "2024-02-24T10:30:00.000Z",
			"published_at": "2024-02-24T10:30:00.000Z",
			"url":          "https://myblog.ghost.io/new-blog-post/",
			"authors": []map[string]any{
				{
					"id":    "1",
					"name":  "John Doe",
					"slug":  "john-doe",
					"email": "john@example.com",
				},
			},
			"tags": []map[string]any{
				{
					"id":   "1",
					"name": "Technology",
					"slug": "technology",
				},
			},
		},
	}
}

func (t *NewPostTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewPostTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewPostTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("ghost-new-post", "New Post")

	form.SelectField("status", "Status Filter").
		Required(false).
		HelpText("Filter for posts with specific status").
		AddOptions([]*smartform.Option{
			{Value: "all", Label: "All"},
			{Value: "published", Label: "Published"},
			{Value: "draft", Label: "Draft"},
			{Value: "scheduled", Label: "Scheduled"},
		}...).
		DefaultValue("all")

	form.TextField("include", "Include Relations").
		Required(false).
		HelpText("Comma-separated list of relations to include (e.g., 'authors,tags')").
		DefaultValue("authors,tags")

	schema := form.Build()
	return schema
}

func (t *NewPostTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Not needed for polling triggers
	return nil
}

func (t *NewPostTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	// Not needed for polling triggers
	return nil
}

func (t *NewPostTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newPostTriggerProps](ctx)
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

	lastRunTime := ctx.LastRun()

	// Build query parameters
	params := "limit=100&order=created_at desc"

	if input.Status != "" && input.Status != "all" {
		params += fmt.Sprintf("&filter=status:%s", input.Status)
	}

	if input.Include != "" {
		params += fmt.Sprintf("&include=%s", input.Include)
	}

	// Add time filter if this is not the first run
	if lastRunTime != nil {
		// Ghost API expects ISO 8601 format
		timeFilter := lastRunTime.UTC().Format(time.RFC3339)
		if input.Status != "" && input.Status != "all" {
			params += fmt.Sprintf("+created_at:>'%s'", timeFilter)
		} else {
			params += fmt.Sprintf("&filter=created_at:>'%s'", timeFilter)
		}
	}

	endpoint := fmt.Sprintf("/posts/?%s", params)

	response, err := client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	// Extract posts from response
	posts, ok := response["posts"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: posts field is not an array")
	}

	// Return empty array if no new posts
	if len(posts) == 0 {
		return []interface{}{}, nil
	}

	// Return posts in chronological order (oldest first) for proper processing
	reversedPosts := make([]interface{}, len(posts))
	for i, post := range posts {
		reversedPosts[len(posts)-1-i] = post
	}

	return reversedPosts, nil
}

func (t *NewPostTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewPostTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":           "63f8d9a0e4b0d7001c3e3b1a",
		"uuid":         "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"title":        "New Blog Post",
		"slug":         "new-blog-post",
		"html":         "<p>This is the content of the new blog post.</p>",
		"status":       "published",
		"featured":     false,
		"visibility":   "public",
		"created_at":   "2024-02-24T10:30:00.000Z",
		"updated_at":   "2024-02-24T10:30:00.000Z",
		"published_at": "2024-02-24T10:30:00.000Z",
		"url":          "https://myblog.ghost.io/new-blog-post/",
	}
}

func NewNewPostTrigger() sdk.Trigger {
	return &NewPostTrigger{}
}
