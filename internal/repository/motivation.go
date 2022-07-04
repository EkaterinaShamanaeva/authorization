package repository

import (
	"context"
	"fmt"
)

type Motivation struct {
	Id      int    `db:"id"`
	Content string `db:"cont"`
	Author  string `db:"author"`
}

// GetRandomMotivation возвращает рандомную цитату из БД
func (r *Repository) GetRandomMotivation(ctx context.Context) (m Motivation, err error) {
	// запрос
	row := r.pool.QueryRow(ctx, `select * from motivations order by random() limit 1`)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	// записываем в структуру
	err = row.Scan(&m.Id, &m.Content, &m.Author)
	if err != nil {
		err = fmt.Errorf("2 failed to query data: %w", err)
		return
	}
	return
}
