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

// wraps a sql.DB
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `insert into snippets (title, content, created, expires)
        values(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `select id, title, content, created, expires from snippets
             where expires > UTC_TIMESTAMP() and id = ?`
	row := m.DB.QueryRow(stmt, id)

	// initalize a new zeroed Snippet struct
	var s Snippet

	// notice this is getting pointers into the structure
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if query returns no rows, then row.Scan will return sql.ErrNoRows.
		// errors.Is() to check for that specifically, and return our own
		// ErrNoRecord error
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	
	// happiness and light
	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `select id, title, content, created, expires from snippets
     where expires > UTC_TIMESTAMP() order by id desc limit 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// do the defer after checking the Query error, otherwise get a panic
	defer rows.Close()

	// initialize an empty slice
	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// check rows.Err() to get any errors during iteration. might not have
	// the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return snippets, nil
}
