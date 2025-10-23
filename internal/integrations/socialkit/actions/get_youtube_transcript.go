package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/socialkit/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getYouTubeTranscriptActionProps struct {
	VideoURL string `json:"video_url"`
}

type GetYouTubeTranscriptAction struct{}

func (a *GetYouTubeTranscriptAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_youtube_transcript",
		DisplayName:   "Get YouTube Transcript",
		Description:   "Extract accurate, timestamped transcripts from YouTube videos for content analysis, accessibility, and data extraction.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getYoutubeTranscriptDocs,
		SampleOutput: map[string]any{
			"transcript": "This is the complete transcript of the video...",
			"segments": []map[string]any{
				{
					"text":  "This is the first segment",
					"start": 0.0,
					"end":   5.2,
				},
				{
					"text":  "This is the second segment",
					"start": 5.2,
					"end":   10.5,
				},
			},
			"metadata": map[string]any{
				"video_id": "dQw4w9WgXcQ",
				"title":    "Example Video Title",
				"duration": 180,
				"language": "en",
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetYouTubeTranscriptAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_youtube_transcript", "Get YouTube Transcript")

	form.TextField("video_url", "Video URL").
		Required(true).
		Placeholder("https://youtube.com/watch?v=dQw4w9WgXcQ").
		HelpText("The full URL of the YouTube video you want to extract the transcript from")

	schema := form.Build()

	return schema
}

func (a *GetYouTubeTranscriptAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getYouTubeTranscriptActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	accessKey := authCtx.Extra["access_key"]

	// Build query parameters
	queryParams := map[string]string{
		"url": input.VideoURL,
	}

	// Make API request
	endpoint := "/youtube/transcript"
	response, err := shared.GetSocialKitClient(accessKey, endpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *GetYouTubeTranscriptAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (a *GetYouTubeTranscriptAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"transcript": "Welcome to this comprehensive tutorial on building modern web applications. In this video, we'll explore the latest trends in web development and demonstrate practical examples. We'll start with the fundamentals and gradually move to more advanced topics.",
		"segments": []map[string]any{
			{
				"text":  "Welcome to this comprehensive tutorial on building modern web applications.",
				"start": 0.0,
				"end":   5.2,
			},
			{
				"text":  "In this video, we'll explore the latest trends in web development and demonstrate practical examples.",
				"start": 5.2,
				"end":   12.8,
			},
			{
				"text":  "We'll start with the fundamentals and gradually move to more advanced topics.",
				"start": 12.8,
				"end":   18.5,
			},
		},
		"metadata": map[string]any{
			"video_id": "dQw4w9WgXcQ",
			"title":    "Modern Web Development Tutorial",
			"duration": 1800,
			"language": "en",
			"channel":  "Tech Education Channel",
		},
	}
}

func (a *GetYouTubeTranscriptAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetYouTubeTranscriptAction() sdk.Action {
	return &GetYouTubeTranscriptAction{}
}
