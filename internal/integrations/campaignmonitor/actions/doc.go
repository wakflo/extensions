package actions

import _ "embed"

//go:embed list_all_campaigns.md
var listCampaignsDocs string

//go:embed send_campaign.md
var sendCampaignDocs string

//go:embed create_campaign.md
var createCampaignDocs string

//go:embed create_campaign_from_template.md
var createCampaignTemplateDocs string

//go:embed list_subscribers.md
var listSubscribersDocs string

//go:embed get_subscriber_list.md
var getSubscriberList string

//go:embed add_subscriber.md
var addSubscriberDocs string

//go:embed get_campaign_list_and_segments.md
var getCampaignListsAndSegmentsDocs string
