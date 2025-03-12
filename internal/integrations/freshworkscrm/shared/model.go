package shared

type Contact struct {
	ID              int     `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	DisplayName     string  `json:"display_name"`
	Avatar          *string `json:"avatar"`
	JobTitle        *string `json:"job_title"`
	City            *string `json:"city"`
	State           *string `json:"state"`
	Zipcode         *string `json:"zipcode"`
	Country         *string `json:"country"`
	Email           *string `json:"email"`
	TimeZone        *string `json:"time_zone"`
	WorkNumber      *string `json:"work_number"`
	MobileNumber    string  `json:"mobile_number"`
	Address         *string `json:"address"`
	LastSeen        *string `json:"last_seen"`
	LeadScore       int     `json:"lead_score"`
	LastContacted   *string `json:"last_contacted"`
	OpenDealsAmount string  `json:"open_deals_amount"`
	Links           struct {
		Conversations string `json:"conversations"`
		Activities    string `json:"activities"`
	} `json:"links"`
	CustomField struct {
		CfIsActive bool `json:"cf_is_active"`
	} `json:"custom_field"`
	UpdatedAt string  `json:"updated_at"`
	Keyword   *string `json:"keyword"`
	Medium    *string `json:"medium"`
	Facebook  *string `json:"facebook"`
	Twitter   *string `json:"twitter"`
	Linkedin  *string `json:"linkedin"`
}

type ContactWrapper struct {
	Contacts []Contact `json:"contacts"`
}

type FilterWrapper struct {
	Filters []ViewDetails `json:"filters"`
}

type ViewDetails struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ModelClassName string `json:"model_class_name"`
	UserID         int    `json:"user_id"`
	IsDefault      bool   `json:"is_default"`
	IsPublic       bool   `json:"is_public"`
	UpdatedAt      string `json:"updated_at"`
}
