package typeform

type FormsResponse struct {
	Items []Form `json:"items"`
}

type Form struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
