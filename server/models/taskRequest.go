package models

type TaskRequest struct {
	Command []string `json:"command"`
	Timeout int      `json:"timeout"`
}
