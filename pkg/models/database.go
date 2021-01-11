package models

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	DB *pgxpool.Pool
}

func (db *Database) GetSnippet(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, (created::timestamp(0)), (expires::timestamp(0)) FROM snippets WHERE expires > current_timestamp AND id=$1`

	row := db.DB.QueryRow(context.Background(), stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (db *Database) LatestSnippets() (Snippets, error) {
	stmt := `SELECT id, title, content, (created::timestamp(0)), (expires::timestamp(0)) FROM snippets "+
			"WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10`

	rows, err := db.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := Snippets{}

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (db *Database) InsertSnippet(title, content, expires string) (int, error) {
	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES($1, $2, current_timestamp, current_date + interval '" + expires + " second') returning id"

	_, err := db.DB.Exec(context.Background(), stmt, title, content)
	result := db.DB.QueryRow(context.Background(), "SELECT currval('snippets_id_seq');")
	if err != nil {
		return 0, err
	}
	var id int
	_ = result.Scan(&id)

	return id, nil
}
