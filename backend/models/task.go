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
	DateTime time.Time `json:"datetime"`
	Priority string    `json:"priority"`
}

// variavel global
var ( //use o var quando quiser delcarar variaveis globais
	Tasks  = []Task{} //variaveis que começam com maiuscula são globais e podem ser acessasa por models.Tasks
	LastID int
	Mu     sync.Mutex
)

func AddTask(newText string, taskTime string, date string, priority string) (Task, error) {
	Mu.Lock()
	defer Mu.Unlock()
	//Revise o modo como esta sendo feito para deixar mais eficiente apos terminar a proxima função a ser adicionada
	year, month, day := 0, time.January, 1
	hour, minute := 0, 0
	var err error

	if date != "" {
		parsedDate, parseErr := time.Parse("02/01/2006", date)
		if parseErr != nil {
			err = parseErr
		} else {
			year, month, day = parsedDate.Date()
		}
	}

	if taskTime != "" {
		parsedTime, parseErr := time.Parse("15:04", taskTime)
		if parseErr != nil {
			err = parseErr
		} else {
			hour, minute, _ = parsedTime.Clock()
		}
	}

	finalDateTime := time.Date(year, month, day, hour, minute, 0, 0, time.Local)

	LastID++
	newTask := Task{
		ID:       LastID,
		Text:     newText,
		DateTime: finalDateTime,
		Priority: priority,
		Complete: false,
	}
	Tasks = append(Tasks, newTask)

	return newTask, err
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

func GetDate(date time.Time) string {
	return date.Format("02/01/2006")
}

func GetTime(time time.Time) string {
	return time.Format("15:04")
}
