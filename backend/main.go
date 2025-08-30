package main

import (
	"log"
	"net/http"
	"os"

	controllers "messenger/controllers"
	services "messenger/services"

	"github.com/joho/godotenv"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load(".env")
	services.InitDB()
	services.InitRedis()

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("pong")) })
	mux.HandleFunc("/register", controllers.RegisterHandler)
	mux.HandleFunc("/login", controllers.LoginHandler)
	mux.HandleFunc("/messages", controllers.MessagesHandler)
	mux.HandleFunc("/ws", controllers.WsHandler)

	log.Println("🚀 Server on http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}
