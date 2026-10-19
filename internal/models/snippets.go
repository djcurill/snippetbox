package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// posstgres does not support LastInsertId(), so we must specify
	// in our query to return id. we use place holders to safely inject
	// values into the statement
	stmt := `insert into snippets (title, content, created, expires) values
	($1, $2, now(), now() + make_interval(days => $3))
	returning id`

	var id int

	err := sm.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (sm *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func (sm *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
