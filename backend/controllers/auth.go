package controllers

import (
	"encoding/json"
	"messenger/models"
	"messenger/services"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	user, err := services.RegisterUser(u.Username, u.Password)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var in models.User
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	token, id, err := services.LoginUser(in.Username, in.Password)
	if err != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"token": token, "user_id": id})
}
