package repository

type Customer struct {
	Id             int `json:"id"`
	OuterInterface int `json:"outer-interface"`
	OuterVlan      int `json:"outer-vlan"`
	Counter        int `json:"counter"`
}

func (s *Storage) GetCustomers(customers *[]Customer) error {
	rows, err := s.db.Query("SELECT Id, OuterInterface, OuterVlan FROM customers")

	if err != nil {
		return err
	}

	var id int
	var outerInterface int
	var outerVlan int

	for rows.Next() {
		err = rows.Scan(&id, &outerInterface, &outerInterface, &outerVlan)
		if err != nil {
			return err
		}
		(*customers) = append((*customers), Customer{Id: id, OuterInterface: outerInterface,
			OuterVlan: outerVlan})
	}

	return nil
}

func (s *Storage) GetCustomer(customerId int, customer *Customer) error {
	var id int
	var outerInterface int
	var outerVlan int
	err := s.db.QueryRow("SELECT Id, OuterInterface, OuterVlan FROM customers WHERE Id = ?", customerId).Scan(&id, &outerInterface, &outerVlan)

	if err != nil {
		return err
	}

	(*customer) = Customer{
		Id:             id,
		OuterInterface: outerInterface,
		OuterVlan:      outerVlan,
	}

	return nil
}

func (s *Storage) IncrementCounterCustomer(customerId int) (int, error) {
	var counter int
	tx, _ := s.db.Begin()
	err := tx.QueryRow("SELECT Counter FROM customers WHERE Id = ?", customerId).Scan(&counter)

	if err != nil {
		tx.Rollback()
		return counter, err
	}

	_, err = tx.Exec("UPDATE customers SET Counter = ?", counter+1)

	if err != nil {
		tx.Rollback()
		return counter, err
	}

	defer tx.Commit()

	return counter, nil
}

func (s *Storage) InsertCustomer(customer *Customer) error {
	stmt, err := s.db.Prepare("INSERT INTO customers (OuterInterface, OuterVlan, Counter) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(customer.OuterInterface, customer.OuterInterface, customer.OuterVlan)
	defer stmt.Close()

	return nil
}
