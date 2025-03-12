package shared

import "time"

type ContactRequest struct {
	Properties map[string]interface{} `json:"properties"`
}

type ListResponse struct {
	Lists []List `json:"lists"`
}

type List struct {
	ProcessingType   string       `json:"processingType"`
	ObjectTypeID     string       `json:"objectTypeId"`
	UpdatedByID      string       `json:"updatedById"`
	FiltersUpdatedAt time.Time    `json:"filtersUpdatedAt"`
	ListID           string       `json:"listId"`
	CreatedAt        time.Time    `json:"createdAt"`
	ProcessingStatus string       `json:"processingStatus"`
	DeletedAt        time.Time    `json:"deletedAt"`
	ListVersion      int          `json:"listVersion"`
	Size             int          `json:"size"`
	Name             string       `json:"name"`
	CreatedByID      string       `json:"createdById"`
	FilterBranch     FilterBranch `json:"filterBranch"`
	UpdatedAt        time.Time    `json:"updatedAt"`
}

type FilterBranch struct {
	FilterBranchType     string         `json:"filterBranchType"`
	FilterBranches       []FilterBranch `json:"filterBranches"`
	FilterBranchOperator string         `json:"filterBranchOperator"`
	Filters              []interface{}  `json:"filters"`
}

type DealPipelineResponse struct {
	Results []PipelineResult `json:"results"`
}

type PipelineResult struct {
	CreatedAt    time.Time `json:"createdAt"`
	ArchivedAt   time.Time `json:"archivedAt"`
	Archived     bool      `json:"archived"`
	DisplayOrder int       `json:"displayOrder"`
	Stages       []Stage   `json:"stages"`
	Label        string    `json:"label"`
	ID           string    `json:"id"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Stage struct {
	CreatedAt        time.Time     `json:"createdAt"`
	ArchivedAt       time.Time     `json:"archivedAt"`
	Archived         bool          `json:"archived"`
	Metadata         StageMetadata `json:"metadata"`
	DisplayOrder     int           `json:"displayOrder"`
	WritePermissions string        `json:"writePermissions"`
	Label            string        `json:"label"`
	ID               string        `json:"id"`
	UpdatedAt        time.Time     `json:"updatedAt"`
}

type StageMetadata struct {
	TicketState string `json:"ticketState"`
}
