package actions

import (
	_ "embed"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/ghostcms/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updatePostActionProps struct {
	PostID          string   `json:"post_id"`
	Title           string   `json:"title,omitempty"`
	Content         string   `json:"content,omitempty"`
	ContentFormat   string   `json:"content_format,omitempty"`
	Slug            string   `json:"slug,omitempty"`
	Status          string   `json:"status,omitempty"`
	Featured        *bool    `json:"featured,omitempty"`
	FeatureImage    string   `json:"feature_image,omitempty"`
	CustomExcerpt   string   `json:"custom_excerpt,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Authors         []string `json:"authors,omitempty"`
	MetaTitle       string   `json:"meta_title,omitempty"`
	MetaDescription string   `json:"meta_description,omitempty"`
}

type UpdatePostAction struct{}

func (a *UpdatePostAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_post",
		DisplayName:   "Update Post",
		Description:   "Updates an existing post in Ghost",
		Type:          core.ActionTypeAction,
		Documentation: updatePostDocs,
		SampleOutput: map[string]any{
			"id":           "63f8d9a0e4b0d7001c3e3b1a",
			"uuid":         "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			"title":        "Updated Blog Post",
			"slug":         "updated-blog-post",
			"html":         "<p>This is the updated content.</p>",
			"status":       "published",
			"featured":     true,
			"visibility":   "public",
			"created_at":   "2024-02-24T10:30:00.000Z",
			"updated_at":   "2024-02-24T15:45:00.000Z",
			"published_at": "2024-02-24T10:30:00.000Z",
			"url":          "https://myblog.ghost.io/updated-blog-post/",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *UpdatePostAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_post", "Update Post")

	form.TextField("post_id", "Post ID").
		Required(true).
		HelpText("The ID of the post to update")

	form.TextField("title", "Title").
		Required(false).
		HelpText("The new title of the post (leave empty to keep current)")

	form.TextareaField("content", "Content").
		Required(false).
		HelpText("The new content of the post (leave empty to keep current)")

	form.SelectField("content_format", "Content Format").
		Required(false).
		HelpText("The format of the content").
		AddOptions([]*smartform.Option{
			{Value: "html", Label: "HTML"},
			{Value: "markdown", Label: "Markdown"},
			{Value: "lexical", Label: "Lexical"},
		}...).
		DefaultValue("html")

	form.TextField("slug", "Slug").
		Required(false).
		HelpText("The URL slug for the post")

	form.SelectField("status", "Status").
		Required(false).
		HelpText("The publication status of the post").
		AddOptions([]*smartform.Option{
			{Value: "", Label: "Keep Current"},
			{Value: "draft", Label: "Draft"},
			{Value: "published", Label: "Published"},
			{Value: "scheduled", Label: "Scheduled"},
		}...)

	form.CheckboxField("featured", "Featured").
		Required(false).
		HelpText("Whether to feature this post")

	form.TextField("feature_image", "Feature Image URL").
		Required(false).
		HelpText("URL of the feature image")

	form.TextareaField("custom_excerpt", "Custom Excerpt").
		Required(false).
		HelpText("Custom excerpt for the post")

	form.ArrayField("tags", "Tags").
		Required(false).
		HelpText("Tags for the post (use tag names or IDs)")

	form.ArrayField("authors", "Authors").
		Required(false).
		HelpText("Authors for the post (use author IDs or email addresses)")

	form.TextField("meta_title", "Meta Title").
		Required(false).
		HelpText("SEO meta title")

	form.TextareaField("meta_description", "Meta Description").
		Required(false).
		HelpText("SEO meta description")

	schema := form.Build()
	return schema
}

func (a *UpdatePostAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *UpdatePostAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updatePostActionProps](ctx)
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

	// First, get the current post to get the updated_at timestamp
	currentPost, err := client.Get(fmt.Sprintf("/posts/%s/", input.PostID))
	if err != nil {
		return nil, fmt.Errorf("failed to get current post: %w", err)
	}

	posts, ok := currentPost["posts"].([]interface{})
	if !ok || len(posts) == 0 {
		return nil, fmt.Errorf("post not found")
	}

	post, ok := posts[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid post format")
	}

	// Build the update object with only changed fields
	updateData := map[string]interface{}{
		"updated_at": post["updated_at"], // Required for conflict detection
	}

	if input.Title != "" {
		updateData["title"] = input.Title
	}

	if input.Content != "" {
		switch input.ContentFormat {
		case "markdown":
			updateData["mobiledoc"] = fmt.Sprintf(`{"version":"0.3.1","atoms":[],"cards":[["markdown",{"markdown":"%s"}]],"markups":[],"sections":[[10,0]]}`, input.Content)
		case "lexical":
			updateData["lexical"] = input.Content
		case "html":
			fallthrough
		default:
			updateData["html"] = input.Content
		}
	}

	if input.Slug != "" {
		updateData["slug"] = input.Slug
	}

	if input.Status != "" {
		updateData["status"] = input.Status
	}

	if input.Featured != nil {
		updateData["featured"] = *input.Featured
	}

	if input.FeatureImage != "" {
		updateData["feature_image"] = input.FeatureImage
	}

	if input.CustomExcerpt != "" {
		updateData["custom_excerpt"] = input.CustomExcerpt
	}

	if input.MetaTitle != "" {
		updateData["meta_title"] = input.MetaTitle
	}

	if input.MetaDescription != "" {
		updateData["meta_description"] = input.MetaDescription
	}

	// Handle tags
	if len(input.Tags) > 0 {
		tags := make([]map[string]string, len(input.Tags))
		for i, tag := range input.Tags {
			tags[i] = map[string]string{"name": tag}
		}
		updateData["tags"] = tags
	}

	// Handle authors
	if len(input.Authors) > 0 {
		authors := make([]map[string]string, len(input.Authors))
		for i, author := range input.Authors {
			authors[i] = map[string]string{"id": author}
		}
		updateData["authors"] = authors
	}

	// Update the post
	response, err := client.Put(fmt.Sprintf("/posts/%s/", input.PostID), map[string]interface{}{
		"posts": []interface{}{updateData},
	})
	if err != nil {
		return nil, err
	}

	// Extract the updated post
	if posts, ok := response["posts"].([]interface{}); ok && len(posts) > 0 {
		return posts[0], nil
	}

	return response, nil
}

func NewUpdatePostAction() sdk.Action {
	return &UpdatePostAction{}
}
