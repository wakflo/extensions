package shared

type NotionQueryResponse struct {
	Results []NotionPage `json:"results"`
}

type NotionPage struct {
	ID         string                    `json:"id"`
	Properties map[string]NotionProperty `json:"properties"`
	URL        string                    `json:"url"`
}

type NotionProperty struct {
	Title []NotionText `json:"title,omitempty"`
}

type NotionText struct {
	Text NotionContent `json:"text"`
}

type NotionContent struct {
	Content string `json:"content"`
}

type NotionSearchResponse struct {
	Results []NotionDatabase `json:"results"`
}

type NotionDatabase struct {
	ID    string       `json:"id"`
	Title []NotionText `json:"title"`
}
