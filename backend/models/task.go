package models

import (
	"sync"
)

type Task struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
}

// variavel global
var ( //use o var quando quiser delcarar variaveis globais
	Tasks  []Task //variaveis que começam com maiuscula são globais e podem ser acessasa por models.Tasks
	LastID int
	Mu     sync.Mutex
)

func AddTask(newText string) Task { //começar com maiscula permite q a função seja exportada
	Mu.Lock()
	defer Mu.Unlock()
	LastID++
	newTask := Task{ID: LastID, Text: newText, Complete: false}
	Tasks = append(Tasks, newTask)

	return newTask
}

func DelTask(id int) {
	Mu.Lock()
	defer Mu.Unlock()
	for i, task := range Tasks {
		if task.ID == id {
			Tasks = append(Tasks[:id], Tasks[i+1:]...)
		}

	}

}

func CompleteTask(id int) {
	Mu.Lock()
	defer Mu.Unlock()
	for i, task := range Tasks {
		if task.ID == id {
			Tasks[i].Complete = !Tasks[i].Complete
		}
	}
}
