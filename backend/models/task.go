package models

import (
	"fmt"
	"sync"
)

type Task struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
}

// variavel global
var ( //use o var quando quiser delcarar variaveis globais
	Tasks  = []Task{} //variaveis que começam com maiuscula são globais e podem ser acessasa por models.Tasks
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

func DelTask(id string) bool {
	Mu.Lock()
	defer Mu.Unlock()
	var err bool = true
	for i, task := range Tasks {
		if id == fmt.Sprintf("%d", task.ID) {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
			err = false
		}
	}
	return err
}

func CompleteTask(id string) Task {
	Mu.Lock()
	defer Mu.Unlock()

	for i, task := range Tasks {
		if fmt.Sprintf("%d", task.ID) == id {
			Tasks[i].Complete = !Tasks[i].Complete //n atualizava pq vc esqueceu de colocar o c.html
			return Tasks[i]
		}
	}
	return Task{}
}
