package main

//tecnologias
//framework gin
//htmx
//tailwing

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
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
		c.HTML(http.StatusOK, "task.html", gin.H{"tasks": tasks}) //faltou informar oq precisa no /tasks
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
				c.HTML(http.StatusOK, "partials/task.html", gin.H{"task": tasks[i], "oob": true})
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

	//teste para debug
	// for _, tmpl := range r.HTMLRender.Instance("").(*gin.Template).Templates() {
	// 	fmt.Println("Template carregado:", tmpl.Name())
	// }
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
