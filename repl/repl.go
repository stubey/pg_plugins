package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"database/sql"

	_ "github.com/lib/pq"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type parms struct {
	cxstring string
	input    string
	verbose  bool
}

func getParms() parms {
	p := parms{}
	flag.StringVar(&p.cxstring, "cxstring", "host=localhost port=5432 user=tom dbname=tom sslmode=disable", "dbx connection string")
	flag.StringVar(&p.input, "input", "in.sql", "file from which to read sql commands")
	flag.BoolVar(&p.verbose, "verbose", false, "print commands prior to sql execution")
	flag.Parse()
	return p
}

func openFile(f string) *bufio.Scanner {
	if f == "-" {
		return bufio.NewScanner(os.Stdin)
	}

	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	return bufio.NewScanner(file)
}

func main() {
	p := getParms()
	if p.verbose {
		log.Printf("parms = %+v", p)
	}

	db, err := sql.Open("postgres", p.cxstring)
	if err != nil {
		log.Fatal(err)
	}

	scanner := openFile(p.input)
	for scanner.Scan() {
		line := scanner.Text()
		if p.verbose {
			log.Printf("line = %v", line)
		}

		rows, err := db.Query(line)
		if err != nil {
			log.Printf("err = %v, rows = %T - %+v", err, rows, rows)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
