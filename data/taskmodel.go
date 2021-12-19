package data

import "log"

type Task struct {
	Name    string     `json:"tasks_name"`
	Details string     `json:"tasks_deatils"`
	Subtask []*SubTask `json:"task_sub"`
}

type SubTask struct {
	Name    string `json:"subtask_name"`
	Details string `json:"subtask_details"`
	Notes   string `json:"subtask_notes"`
}

type Tasks []*Task
type SubTasks []*SubTask

// GetSubtask using TaskName
func GetSubTaskByName(task Task) (SubTasks, error) {

	tasks, err := GetAllTask()
	if err != nil {
		return nil, err
	}

	var subtask SubTasks

	for _, t := range tasks {
		if t.Name == task.Name {

			subtask = append(subtask, t.Subtask...)
		}
	}

	return subtask, nil

}

//Add SubTask to Task specified
func AddSubTask(sub SubTask, task Task) error {

	filedata, err := ReadFile()

	if err != nil {
		return err
	}

	tasks, err := FromJson(filedata)
	if err != nil {

		return err
	}

	for _, t := range tasks {
		if t.Name == task.Name {

			pos := findsubPos(t.Subtask)
			if pos != -1 {

				t.Subtask = append(t.Subtask[:pos], t.Subtask[pos+1:]...)
			}
			t.Subtask = append(t.Subtask, &sub)

		}
	}

	data, err := ToJson(tasks)
	if err != nil {

		return err
	}

	return WriteFile(data)
}

//Add Task to Json
func AddTaskToList(task Task) error {

	// get tasks list from the file

	filedata, err := ReadFile()

	if err != nil {
		return err
	}

	tasks, err := FromJson(filedata)
	if err != nil {

		return err
	}

	pos := findtaskPos(tasks)
	if pos != -1 {

		tasks = append(tasks[:pos], tasks[pos+1:]...)

	}

	tasks = append(tasks, &task)

	data, err := ToJson(tasks)
	if err != nil {

		return err
	}

	return WriteFile(data)

}

//Function to GetAll task

func GetAllTask() (Tasks, error) {

	filedata, err := ReadFile()

	if err != nil {
		return nil, err
	}

	tasks, err := FromJson(filedata)
	if err != nil {

		return nil, err
	}

	return tasks, nil

}

func findsubPos(subs SubTasks) int {

	for idx, sub := range subs {

		if sub.Name == "" {
			return idx
		}
	}

	return -1
}

func findtaskPos(tasks Tasks) int {

	for idx, task := range tasks {

		if task.Name == "" {

			return idx
		}
	}
	return -1
}

//Function to create new Task and SubTask

func CreateTask(name string) *Task {

	return &Task{
		Name: name,
		Subtask: []*SubTask{
			{
				Name: "",
			},
		},
	}

}

func CreateSubTask(name string) *SubTask {
	return &SubTask{
		Name: name,
	}
}

// Function to init TasksLists

func Init() *Tasks {

	task := &Tasks{
		&Task{
			Subtask: []*SubTask{
				{},
			},
		},
	}

	data, err := ToJson(task)
	if err != nil {
		log.Panic(err)
	}

	err = WriteFile(data)
	if err != nil {
		log.Panic(err)
	}

	return task
}
