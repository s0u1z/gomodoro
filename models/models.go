package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name        string `gorm:"not null; unique"`
	Description string
	Subtasks    []SubTask
	TimeElapsed time.Duration
}

type SubTask struct {
	gorm.Model
	TaskID      uint `gorm:"index;not null"`
	Name        string
	Description string
	Notes       string
	TimeElapsed time.Duration
}
