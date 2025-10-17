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

type createPostActionProps struct {
	Title                  string   `json:"title"`
	Content                string   `json:"content"`
	ContentFormat          string   `json:"content_format"`
	Slug                   string   `json:"slug,omitempty"`
	Status                 string   `json:"status"`
	Featured               bool     `json:"featured"`
	FeatureImage           string   `json:"feature_image,omitempty"`
	CustomExcerpt          string   `json:"custom_excerpt,omitempty"`
	Tags                   []string `json:"tags,omitempty"`
	Authors                []string `json:"authors,omitempty"`
	MetaTitle              string   `json:"meta_title,omitempty"`
	MetaDescription        string   `json:"meta_description,omitempty"`
	OgTitle                string   `json:"og_title,omitempty"`
	OgDescription          string   `json:"og_description,omitempty"`
	TwitterTitle           string   `json:"twitter_title,omitempty"`
	TwitterDescription     string   `json:"twitter_description,omitempty"`
	CodeinjectionHead      string   `json:"codeinjection_head,omitempty"`
	CodeinjectionFoot      string   `json:"codeinjection_foot,omitempty"`
	EmailSubject           string   `json:"email_subject,omitempty"`
	SendEmailWhenPublished bool     `json:"send_email_when_published"`
}

type CreatePostAction struct{}

func (a *CreatePostAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_post",
		DisplayName:   "Create Post",
		Description:   "Creates a new post in your Ghost publication",
		Type:          core.ActionTypeAction,
		Documentation: createPostDocs,
		SampleOutput: map[string]any{
			"id":           "63f8d9a0e4b0d7001c3e3b1a",
			"uuid":         "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			"title":        "My New Blog Post",
			"slug":         "my-new-blog-post",
			"html":         "<p>This is the content of my blog post.</p>",
			"status":       "draft",
			"featured":     false,
			"visibility":   "public",
			"created_at":   "2024-02-24T10:30:00.000Z",
			"updated_at":   "2024-02-24T10:30:00.000Z",
			"published_at": nil,
			"url":          "https://myblog.ghost.io/my-new-blog-post/",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreatePostAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_post", "Create Post")

	form.TextField("title", "Title").
		Required(true).
		HelpText("The title of the post")

	form.TextareaField("content", "Content").
		Required(true).
		HelpText("The content of the post (HTML or Markdown based on format)")

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
		HelpText("The URL slug for the post (auto-generated if not provided)")

	form.SelectField("status", "Status").
		Required(true).
		HelpText("The publication status of the post").
		AddOptions([]*smartform.Option{
			{Value: "draft", Label: "Draft"},
			{Value: "published", Label: "Published"},
			{Value: "scheduled", Label: "Scheduled"},
		}...).
		DefaultValue("draft")

	form.CheckboxField("featured", "Featured").
		Required(false).
		HelpText("Whether to feature this post").
		DefaultValue(false)

	form.TextField("feature_image", "Feature Image URL").
		Required(false).
		HelpText("URL of the feature image")

	form.TextareaField("custom_excerpt", "Custom Excerpt").
		Required(false).
		HelpText("Custom excerpt for the post")

	form.TextField("meta_title", "Meta Title").
		Required(false).
		HelpText("SEO meta title")

	form.TextareaField("meta_description", "Meta Description").
		Required(false).
		HelpText("SEO meta description")

	form.TextField("og_title", "Open Graph Title").
		Required(false).
		HelpText("Open Graph title for social sharing")

	form.TextareaField("og_description", "Open Graph Description").
		Required(false).
		HelpText("Open Graph description for social sharing")

	form.TextField("twitter_title", "Twitter Title").
		Required(false).
		HelpText("Twitter card title")

	form.TextareaField("twitter_description", "Twitter Description").
		Required(false).
		HelpText("Twitter card description")

	form.TextareaField("codeinjection_head", "Code Injection - Head").
		Required(false).
		HelpText("Code to inject into the <head> tag")

	form.TextareaField("codeinjection_foot", "Code Injection - Footer").
		Required(false).
		HelpText("Code to inject before the closing </body> tag")

	form.TextField("email_subject", "Email Subject").
		Required(false).
		HelpText("Subject line for the email newsletter")

	form.CheckboxField("send_email_when_published", "Send Email When Published").
		Required(false).
		HelpText("Whether to send an email newsletter when the post is published").
		DefaultValue(false)

	schema := form.Build()
	return schema
}

func (a *CreatePostAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *CreatePostAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createPostActionProps](ctx)
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

	// Build the post object
	post := map[string]interface{}{
		"title":    input.Title,
		"status":   input.Status,
		"featured": input.Featured,
	}

	// Set content based on format
	switch input.ContentFormat {
	case "markdown":
		post["mobiledoc"] = fmt.Sprintf(`{"version":"0.3.1","atoms":[],"cards":[["markdown",{"markdown":"%s"}]],"markups":[],"sections":[[10,0]]}`, input.Content)
	case "lexical":
		post["lexical"] = input.Content
	case "html":
		fallthrough
	default:
		post["html"] = input.Content
	}

	// Add optional fields
	if input.Slug != "" {
		post["slug"] = input.Slug
	}
	if input.FeatureImage != "" {
		post["feature_image"] = input.FeatureImage
	}
	if input.CustomExcerpt != "" {
		post["custom_excerpt"] = input.CustomExcerpt
	}
	if input.MetaTitle != "" {
		post["meta_title"] = input.MetaTitle
	}
	if input.MetaDescription != "" {
		post["meta_description"] = input.MetaDescription
	}
	if input.OgTitle != "" {
		post["og_title"] = input.OgTitle
	}
	if input.OgDescription != "" {
		post["og_description"] = input.OgDescription
	}
	if input.TwitterTitle != "" {
		post["twitter_title"] = input.TwitterTitle
	}
	if input.TwitterDescription != "" {
		post["twitter_description"] = input.TwitterDescription
	}
	if input.CodeinjectionHead != "" {
		post["codeinjection_head"] = input.CodeinjectionHead
	}
	if input.CodeinjectionFoot != "" {
		post["codeinjection_foot"] = input.CodeinjectionFoot
	}
	if input.EmailSubject != "" {
		post["email_subject"] = input.EmailSubject
	}

	// Handle tags
	if len(input.Tags) > 0 {
		tags := make([]map[string]string, len(input.Tags))
		for i, tag := range input.Tags {
			tags[i] = map[string]string{"name": tag}
		}
		post["tags"] = tags
	}

	// Handle authors
	if len(input.Authors) > 0 {
		authors := make([]map[string]string, len(input.Authors))
		for i, author := range input.Authors {
			authors[i] = map[string]string{"id": author}
		}
		post["authors"] = authors
	}

	if input.SendEmailWhenPublished && input.Status == "published" {
		post["send_email_when_published"] = "all"
	}

	// Create the post
	response, err := client.Post("/posts/", map[string]interface{}{
		"posts": []interface{}{post},
	})
	if err != nil {
		return nil, err
	}

	// Extract the created post
	if posts, ok := response["posts"].([]interface{}); ok && len(posts) > 0 {
		return posts[0], nil
	}

	return response, nil
}

func NewCreatePostAction() sdk.Action {
	return &CreatePostAction{}
}
