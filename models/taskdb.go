package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type TaskModelDb struct {
	db *gorm.DB
}

func NewTaskModel(db *gorm.DB) *TaskModelDb {

	return &TaskModelDb{db: db}
}

func (m *TaskModelDb) AllTask() ([]*Task, error) {

	var t []*Task

	result := m.db.Find(&t)
	if result.Error != nil {
		return nil, result.Error
	}

	return t, nil

}
func (m *TaskModelDb) AddTask(name string) (*Task, error) {

	var t Task
	name = strings.TrimSpace(name)

	result := m.db.FirstOrCreate(&t, &Task{

		Name: name,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	return &t, nil
}

func (m *TaskModelDb) GetTask(name string) (*Task, error) {
	var t Task
	name = strings.TrimSpace(name)

	result := m.db.Where(&Task{Name: name}).First(&t)
	if result.Error != nil {
		return nil, result.Error
	}

	return &t, nil
}
