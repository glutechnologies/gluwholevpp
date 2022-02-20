package repository

type Bitstream struct {
	Id         string `json:"id"`
	CustomerId int    `json:"customer-id"`
	SrcId      int    `json:"src-id"`
	SrcOuter   int    `json:"src-outer"`
	SrcInner   int    `json:"src-inner"`
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
	var customerId int
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
		checkErr(err)
		(*bitstreams) = append((*bitstreams), Bitstream{Id: id, CustomerId: customerId, SrcId: srcId, SrcOuter: srcOuter,
			SrcInner: srcInner, DstId: dstId, DstOuter: dstOuter, DstInner: dstInner, Comment: comment})
	}

	return nil
}

func (s *Storage) InsertBitstream(bitstream *Bitstream) error {
	stmt, err := s.db.Prepare("INSERT INTO bitstreams (Id, CustomerId, SrcOuter, SrcInner, DstOuter, DstInner, Comment) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	stmt.Exec(bitstream.Id, bitstream.CustomerId, bitstream.SrcOuter, bitstream.SrcInner, bitstream.DstOuter, bitstream.DstInner, bitstream.Comment)
	defer stmt.Close()

	return nil
}
