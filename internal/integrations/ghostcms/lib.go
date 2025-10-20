package ghostcms

import (
	_ "embed"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/ghostcms/actions"
	// "github.com/wakflo/extensions/internal/integrations/ghostcms/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var (
	form = smartform.NewAuthForm("ghost-auth", "Ghost CMS API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("site_url", "Site URL").
		Required(true).
		Placeholder("https://yourblog.ghost.io").
		HelpText("Your Ghost site URL (e.g., https://yourblog.ghost.io)")

	_ = form.TextField("admin_api_key", "Admin API Key").
		Required(true).
		Placeholder("site-name:key").
		HelpText("Your Ghost Admin API key (found in Settings → Integrations)")

	GhostSharedAuth = form.Build()
)

var Integration = sdk.Register(NewGhostCMS())

type GhostCMS struct{}

func (g *GhostCMS) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (g *GhostCMS) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   GhostSharedAuth,
	}
}

func (g *GhostCMS) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		// triggers.NewNewPostTrigger(),
		// triggers.NewPostUpdatedTrigger(),
		// triggers.NewPostPublishedTrigger(),
		// triggers.NewNewMemberTrigger(),
		// triggers.NewMemberUpdatedTrigger(),
		// triggers.NewNewPageTrigger(),
		// triggers.NewPageUpdatedTrigger(),
		// triggers.NewNewTagTrigger(),
	}
}

func (g *GhostCMS) Actions() []sdk.Action {
	return []sdk.Action{
		// Posts
		actions.NewCreatePostAction(),
		// actions.NewUpdatePostAction(),
		// actions.NewDeletePostAction(),
		// actions.NewGetPostAction(),
		actions.NewListPostsAction(),

		// // Pages
		// actions.NewCreatePageAction(),
		// actions.NewUpdatePageAction(),
		// actions.NewDeletePageAction(),
		// actions.NewGetPageAction(),
		// actions.NewListPagesAction(),

		// // Members
		// actions.NewCreateMemberAction(),
		// actions.NewUpdateMemberAction(),
		// actions.NewDeleteMemberAction(),
		// actions.NewGetMemberAction(),
		// actions.NewListMembersAction(),

		// // Tags
		// actions.NewCreateTagAction(),
		// actions.NewUpdateTagAction(),
		// actions.NewDeleteTagAction(),
		// actions.NewGetTagAction(),
		// actions.NewListTagsAction(),

		// // Media & Site
		// actions.NewUploadImageAction(),
		// actions.NewGetSiteAction(),

		// // Authors
		// actions.NewListAuthorsAction(),
		// actions.NewGetAuthorAction(),
	}
}

func NewGhostCMS() sdk.Integration {
	return &GhostCMS{}
}
