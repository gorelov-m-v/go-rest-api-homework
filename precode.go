package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getAllTasks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(tasksJSON)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if newTask.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask

	w.WriteHeader(http.StatusCreated)
}

func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id := path

	switch r.Method {
	case http.MethodGet:
		getTaskByID(w, r, id)
	case http.MethodDelete:
		deleteTaskByID(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getTaskByID(w http.ResponseWriter, _ *http.Request, id string) {
	task, exists := tasks[id]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(taskJSON)
}

func deleteTaskByID(w http.ResponseWriter, _ *http.Request, id string) {
	_, exists := tasks[id]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.WriteHeader(http.StatusOK)
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		getAllTasks(w, r)
	case http.MethodPost:
		createTask(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", handleTasks)
	mux.HandleFunc("/tasks/", handleTaskByID)

	fmt.Println("Сервер запущен на порту :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
