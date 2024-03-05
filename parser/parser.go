package parser

import (
	"encoding/binary"
	"fmt"
	"io"
)

type DB struct {
	Size    int64
	Entries []Entry
}

type Entry struct {
	Key       string
	Value     float64
	Timestamp float64
}

func Parse(r io.ReadSeeker) (db *DB, err error) {
	db = new(DB)
	if err = binary.Read(r, binary.LittleEndian, &db.Size); err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}
	var pos int64 = 8

	for pos < db.Size {
		e := Entry{}

		var length uint32
		if err = binary.Read(r, binary.LittleEndian, &length); err != nil {
			return nil, fmt.Errorf("parser error: %w", err)
		}
		key := make([]byte, length)
		if _, err = io.ReadFull(r, key); err != nil {
			return nil, fmt.Errorf("parser error: %w", err)
		}
		e.Key = string(key)

		padding := 8 - (length+4)%8
		if pos, err = r.Seek(int64(padding), io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("parser error: %w", err)
		}

		if err = binary.Read(r, binary.LittleEndian, &e.Value); err != nil {
			return nil, fmt.Errorf("parser error: %w", err)
		}
		if err = binary.Read(r, binary.LittleEndian, &e.Timestamp); err != nil {
			return nil, fmt.Errorf("parser error: %w", err)
		}
		pos += 16

		db.Entries = append(db.Entries, e)
	}
	return db, nil
}
