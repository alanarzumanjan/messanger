package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func handleMessages(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 50
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	beforeID := int64(0)
	if v := q.Get("before_id"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			beforeID = n
		}
	}

	var rows RowsCloser
	var err error
	if beforeID > 0 {
		rows, err = DB.Query(Ctx, `
			select m.id, m.user_id, u.username, m.body, m.created_at
			from messages m
			join users u on u.id = m.user_id
			where m.id < $1
			order by m.id desc
			limit $2
		`, beforeID, limit)
	} else {
		rows, err = DB.Query(Ctx, `
			select m.id, m.user_id, u.username, m.body, m.created_at
			from messages m
			join users u on u.id = m.user_id
			order by m.id desc
			limit $1
		`, limit)
	}
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	defer rows.Close()

	type item struct {
		ID        int64  `json:"id"`
		UserID    int64  `json:"user_id"`
		Username  string `json:"username"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
	}
	var out []item
	for rows.Next() {
		var it item
		var createdAt time.Time
		if err := rows.Scan(&it.ID, &it.UserID, &it.Username, &it.Body, &createdAt); err != nil {
			http.Error(w, "scan error", 500)
			return
		}
		it.CreatedAt = createdAt.UTC().Format(time.RFC3339Nano)
		out = append(out, it)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// RowsCloser — маленький интерфейс совместимый с pgx.Rows
type RowsCloser interface {
	Next() bool
	Scan(...any) error
	Close()
}
