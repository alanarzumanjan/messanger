package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	services "messenger/services"
)

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 50
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	beforeID := int64(0)
	if v := q.Get("before_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			beforeID = n
		}
	}

	msgs, err := services.GetMessages(limit, beforeID)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(msgs)
}
