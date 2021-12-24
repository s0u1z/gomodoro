package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type ModelDb struct {
	db *gorm.DB
}

func NewModelDB(db *gorm.DB) *ModelDb {

	return &ModelDb{db: db}
}

func (m *ModelDb) AllTask() ([]*Task, error) {

	var t []*Task

	result := m.db.Find(&t)
	if result.Error != nil {
		return nil, result.Error
	}

	return t, nil

}

func (m *ModelDb) AllSubTask(t *Task) ([]*SubTask, error) {

	var s []*SubTask

	result := m.db.Where(&Task{Name: t.Name}).Find(&s)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil

}

func (m *ModelDb) AddSubtask(name string, task Task) (*SubTask, error) {

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

func (m *ModelDb) AddTask(name string) (*Task, error) {

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

func (m *ModelDb) GetTask(name string) (*Task, error) {
	var t Task
	name = strings.TrimSpace(name)

	result := m.db.Where(&Task{Name: name}).First(&t)
	if result.Error != nil {
		return nil, result.Error
	}

	return &t, nil
}

func (m *ModelDb) GetSubTask(name string) (*SubTask, error) {

	var s SubTask

	result := m.db.Where(&SubTask{Name: strings.TrimSpace(name)}).First(&s)
	if result.Error != nil {
		return nil, result.Error
	}
	return &s, nil

}

func (m *ModelDb) GetSubtaskByTask(t *Task) ([]*SubTask, error) {
	var s []*SubTask

	result := m.db.Model(&t).Where(&SubTask{TaskID: t.ID}).Related(&s)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil
}
