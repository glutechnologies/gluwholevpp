package repository

import (
	"database/sql"

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
	s.db = db

	return err
}

func (s *Storage) CloseDB() error {
	return s.db.Close()
}
