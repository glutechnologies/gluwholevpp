package repository

type Bitstream struct {
	Id         string `json:"id" validate:"required"`
	CustomerId string `json:"customer-id" validate:"required"`
	SrcId      int    `json:"src-id"`
	SrcOuter   int    `json:"src-outer" validate:"required"`
	SrcInner   int    `json:"src-inner" validate:"required"`
	DstId      int    `json:"dst-id"`
	DstOuter   int    `json:"dst-outer"`
	DstInner   int    `json:"dst-inner"`
	Comment    string `json:"comment"`
}

func (s *Storage) GetBitstreams(bitstreams *[]Bitstream) error {
	rows, err := s.db.Query("SELECT Id, CustomerId, SrcId, SrcOuter, SrcInner, DstId, DstOuter, DstInner, Comment FROM bitstreams")

	if err != nil {
		return err
	}

	var id string
	var customerId string
	var srcId int
	var srcOuter int
	var srcInner int
	var dstId int
	var dstOuter int
	var dstInner int
	var comment string

	for rows.Next() {
		err = rows.Scan(&id, &customerId, &srcId, &srcOuter, &srcInner, &dstId, &dstOuter, &dstInner, &comment)
		if err != nil {
			return err
		}
		(*bitstreams) = append((*bitstreams), Bitstream{Id: id, CustomerId: customerId, SrcId: srcId, SrcOuter: srcOuter,
			SrcInner: srcInner, DstId: dstId, DstOuter: dstOuter, DstInner: dstInner, Comment: comment})
	}

	return nil
}

func (s *Storage) GetBitstreamsFromCustomer(customerId string, bitstreams *[]Bitstream) error {
	rows, err := s.db.Query("SELECT Id, CustomerId, SrcId, SrcOuter, SrcInner, DstId, DstOuter, DstInner, Comment FROM bitstreams WHERE customerId = ?", customerId)

	if err != nil {
		return err
	}

	for rows.Next() {
		var id string
		var customerId string
		var srcId int
		var srcOuter int
		var srcInner int
		var dstId int
		var dstOuter int
		var dstInner int
		var comment string

		err = rows.Scan(&id, &customerId, &srcId, &srcOuter, &srcInner, &dstId, &dstOuter, &dstInner, &comment)
		if err != nil {
			return err
		}
		(*bitstreams) = append((*bitstreams), Bitstream{Id: id, CustomerId: customerId, SrcId: srcId, SrcOuter: srcOuter,
			SrcInner: srcInner, DstId: dstId, DstOuter: dstOuter, DstInner: dstInner, Comment: comment})
	}

	return nil
}

func (s *Storage) GetBitstream(id string, bitstream *Bitstream) error {

	var customerId string
	var srcId int
	var srcOuter int
	var srcInner int
	var dstId int
	var dstOuter int
	var dstInner int
	var comment string

	err := s.db.QueryRow("SELECT Id, CustomerId, SrcId, SrcOuter, SrcInner, DstId, DstOuter, DstInner, Comment FROM bitstreams WHERE Id = ?",
		id).Scan(&id, &customerId, &srcId, &srcOuter, &srcInner, &dstId, &dstOuter, &dstInner, &comment)

	if err != nil {
		return err
	}

	*(bitstream) = Bitstream{Id: id, CustomerId: customerId, SrcId: srcId, SrcOuter: srcOuter,
		SrcInner: srcInner, DstId: dstId, DstOuter: dstOuter, DstInner: dstInner, Comment: comment}

	return nil
}

func (s *Storage) DeleteBitstream(id string) error {
	stmt, err := s.db.Prepare("DELETE FROM bitstreams WHERE Id = ?")
	if err != nil {
		return err
	}

	stmt.Exec(id)
	defer stmt.Close()

	return nil
}

func (s *Storage) InsertBitstream(bitstream *Bitstream) error {
	stmt, err := s.db.Prepare("INSERT INTO bitstreams (Id, CustomerId, SrcId, SrcOuter, SrcInner, DstId, DstOuter, DstInner, Comment) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(bitstream.Id, bitstream.CustomerId, bitstream.SrcId, bitstream.SrcOuter, bitstream.SrcInner, bitstream.DstId, bitstream.DstOuter, bitstream.DstInner, bitstream.Comment)
	defer stmt.Close()

	return nil
}
