package models

type Calendar struct {
	Year          int
	Month         string
	Days          []int
	EmptyDays     []int
	PreviousMonth string
	NextMonth     string
}
