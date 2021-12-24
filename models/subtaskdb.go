package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type SubTaskModelDb struct {
	db *gorm.DB
}

func NewSubTaskModel(db *gorm.DB) *SubTaskModelDb {

	return &SubTaskModelDb{db: db}
}

func (m *SubTaskModelDb) AllSubTask(t *Task) ([]*SubTask, error) {

	var s []*SubTask

	result := m.db.Where(&Task{Name: t.Name}).Find(&s)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil

}

func (m *SubTaskModelDb) AddSubtask(name string, task Task) (*SubTask, error) {

	subtask := SubTask{

		Name:   strings.TrimSpace(name),
		TaskID: task.ID,
	}
	result := m.db.Create(&subtask)
	if result.Error != nil {
		return nil, result.Error
	}

	return &subtask, nil

}

func (m *SubTaskModelDb) GetSubTask(name string) (*SubTask, error) {

	var s SubTask

	result := m.db.Where(&SubTask{Name: strings.TrimSpace(name)}).First(&s)
	if result.Error != nil {
		return nil, result.Error
	}
	return &s, nil

}

func (m *SubTaskModelDb) GetSubtaskByTask(t *Task) ([]*SubTask, error) {
	var s []*SubTask

	result := m.db.Model(&t).Where(&SubTask{TaskID: t.ID}).Related(&s)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil
}
