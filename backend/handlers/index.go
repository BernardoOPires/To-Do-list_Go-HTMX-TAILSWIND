package handlers

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadIndex(c *gin.Context) {
	models.Mu.Lock()
	tasks := models.Tasks
	models.Mu.Unlock()

	c.HTML(http.StatusOK, "index.html", gin.H{"tasks": tasks}) // index vai renderizar a pag principal
}
