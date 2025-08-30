package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	if u.Username == "" || u.Password == "" {
		http.Error(w, "username/password required", 400)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	err := DB.QueryRow(Ctx,
		"insert into users (username, password_hash) values ($1,$2) returning id",
		u.Username, string(hash)).Scan(&u.ID)
	if err != nil {
		http.Error(w, "user exists?", 400)
		return
	}
	u.Password = ""
	writeJSON(w, 201, u)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var in User
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	var id int64
	var hash string
	err := DB.QueryRow(Ctx, "select id, password_hash from users where username=$1", in.Username).Scan(&id, &hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Password)) != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"usr": in.Username,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte(mustEnv("JWT_SECRET")))
	writeJSON(w, 200, map[string]any{"token": signed, "user_id": id})
}
