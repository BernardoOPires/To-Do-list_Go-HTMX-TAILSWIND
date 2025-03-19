package handlers

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
	if text == "" {
		c.String(http.StatusBadRequest, "Texto não pode ser vazio")
		return
	}

	task := models.AddTask(text)

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
	}
	defer openedFile.Close() //fecha o arquivo ao terminar a execução da função

	//passo 3 - ler o arquivo(excel)
	tempFilePath := "temp.xlsx"
	err = c.SaveUploadedFile(file, tempFilePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao salvar o arquivo excel temporario")
	}

	excelFile, err := excelize.OpenFile(tempFilePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao ler os arquivo")
	}

	// so deve ter o texto na primeira coluna no arquivo
	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro no processamento das rows")
	}

	for _, row := range rows {
		if len(row) > 0 {
			lastID++
			task := Task{
				ID:       lastID,
				Text:     row[0],
				Complete: false,
			}
			tasks = append(tasks, task)
		}
	}
	c.HTML(http.StatusOK, "tasks.html", gin.H{"tasks": tasks})

}
