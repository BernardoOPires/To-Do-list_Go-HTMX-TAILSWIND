package handlers

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var tasks []models.Task
var lastID int

func ListTasks(c *gin.Context) {
	models.Mu.Lock()
	defer models.Mu.Unlock()

	tasks := make([]models.Task, len(models.Tasks))
	copy(tasks, models.Tasks)

	c.HTML(http.StatusOK, "tasks.html", gin.H{"tasks": tasks}) // o erro era q vc passou task.html que é um tarefa e n o loop de tasks q ia gerar a lista
}

func AddTask(c *gin.Context) {
	text := c.PostForm("text")
	if text == "" {
		c.String(http.StatusBadRequest, "Texto não pode ser vazio")
		return
	}

	task := models.AddTask(text)

	c.HTML(http.StatusOK, "task.html", task)
}

func DelTask(c *gin.Context) {
	TaskID := c.Param("id")

	models.DelTask(TaskID)

}

func CompleteTask(c *gin.Context) {
	TaskID := c.Param("id")

	task := models.CompleteTask(TaskID)
	if task.ID == 0 {
		c.String(http.StatusInternalServerError, "Tarefa não encontrada")
	}

	c.HTML(http.StatusOK, "task.html", task)

}
