package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	filename string
	db       *sql.DB
}

func (s *Storage) Init(filename string) {
	s.filename = filename
}

func (s *Storage) OpenDB() error {
	db, err := sql.Open("sqlite3", s.filename)

	checkErr(err)
	s.db = db

	return err
}

func (s *Storage) CloseDB() error {
	return s.db.Close()
}

func checkErr(err error) {
	fmt.Println(err)
}
