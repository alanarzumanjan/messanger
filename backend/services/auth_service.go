package services

import (
	"errors"
	"messenger/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, password string) (models.User, error) {
	if username == "" || len(password) < 6 {
		return models.User{}, errors.New("invalid username/password")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	var u models.User
	err := DB.QueryRow(Ctx,
		"insert into users (username, password_hash) values ($1,$2) returning id, created_at",
		username, string(hash),
	).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return models.User{}, errors.New("user exists")
	}
	u.Username = username
	return u, nil
}

func LoginUser(username, password string) (string, int64, error) {
	var id int64
	var hash string
	err := DB.QueryRow(Ctx, "select id, password_hash from users where username=$1", username).Scan(&id, &hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return "", 0, errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"usr": username,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	signed, _ := token.SignedString([]byte(MustEnv("JWT_SECRET")))
	return signed, id, nil
}
