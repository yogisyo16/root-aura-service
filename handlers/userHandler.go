package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-mongo-todos/services"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Service services.User
}

func NewUserHandler(service services.User) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) insertUser(w http.ResponseWriter, r *http.Request) {
	var newUser services.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Println(err)
		errorRes := Response{
			Msg:  "Invalid request body",
			Code: 400,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		errorRes := Response{
			Msg:  "Error processing password",
			Code: 500,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}
	newUser.Password = string(hashedPassword)

	err = h.Service.InsertUser(newUser)
	if err != nil {
		errorRes := Response{
			Msg:  "Error creating user",
			Code: 500,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorRes.Code)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	res := Response{
		Msg:  "Successfully Created User",
		Code: 201,
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

func (h *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		log.Println(err)
		return
	}

	response := struct {
		Code int `json:"code"`
		Data struct {
			Items []services.User `json:"items"`
		} `json:"data"`
	}{
		Code: 200,
		Data: struct {
			Items []services.User `json:"items"`
		}{
			Items: users,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		log.Println(err)
		return
	}

	response := struct {
		Code int `json:"code"`
		Data struct {
			Items services.User `json:"items"`
		} `json:"data"`
	}{
		Code: 200,
		Data: struct {
			Items services.User `json:"items"`
		}{
			Items: user,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}
