package handlers

import (
	"backend/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var tasks []models.Task
var lastID int

func ListTasks(c *gin.Context) {
	models.Mu.Lock()
	tasks := models.Tasks
	models.Mu.Unlock()

	for i, task := range tasks {
		fmt.Printf("Task %d: ID = %v\n", i, task.ID)
	}

	c.HTML(http.StatusOK, "task.html", gin.H{"tasks": tasks})
}

func AddTask(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		c.String(http.StatusBadRequest, "Texto não pode ser vazio")
		return
	}

	task := models.AddTask(text)

	// models.Mu.Lock()
	// lastID++
	// task := models.Task{ID: lastID, Text: text, Complete: false}
	// tasks = append(tasks, task)
	// models.Mu.Unlock()

	c.HTML(http.StatusOK, "task.html", task)
}

func DelTask(c *gin.Context) {
	TaskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "id invalido")
		return
	}

	models.DelTask(TaskID)

	// for i, task := range tasks {
	// 	if task.ID == TaskID {
	// 		tasks = append(tasks[:i], tasks[i+1:]...)
	// 		c.String(http.StatusOK, "Tarefa removida") // mude o task.html para o caminho correto
	// 		return
	// 	}
	// }

	c.String(http.StatusInternalServerError, "Tarefa não encontrada")
}

func CompleteTask(c *gin.Context) {
	TaskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "id invalido")
		return
	}

	models.CompleteTask(TaskID)

	// for i, task := range tasks {
	// 	if task.ID == TaskID {
	// 		tasks[i].Complete = !tasks[i].Complete
	// 		c.HTML(http.StatusOK, "task.html", tasks[i]) // mude o task.html para o caminho correto
	// 		return
	// 	}
	// }

	c.String(http.StatusInternalServerError, "Tarefa não encontrada")
}
