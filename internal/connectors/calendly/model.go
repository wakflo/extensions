package calendly

import "time"

type User struct {
	URI                 string `json:"uri"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	Email               string `json:"email"`
	SchedulingURL       string `json:"scheduling_url"`
	Timezone            string `json:"timezone"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	CurrentOrganization string `json:"current_organization"`
	AvatarURL           string `json:"avatar_url"`
}
type UsersResponse struct {
	Users []User `json:"collection"`
}

type CurrentUserResponse struct {
	Resource User `json:"resource"`
}

type EventsResponse struct {
	Events []Event `json:"collection"`
}

type Event struct {
	URI       string    `json:"uri"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
