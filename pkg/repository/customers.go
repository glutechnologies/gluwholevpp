package repository

type Customer struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	OuterInterface int    `json:"outer-interface"`
	OuterVlan      int    `json:"outer-vlan"`
	Counter        int    `json:"counter"`
}

func (s *Storage) GetCustomers(customers *[]Customer) error {
	rows, err := s.db.Query("SELECT Id, Name, OuterInterface, OuterVlan, Counter FROM customers")

	if err != nil {
		return err
	}

	var id int
	var name string
	var outerInterface int
	var outerVlan int
	var counter int

	for rows.Next() {
		err = rows.Scan(&id, &name, &outerInterface, &outerVlan, &counter)
		if err != nil {
			return err
		}
		(*customers) = append((*customers), Customer{Id: id, Name: name, OuterInterface: outerInterface,
			OuterVlan: outerVlan, Counter: counter})
	}

	return nil
}

func (s *Storage) GetCustomer(customerId int, customer *Customer) error {
	var id int
	var name string
	var outerInterface int
	var outerVlan int
	var counter int
	err := s.db.QueryRow("SELECT Id, Name, OuterInterface, OuterVlan, Counter FROM customers WHERE Id = ?", customerId).Scan(&id, &name, &outerInterface, &outerVlan, &counter)

	if err != nil {
		return err
	}

	(*customer) = Customer{
		Id:             id,
		Name:           name,
		OuterInterface: outerInterface,
		OuterVlan:      outerVlan,
		Counter:        counter,
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

	_, err = tx.Exec("UPDATE customers SET Counter = ? WHERE Id = ?", counter+1, customerId)

	if err != nil {
		tx.Rollback()
		return counter, err
	}

	defer tx.Commit()

	return counter, nil
}

func (s *Storage) InsertCustomer(customer *Customer) error {
	stmt, err := s.db.Prepare("INSERT INTO customers (Name, OuterInterface, OuterVlan) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(customer.Name, customer.OuterInterface, customer.OuterVlan)
	defer stmt.Close()

	return nil
}