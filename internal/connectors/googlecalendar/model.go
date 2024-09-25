package googlecalendar

type Calendar struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type CalendarList struct {
	Items []Calendar `json:"items"`
}

type EventList struct {
	Items []Event `json:"items"`
}

type Event struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}
