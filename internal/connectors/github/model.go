package github

type Response struct {
	Data Data `json:"data"`
}

// Data contains the node information
type Data struct {
	Node Node `json:"node"`
}

// Node contains assignable users information
type Node struct {
	AssignableUsers AssignableUsers `json:"assignableUsers"`
}

// AssignableUsers contains a slice of user nodes
type AssignableUsers struct {
	Nodes []User `json:"nodes"`
}

// User represents each user with login and ID
type User struct {
	Login string `json:"login"`
	ID    string `json:"id"`
}
