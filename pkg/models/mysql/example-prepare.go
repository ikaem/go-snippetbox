package mysql

import (
	"database/sql"
)

// this is just an example model
// we will store teh stirng query here

type ExampleModel struct {
	DB              *sql.DB
	InsertStatement *sql.Stmt
	// and then we could add bunch of other fields, depenmding on which queries we have
}

// this is a constructor for the model
func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
	// here we use prepare method to creeate a new prperated statement
	// it returns statement object

	// WE CAN CREATE BUNCH OF others statments too
	insertStatement, err := db.Prepare("INSERT INTO...")

	if err != nil {
		return nil, err
	}

	// and then we just store the thing  in the struct model
	// we also stroe db pool, in case we want to use it

	return &ExampleModel{db, insertStatement}, nil
}

func (m *ExampleModel) Insert(args ...string) error {
	// we call exec on the prepared statement
	// we dont call it on the connection pool

	// note how we pass variadic props
	_, err := m.InsertStatement.Exec(args)

	return err
}
