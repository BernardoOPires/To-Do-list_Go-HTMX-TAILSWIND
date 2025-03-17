package handlers

import (
	"backend/models"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
)

func CalendarHandler(c *gin.Context) {
	//c.query("month") serve para obter os valores dos parametors da url
	monthQuery := c.Query("month") //pesquise mais sobre o .query
	var currentTime time.Time
	if monthQuery != "" {
		parsedTime, err := time.Parse("2006-01", monthQuery) // time.Parse("2006-01", monthQuery) serve para formatar da forma que o go usa nos parsing dates
		if err == nil {
			currentTime = parsedTime
		} else {
			currentTime = time.Now()
		}
	} else {
		currentTime = time.Now()
	}

	year, month, _ := currentTime.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentTime.Location())
	emptyDays := int(firstOfMonth.Weekday())

	data := models.Calendar{
		Year:          year,
		Month:         month.String(),
		Days:          daysInMonthSlice(daysInMonth(year, month)),
		EmptyDays:     make([]int, emptyDays),
		PreviousMonth: firstOfMonth.AddDate(0, -1, 0).Format("2006-01"),
		NextMonth:     firstOfMonth.AddDate(0, 1, 0).Format("2006-01"),
	}

	tmpl := template.Must(template.ParseFiles("templates/calendario.html"))
	tmpl.Execute(c.Writer, data)
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day() //pq month+1, 0, 0, 0, 0, 0, ??
}

func daysInMonthSlice(days int) []int {
	slice := make([]int, days) //oq make faz?
	for i := 1; i <= days; i++ {
		slice[i-1] = i
	}
	return slice
}
