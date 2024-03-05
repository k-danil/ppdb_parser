package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"ppdb_parser/parser"
)

func main() {
	input := flag.String("input", "", "Input .db file to parse")
	output := flag.String("output", "", "Output .json file")
	flag.Parse()

	if input == nil || *input == "" {
		log.Fatalln("--input is mandatory argument")
	}
	i, err := os.Open(*input)
	if err != nil {
		log.Fatalf("can't open input file: %s", err)
	}

	var db *parser.DB
	if db, err = parser.Parse(i); err != nil {
		log.Fatalln(err)
	}
	uniq := make(map[string]struct{}, len(db.Entries))

	log.Printf("prom_db length: %d", db.Size)
	log.Printf("prom_db entries: %d", len(db.Entries))
	for _, s := range db.Entries {
		uniq[s.Key] = struct{}{}
	}
	log.Printf("uniq entrie count: %d", len(uniq))

	if output != nil && *output != "" {
		var o *os.File
		if o, err = os.Create(*output); err != nil {
			log.Fatalf("can't create output file: %s", err)
		}
		defer func() { _ = o.Close() }()
		if err = json.NewEncoder(o).Encode(&db); err != nil {
			log.Fatalf("can't write to output file: %s", err)
		}
	}
}
