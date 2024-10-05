package notion

// NotionQueryResponse represents the response structure for database query
type NotionQueryResponse struct {
	Results []NotionPage `json:"results"`
}

// NotionPage represents the structure of an individual page in Notion
type NotionPage struct {
	ID         string                    `json:"id"`
	Properties map[string]NotionProperty `json:"properties"`
	URL        string                    `json:"url"`
}

// NotionProperty represents a single property of a Notion page
type NotionProperty struct {
	Title []NotionText `json:"title,omitempty"`
}

// NotionText represents a text object within a property
type NotionText struct {
	Text NotionContent `json:"text"`
}

// NotionContent represents the actual content within the text object
type NotionContent struct {
	Content string `json:"content"`
}

// NotionSearchResponse represents the response structure for a search request
type NotionSearchResponse struct {
	Results []NotionDatabase `json:"results"`
}

// NotionDatabase represents the structure of a Notion database in the search response
type NotionDatabase struct {
	ID    string       `json:"id"`
	Title []NotionText `json:"title"`
}

//// NotionText represents a text object within a title or other property
//type NotionText struct {
//	Text NotionContent `json:"text"`
//}

//// NotionContent represents the actual content within the text object
//type NotionContent struct {
//	Content string `json:"content"`
//}
