package models

import "sync"

type Task struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
}

// variavel global
var (
	Tasks  []Task
	LastID int
	Mu     sync.Mutex
)
