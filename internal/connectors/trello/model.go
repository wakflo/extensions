package trello

type Board struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BoardList struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard"`
	Name    string `json:"name"`
}

type Cards struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CardRequest struct {
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Pos      string   `json:"pos,omitempty"`
	IDLabels []string `json:"idLabels,omitempty"`
}

type ListRequest struct {
	Name    string `json:"name"`
	IDBoard string `json:"idBoard"`
	Pos     string `json:"pos,omitempty"`
}

type Limits struct{}

type Membership struct{}

type Preferences struct{}
