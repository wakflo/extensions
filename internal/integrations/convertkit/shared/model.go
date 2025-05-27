package shared

type SubscribersResponse struct {
	TotalSubscribers int `json:"total_subscribers"`
	Page             int `json:"page"`
	TotalPages       int `json:"total_pages"`
	Subscribers      []struct {
		ID           int                    `json:"id"`
		FirstName    string                 `json:"first_name"`
		EmailAddress string                 `json:"email_address"`
		State        string                 `json:"state"`
		CreatedAt    string                 `json:"created_at"`
		Fields       map[string]interface{} `json:"fields"`
	} `json:"subscribers"`
}
