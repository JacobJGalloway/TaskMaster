package handler

import (
	"encoding/json"
	"net/http"
	"TaskMaster-api/service"
	"TaskMaster-api/model"
	"strconv"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task model.StudyTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTask(r.Context(), &task); err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TaskHandler) ListActiveByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	tasks, err := h.service.ListActiveTasksByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) MarkComplete(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	if err := h.service.CompleteTask(r.Context(), taskID); err != nil {
		http.Error(w, "failed to mark task complete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
