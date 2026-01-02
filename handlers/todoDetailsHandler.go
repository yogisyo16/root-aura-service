package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yogisyo16/root-aura-service/services"
)

type TodoDetailsHandler struct {
	Service services.TodoDetails
}

func NewTodoDetailsHandler(service services.TodoDetails) *TodoDetailsHandler {
	return &TodoDetailsHandler{
		Service: service,
	}
}

type CreateTodoDetailsRequest struct {
	TodoID          string `json:"todo_id"`
	TaskDetails     string `json:"task_details"`
	NotesDetails    string `json:"notes_details"`
	StatusDetails   string `json:"status_details"`
	PriorityDetails string `json:"priority_details"`
}

func (h *TodoDetailsHandler) getTodoDetails(w http.ResponseWriter, r *http.Request) {
	todoDetails, err := h.Service.GetAllTodosDetails()
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	response := struct {
		Code int `json:"code"`
		Data struct {
			Items []services.TodoDetails `json:"items"`
		} `json:"data"`
	}{
		Code: 200,
		Data: struct {
			Items []services.TodoDetails `json:"items"`
		}{
			Items: todoDetails,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

func (h *TodoDetailsHandler) getTodoDetailsByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	todoDetail, err := h.Service.GetTodoDetailsById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todoDetail)
}

func (h *TodoDetailsHandler) createTodoDetails(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoDetailsRequest

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

	newTodoDetails := services.TodoDetails{
		TodoID:          req.TodoID,
		TaskDetails:     req.TaskDetails,
		NotesDetails:    req.NotesDetails,
		StatusDetails:   req.StatusDetails,
		PriorityDetails: req.PriorityDetails,
	}

	err = h.Service.InsertTodoDetails(newTodoDetails)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response{
			Msg:  "Failed to create todo details",
			Code: 500,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(Response{
		Msg:  "Successfully Created Todo Details",
		Code: 201,
	})
}

func (h *TodoDetailsHandler) deleteTodoDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Service.DeleteTodoDetails(id)
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
