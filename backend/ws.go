package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type inboundMsg struct {
	Body string `json:"body"`
}

type ChatMessage struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username,omitempty"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
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
	uname, _ := claims["usr"].(string)

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
			var in inboundMsg
			if err := json.Unmarshal(data, &in); err != nil || strings.TrimSpace(in.Body) == "" {
				continue
			}
			// Сохраним в БД с возвратом id/created_at
			var id int64
			var createdAt time.Time
			if err := DB.QueryRow(Ctx,
				`insert into messages (user_id, body) values ($1,$2) returning id, created_at`,
				int64(uid), in.Body,
			).Scan(&id, &createdAt); err != nil {
				continue
			}
			// Отправим в общий канал
			out := ChatMessage{
				ID:        id,
				UserID:    int64(uid),
				Username:  uname,
				Body:      in.Body,
				CreatedAt: createdAt.UTC().Format(time.RFC3339Nano),
			}
			b, _ := json.Marshal(out)
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
