package models

// `Task` defines the structure of a task.
type Task struct {
	StoreID int      `json:"store_id"`
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Tags    []string `json:"tags"`
	Due     string   `json:"due"`
}
