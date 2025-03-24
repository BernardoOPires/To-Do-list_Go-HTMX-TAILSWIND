package models

import (
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

type Task struct { //adicione data dps
	ID       int       `json:"id"`
	Text     string    `json:"text"`
	Complete bool      `json:"complete"`
	time     time.Time `json:"time"`
	// date
}

// variavel global
var ( //use o var quando quiser delcarar variaveis globais
	Tasks  = []Task{} //variaveis que começam com maiuscula são globais e podem ser acessasa por models.Tasks
	LastID int
	Mu     sync.Mutex
)

func AddTask(newText string, time string) Task { //começar com maiscula permite q a função seja exportada
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
	for index, task := range Tasks {
		if id == fmt.Sprintf("%d", task.ID) {
			Tasks = append(Tasks[:index], Tasks[index+1:]...)
			err = false
		}
	}
	return err
}

func CompleteTask(id string) Task {
	Mu.Lock()
	defer Mu.Unlock()

	for index, task := range Tasks {
		if fmt.Sprintf("%d", task.ID) == id {
			Tasks[index].Complete = !Tasks[index].Complete //n atualizava pq vc esqueceu de colocar o c.html
			return Tasks[index]
		}
	}
	return Task{}
}

func ExcelToTask(file multipart.File) error {
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		return err
	}

	// so deve ter o texto na primeira coluna no arquivo
	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return err
	}

	// var time string
	for _, row := range rows {
		if len(row) > 0 {
			text := row[0]
			// data = row[1]  // para dps
			if text != "" {
				LastID++
				Tasktemp := Task{ID: LastID, Text: text, Complete: false}
				Tasks = append(Tasks, Tasktemp)
			}
		}
	}

	return err
}
