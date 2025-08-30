package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatMessage struct {
	UserID  int64  `json:"user_id"`
	Body    string `json:"body"`
	Created string `json:"created,omitempty"`
}

const ChatChannel = "chat:global"

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Авторизация по JWT из заголовка Authorization: Bearer <token>
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		http.Error(w, "missing bearer", 401)
		return
	}
	raw := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) { return []byte(mustEnv("JWT_SECRET")), nil })
	if err != nil {
		http.Error(w, "bad token", 401)
		return
	}
	uid, _ := claims["sub"].(float64)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Подписка на общий канал
	sub := RDB.Subscribe(Ctx, ChatChannel)
	defer sub.Close()
	ch := sub.Channel()

	// Чтение из сокета и публикация в Redis
	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				break
			}
			msg := ChatMessage{UserID: int64(uid), Body: string(data)}
			// Сохраним в БД
			DB.Exec(Ctx, "insert into messages (user_id, body) values ($1,$2)", msg.UserID, msg.Body)
			// Отправим в общий канал
			b, _ := json.Marshal(msg)
			RDB.Publish(Ctx, ChatChannel, string(b))
		}
	}()

	// Релей сообщений из Redis в сокет
	for m := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(m.Payload)); err != nil {
			break
		}
	}
}
