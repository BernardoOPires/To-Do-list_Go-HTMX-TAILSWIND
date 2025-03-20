package handlers

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListTasks(c *gin.Context) {
	models.Mu.Lock()
	defer models.Mu.Unlock()

	tasks := make([]models.Task, len(models.Tasks))
	copy(tasks, models.Tasks)

	c.HTML(http.StatusOK, "tasks.html", gin.H{"tasks": tasks}) // o erro era q vc passou task.html que é um tarefa e n o loop de tasks q ia gerar a lista
}

func AddTask(c *gin.Context) {
	text := c.PostForm("text")
	time := c.PostForm("time")
	if text == "" {
		c.String(http.StatusBadRequest, "Texto não pode ser vazio")
		return
	}

	task := models.AddTask(text, time)

	c.HTML(http.StatusOK, "task.html", task)
}

func DelTask(c *gin.Context) {
	TaskID := c.Param("id")

	err := models.DelTask(TaskID)
	if err == false {
		c.Status(http.StatusOK)
		return
	}

}

func CompleteTask(c *gin.Context) {
	TaskID := c.Param("id")

	task := models.CompleteTask(TaskID)
	if task.ID == 0 {
		c.String(http.StatusInternalServerError, "Tarefa não encontrada")
		return
	}

	c.HTML(http.StatusOK, "task.html", task)
}

func UploadExcelHandler(c *gin.Context) {
	//passo 1 - obtenha o arquivo
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao carregar o arquivo")
		return
	}

	//passo 2 - abrir o arquivo
	openedFile, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao abrir o arquivo excel")
		return
	}
	defer openedFile.Close() //fecha o arquivo ao terminar a execução da função

	err = models.ExcelToTask(openedFile)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro processar excel para task")
		return
	}
	c.HTML(http.StatusOK, "tasks.html", gin.H{"tasks": models.Tasks})

}
