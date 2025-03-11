package main

//tecnologias
//framework gin
//htmx
//tailwing

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type Task struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
}

var tasks = []Task{}
var lastID = 0
var mu sync.Mutex // verifique oq esse faz em detalhes

func main() {
	r := gin.Default()

	// Carregar todos os templates da pasta "template/"
	templateFiles, err := loadTemplates("template/")
	if err != nil {
		panic(err)
	}

	fmt.Println("Templates encontrados:", templateFiles)

	// Criar um novo template do Go
	templ := template.New("")

	// Adicionar arquivos ao template
	templ, err = templ.ParseFiles(templateFiles...)
	if err != nil {
		panic(err)
	}

	// Definir o template no Gin
	r.SetHTMLTemplate(templ)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"tasks": tasks}) // index vai renderizar a pag principal
	})

	r.GET("/tasks", func(c *gin.Context) {
		for _, task := range tasks {
			fmt.Println("Tarefa:", task.ID, task.Text) // 🔍 Verifique se o ID está correto
		} // a lista esta salvando
		c.HTML(http.StatusOK, "tasks.html", gin.H{"tasks": tasks}) //estou passando uma lista para tasks q carrega o text de uma so
		//passe apenas uma task passe isso tudo pelo html
	})

	//adicionar tarefa
	r.POST("/add", func(c *gin.Context) {
		text := c.PostForm("text")
		if text != "" {
			mu.Lock()
			lastID++
			task := Task{ID: lastID, Text: text, Complete: false}
			tasks = append(tasks, task) //adiciona a lista, mas continua dando erros
			mu.Unlock()
			// c.HTML(http.StatusOK, "task.html", gin.H{"task": task})
			c.HTML(http.StatusOK, "task.html", gin.H{"ID": task.ID, "Text": task.Text, "Complete": task.Complete})
		}
	})

	//marcar/desmarcar como concluido
	r.PATCH("/complete/:id", func(c *gin.Context) {
		id := c.Param("id")
		mu.Lock()
		defer mu.Unlock()
		for i := range tasks {
			if id == fmt.Sprintf("%d", tasks[i].ID) {
				tasks[i].Complete = !tasks[i].Complete //Alterna entre concluído/não concluído
				// Renderiza o HTML atualizado da tarefa para substituir no front-end via HTMX
				c.HTML(http.StatusOK, "task.html", gin.H{"task": tasks[i], "oob": true})
				return
			}
		}
		c.String(http.StatusNotFound, "Tarefa não encontrada")
	})

	r.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
		mu.Lock()
		defer mu.Unlock()
		for i, task := range tasks {
			if id == fmt.Sprintf("%d", task.ID) {
				tasks = append(tasks[:i], tasks[i+1:]...) // junta as partes, remove o i por juntar as parte que vem asntescom a parte q vem dps
				c.String(http.StatusOK, "Tarefa removida")
				break
			}
		}
		c.String(http.StatusNotFound, "Tarefa não encontrada")
	})

	//somnete o texto é necessario no csv, e deve estar na primeira coluna
	r.POST("/upload", func(c *gin.Context) {
		//passo 1 - obtenha o arquivo
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusInternalServerError, "Erro ao carregar o arquivo")
			return
		}

		//passo 2 - abrir o arquivo
		openedFile, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "Erro ao abrir o arquivo")
		}
		defer openedFile.Close() //fecha o arquivo ao terminar a execução da função

		//passo 3 - ler o arquivo(csv)
		reader := csv.NewReader(openedFile)
		records, err := reader.ReadAll() //slice de cada linha do arquivo
		if err != nil {
			c.String(http.StatusInternalServerError, "Erro ao ler o arquivo(csv)")
		}

		//passo 4 - converter os registros para tarefas
		// bloqueia o acesso a variavel tasks durante o processo,
		// fazendo com que somente uma go routine possa modificar a tasks por vez
		mu.Lock()
		defer mu.Unlock()
		for _, record := range records {
			lastID++
			task := Task{
				ID:       lastID,
				Text:     record[0],
				Complete: false,
			}
			tasks = append(tasks, task)
		}

		//passo 5 - retornar a lista atualizada
		c.HTML(http.StatusOK, "task.html", gin.H{"task": tasks})
	})

	r.POST("/upload-excel", func(c *gin.Context) {
		print("entrada")
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
	})
	r.Run(":8080")
}
func loadTemplates(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if path != root {
				loadTemplates(path)
			}
		} else {
			files = append(files, path)
		}
		return err
	})
	return files, err
}
