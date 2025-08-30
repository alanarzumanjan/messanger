package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	services "messenger/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type inboundMsg struct {
	Body string `json:"body"`
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		http.Error(w, "missing bearer", 401)
		return
	}
	raw := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
		return []byte(services.MustEnv("JWT_SECRET")), nil
	})
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

	sub := services.RDB.Subscribe(services.Ctx, "chat:global")
	defer sub.Close()
	ch := sub.Channel()

	// чтение из сокета → сохранить → публиковать
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
			msg, err := services.SaveMessage(int64(uid), in.Body, uname)
			if err != nil {
				continue
			}
			b, _ := json.Marshal(msg)
			services.RDB.Publish(services.Ctx, "chat:global", string(b))
		}
	}()

	// релей из Redis → клиенту
	for m := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(m.Payload)); err != nil {
			break
		}
	}
}
