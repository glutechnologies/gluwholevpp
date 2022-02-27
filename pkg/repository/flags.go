package repository

type Flags struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func (s *Storage) IncrementIdCounter() (int, error) {
	var counter int
	tx, _ := s.db.Begin()
	err := tx.QueryRow("SELECT Value FROM flags WHERE Key = ?", "IdCounter").Scan(&counter)

	if err != nil {
		tx.Rollback()
		return counter, err
	}

	_, err = tx.Exec("UPDATE flags SET Value = ? WHERE Key = ?", counter+2, "IdCounter")

	if err != nil {
		tx.Rollback()
		return counter, err
	}

	defer tx.Commit()

	return counter, nil
}
