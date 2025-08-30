package services

import (
	"messenger/models"
	"time"
)

// Сохранение нового сообщения
func SaveMessage(userID int64, body string, username string) (models.Message, error) {
	var m models.Message
	err := DB.QueryRow(Ctx,
		`insert into messages (user_id, body) values ($1,$2) returning id, created_at`,
		userID, body,
	).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		return m, err
	}
	m.UserID = userID
	m.Username = username
	m.Body = body
	return m, nil
}

// История сообщений
func GetMessages(limit int, beforeID int64) ([]models.Message, error) {
	var rows any
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
		return nil, err
	}
	defer rows.(interface{ Close() }).Close()

	var out []models.Message
	for rows.(interface{ Next() bool }).Next() {
		var m models.Message
		var createdAt time.Time
		if err := rows.(interface{ Scan(...any) error }).Scan(&m.ID, &m.UserID, &m.Username, &m.Body, &createdAt); err != nil {
			return nil, err
		}
		m.CreatedAt = createdAt
		out = append(out, m)
	}
	return out, nil
}
