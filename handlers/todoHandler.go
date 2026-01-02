package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yogisyo16/root-aura-service/services"
)

type TodoHandler struct {
	Service        services.Todo
	DetailsService services.TodoDetails // Add this
}

// Create Todo request structure
type CreateTodoRequest struct {
	Task      string `json:"task"`
	DateStart string `json:"date_start"`
	DateDue   string `json:"date_due"`
	Completed bool   `json:"completed"`
}

// Update Todo request structure
type UpdateTodoRequest struct {
	Task      string `json:"task"`
	DateStart string `json:"date_start"`
	DateDue   string `json:"date_due"`
	Completed bool   `json:"completed"`
}

type TodoWithDetails struct {
	ID          string                `json:"id,omitempty"`
	UserID      string                `json:"user_id"`
	Task        string                `json:"task"`
	DateStart   time.Time             `json:"date_start,omitempty"`
	DateDue     time.Time             `json:"date_due,omitempty"`
	Completed   bool                  `json:"completed"`
	TodoDetails *services.TodoDetails `json:"todo_details"`
	CreatedAt   time.Time             `json:"created_at,omitempty"`
	UpdatedAt   time.Time             `json:"update_at,omitempty"`
}

// Generic response structure
func NewTodoHandler(service services.Todo, detailsService services.TodoDetails) *TodoHandler {
	return &TodoHandler{
		Service:        service,
		DetailsService: detailsService, // Initialize this
	}
}

// Health Check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Msg:  "Health Check",
		Code: 200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

// Logic to get all todos
// Updated with details todos included
func (h *TodoHandler) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.Service.GetAllTodos()
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if todos == nil {
		todos = []services.Todo{}
	}

	// Fetch details for each todo
	var todosWithDetails []TodoWithDetails
	for _, todo := range todos {
		todoWithDetails := TodoWithDetails{
			ID:          todo.ID,
			UserID:      todo.UserID,
			Task:        todo.Task,
			DateStart:   todo.DateStart,
			DateDue:     todo.DateDue,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
			TodoDetails: nil, // Default to nil (will show as null in JSON)
		}

		// Try to get the details for this todo
		details, err := h.DetailsService.GetTodoDetailsByTodoId(todo.ID)
		if err == nil && details.ID != "" {
			todoWithDetails.TodoDetails = &details
		}

		todosWithDetails = append(todosWithDetails, todoWithDetails)
	}

	// Response structure wrapper
	response := struct {
		Code int `json:"code"`
		Data struct {
			Items []TodoWithDetails `json:"items"`
		} `json:"data"`
	}{
		Code: 200,
		Data: struct {
			Items []TodoWithDetails `json:"items"`
		}{
			Items: todosWithDetails,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

// Logic to get todo by id
// Updated with details todos included sorted by id
func (h *TodoHandler) getTodoByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	todo, err := h.Service.GetTodoById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	// Create response with details
	todoWithDetails := TodoWithDetails{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Task:        todo.Task,
		DateStart:   todo.DateStart,
		DateDue:     todo.DateDue,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
		TodoDetails: nil, // Default to nil
	}

	// Try to get the details for this todo
	details, err := h.DetailsService.GetTodoDetailsByTodoId(todo.ID)
	if err == nil && details.ID != "" {
		todoWithDetails.TodoDetails = &details
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todoWithDetails)
}

// Create todo
func parseDateTime(dateStr string) (time.Time, error) {
	// Try parsing as full RFC3339 first (2006-01-02T15:04:05Z)
	t, err := time.Parse(time.RFC3339, dateStr)
	if err == nil {
		return t, nil
	}

	// Try parsing as date only (2006-01-02)
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		// If only date is provided, set time to 00:00:00
		return t, nil
	}

	// Try parsing as datetime without timezone (2006-01-02T15:04:05)
	t, err = time.Parse("2006-01-02T15:04:05", dateStr)
	if err == nil {
		return t, nil
	}

	return time.Time{}, err
}

// Create todo
func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid request body",
			Code: 400,
		})
		return
	}

	// Parse dates
	dateStart, err := parseDateTime(req.DateStart)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid date_start format. Use YYYY-MM-DD or ISO 8601 format",
			Code: 400,
		})
		return
	}

	dateDue, err := parseDateTime(req.DateDue)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid date_due format. Use YYYY-MM-DD or ISO 8601 format",
			Code: 400,
		})
		return
	}

	// IMPORTANT: Validate date logic - start must be before or equal to due
	if dateStart.After(dateDue) {
		log.Println("Start date is after due date")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Start date cannot be after due date",
			Code: 400,
		})
		return
	}

	// Additional validation: check if task is not empty
	if len(req.Task) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Task name is required",
			Code: 400,
		})
		return
	}

	// Create the Todo
	newTodo := services.Todo{
		Task:      req.Task,
		DateStart: dateStart,
		DateDue:   dateDue,
		Completed: req.Completed,
	}

	err = h.Service.InsertTodo(newTodo)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Failed to create todo",
			Code: 500,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(Response{
		Msg:  "Successfully Created Todo",
		Code: 201,
	})
}

// Also update the updateTodo function with the same validation
func (h *TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid request body",
			Code: 400,
		})
		return
	}

	dateStart, err := parseDateTime(req.DateStart)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid date_start format",
			Code: 400,
		})
		return
	}

	dateDue, err := parseDateTime(req.DateDue)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Invalid date_due format",
			Code: 400,
		})
		return
	}

	// Validate date logic
	if dateStart.After(dateDue) {
		log.Println("Start date is after due date")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Start date cannot be after due date",
			Code: 400,
		})
		return
	}

	// Validate task
	if len(req.Task) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Task name is required",
			Code: 400,
		})
		return
	}

	updateTodo := services.Todo{
		Task:      req.Task,
		DateStart: dateStart,
		DateDue:   dateDue,
		Completed: req.Completed,
	}

	_, err = h.Service.UpdatedTodo(id, updateTodo)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response{
			Msg:  err.Error(),
			Code: 500,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Response{
		Msg:  "Successfully Updated Todo",
		Code: 200,
	})
}

// Delete Todo it can be used only by id (On going to make delete all todos)
func (h *TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Service.DeleteTodo(id)
	if err != nil {
		errorRes := Response{
			Msg:  "Error",
			Code: 304,
		}
		json.NewEncoder(w).Encode(errorRes)
		w.WriteHeader(errorRes.Code)
		return
	}

	res := Response{
		Msg:  "Succesfully Deleted Todo",
		Code: 200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}
