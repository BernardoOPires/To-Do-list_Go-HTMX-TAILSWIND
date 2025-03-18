package main

//tecnologias
//framework gin
//htmx
//tailwing

import (
	"backend/handlers"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Carregar todos os templates da pasta "template/"
	templateFiles, err := loadTemplates("../frontend/template/")
	if err != nil {
		panic(err)
	}

	fmt.Println("Templates encontrados:", templateFiles) //templates para debug

	// Criar um novo template do Go
	templ := template.New("")

	// Adicionar arquivos ao template
	templ, err = templ.ParseFiles(templateFiles...)
	if err != nil {
		panic(err)
	}

	// Definir o template no Gin
	r.SetHTMLTemplate(templ)

	//tarefas e funções a adicionar
	//pesquisar por tarefa
	//filtrar tarefas
	//deletar todas as tarefas, com confirmação
	//marcar mais de uma por vez
	//agendar e aceitar arquivo
	//receber arquivo xtls e carregar as tarefas dele

	r.GET("/", handlers.LoadIndex)
	r.GET("/getTasks", handlers.ListTasks)
	r.POST("/add", handlers.AddTask) //post pq muda completamente o valor, add
	r.DELETE("/delete/:id", handlers.DelTask)
	r.PATCH("/complete/:id", handlers.CompleteTask) //path pq muda um valor de uma variavel, mudança parcial

	r.GET("/calendario", handlers.CalendarHandler)

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
