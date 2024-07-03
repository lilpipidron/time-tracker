package requests

type EndTaskRequest struct {
	UserID int `json:"user_id"`
	TaskID int `json:"task_id"`
}
