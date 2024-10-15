package zoom

type MeetingRegistrant struct {
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name,omitempty"`
	Email                 string `json:"email"`
	Address               string `json:"address,omitempty"`
	City                  string `json:"city,omitempty"`
	State                 string `json:"state,omitempty"`
	Zip                   string `json:"zip,omitempty"`
	Country               string `json:"country,omitempty"`
	Phone                 string `json:"phone,omitempty"`
	Comments              string `json:"comments,omitempty"`
	Industry              string `json:"industry,omitempty"`
	JobTitle              string `json:"job_title,omitempty"`
	NumberOfEmployees     string `json:"no_of_employees,omitempty"`
	Organization          string `json:"org,omitempty"`
	PurchasingTimeFrame   string `json:"purchasing_time_frame,omitempty"`
	RoleInPurchaseProcess string `json:"role_in_purchase_process,omitempty"`
	Language              string `json:"language,omitempty"`
}
