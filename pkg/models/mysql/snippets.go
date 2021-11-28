// pkg/models/mysql/snippets.go
package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ikaem/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// function to insert new snippet
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// return 0, nil

	// this is sql statement

	query := `
		INSERT INTO snippets(title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`

	// THEN WE call the exect
	// we get result back , result has some useful info about the operation

	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	// then we get that last inert id from the same named method

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// convert and return the id

	return int(id), nil
}

// funciton to return a specific snippet
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// return nil, nil

	// again, the query

	query := `
		SELECT id, title, content, created, expires FROM snippets 
		WHERE expires > UTC_TIMESTAMP() and id = ?
		`
	s := &models.Snippet{}
	// create zeroed initialized struct of the snippet

	// then we get row from the statement
	// row := m.DB.QueryRow(query, id)

	// // then we scan the row
	// // we copy the row values into addresses of the zeroed intiialized Snppet struct

	// err := row.Scan(
	// 	&s.ID,
	// 	&s.Title,
	// 	&s.Content,
	// 	&s.Created,
	// 	&s.Expires,
	// )

	err := m.DB.QueryRow(query, id).Scan(
		&s.ID,
		&s.Title,
		&s.Content,
		&s.Created,
		&s.Expires,
	)

	// now we check for errors

	if err != nil {
		// we check specificaly for now rows error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil

}

// function to return 10 latest snippets
// note the return type - slice of snippet models
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	// note that we can return nil for everythign
	// return nil, nil

	query := `
	SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() 
	ORDER BY created 
	DESC
	LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// we need to call closing of the rows, defer, and do it after error check
	defer rows.Close()

	// initialize empty slice to hold models snippets objects

	snippets := []*models.Snippet{}

	// we now loop of rows in the resultset
	// it will prepare first and each subseuqnet row to be acted on by the rows.scan

	fmt.Printf("%v", rows)

	for rows.Next() {
		// create point to a new zeroed snippet struct

		s := &models.Snippet{}

		// we can now access current result row because we called .Next()
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// now we append values to the snippers slice

		snippets = append(snippets, s)
	}

	// now we call rows.err to check if there was any error during iteration

	if err = rows.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("%v", snippets)

	return snippets, nil

}

func (m *SnippetModel) ExampleTransaction() error {
	// so we need to call Being method on the connection pool
	// this crreates an in progress database transaction

	tx, err := m.DB.Begin()

	if err != nil {
		return err
	}

	// now we just work on the tanaction, instead of database

	_, err = tx.Exec("INSERT INTO...")

	if err != nil {
		// in case of an error, we xcall rollback to abort the transaction
		tx.Rollback()
		return err

	}

	// and if all is good, we have to commit

	err = tx.Commit()

	return err

}
