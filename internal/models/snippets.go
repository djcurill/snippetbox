package models

import (
	"database/sql"
	"errors"
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
	// prepared staement with id placeholder
	stmt := `select id, title, content, created, expires from snippets where expires > now() and id = $1;`

	// query for a single row, returns a sql.Row
	row := sm.DB.QueryRow(stmt, id)

	// create a snippet type var which will be the destination point when calling
	// row.Scan
	var snippet Snippet
	// arguments are pointers specifying where data will be copied to
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Expires, &snippet.Created)

	// error checking
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// return a custom error so we don't rely on the underlying
			// datastore specific errors
			return Snippet{}, ErrNoRecord
		}
		return Snippet{}, err
	}

	return snippet, nil
}

func (sm *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `select id, title, content, created, expires from snippets
				order by created desc
				limit 10;`

	rows, err := sm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// very important to close rows if there is no immediate error
	// forgetting to close would leave a connection open
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return snippets, nil
}
