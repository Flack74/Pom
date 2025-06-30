package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Task represents a task to be completed during Pomodoro sessions
type Task struct {
	ID          string    `json:"id"`           // Unique identifier
	Title       string    `json:"title"`        // Task title
	Description string    `json:"description"`  // Optional description
	CreatedAt   time.Time `json:"created_at"`   // When the task was created
	CompletedAt time.Time `json:"completed_at"` // When the task was completed
	Sessions    int       `json:"sessions"`     // Number of sessions spent on this task
	Minutes     int       `json:"minutes"`      // Total minutes spent on this task
	Tags        []string  `json:"tags"`         // Optional tags for categorization
	IsCompleted bool      `json:"is_completed"` // Whether the task is completed
}

// TaskList represents a list of tasks
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// GetTaskFilePath returns the path to the tasks file
func GetTaskFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "tasks.json"), nil
}

// SaveTasks saves the task list to the configuration file
func SaveTasks(tasks TaskList) error {
	taskPath, err := GetTaskFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(taskPath, data, 0644)
}

// LoadTasks loads the task list from the configuration file
func LoadTasks() (TaskList, error) {
	taskPath, err := GetTaskFilePath()
	if err != nil {
		return TaskList{}, err
	}

	data, err := os.ReadFile(taskPath)
	if err != nil {
		if os.IsNotExist(err) {
			return TaskList{Tasks: []Task{}}, nil
		}
		return TaskList{}, err
	}

	var tasks TaskList
	if err := json.Unmarshal(data, &tasks); err != nil {
		return TaskList{}, err
	}

	return tasks, nil
}

// AddTask adds a new task to the list
func AddTask(title, description string, tags []string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	newTask := Task{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		Tags:        tags,
		IsCompleted: false,
	}

	tasks.Tasks = append(tasks.Tasks, newTask)
	return SaveTasks(tasks)
}

// CompleteTask marks a task as completed
func CompleteTask(id string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	for i := range tasks.Tasks {
		if tasks.Tasks[i].ID == id {
			tasks.Tasks[i].IsCompleted = true
			tasks.Tasks[i].CompletedAt = time.Now()
			return SaveTasks(tasks)
		}
	}

	return fmt.Errorf("task with ID %s not found", id)
}

// UpdateTaskProgress updates the time spent on a task
func UpdateTaskProgress(id string, sessions, minutes int) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	for i := range tasks.Tasks {
		if tasks.Tasks[i].ID == id {
			tasks.Tasks[i].Sessions += sessions
			tasks.Tasks[i].Minutes += minutes
			return SaveTasks(tasks)
		}
	}

	return fmt.Errorf("task with ID %s not found", id)
}

// ListTasks displays all tasks
func ListTasks(showCompleted bool) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	fmt.Println("\nTasks:")
	for _, task := range tasks.Tasks {
		if !showCompleted && task.IsCompleted {
			continue
		}

		status := "[ ]"
		if task.IsCompleted {
			status = "[✓]"
		}

		fmt.Printf("%s %s (ID: %s)\n", status, task.Title, task.ID)
		if task.Description != "" {
			fmt.Printf("   Description: %s\n", task.Description)
		}
		if len(task.Tags) > 0 {
			fmt.Printf("   Tags: %v\n", task.Tags)
		}
		fmt.Printf("   Sessions: %d, Total Time: %d minutes\n", task.Sessions, task.Minutes)
		fmt.Println()
	}

	return nil
}
