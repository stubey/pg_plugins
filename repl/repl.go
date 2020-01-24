package main

import (
	"bufio"
	"context"
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
	flag.StringVar(&p.cxstring, "cxstring", "host=pgprimary port=5432 user=postgres dbname=replsource sslmode=disable", "dbx connection string")
	flag.StringVar(&p.input, "input", "-", "file from which to read sql commands, '-' for stdin")
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

	ctx := context.Background()

	scanner := openFile(p.input)
	for scanner.Scan() {
		line := scanner.Text()
		if p.verbose {
			log.Printf("line = %v", line)
		}

		if false {
		rows, err := db.Query(line)  need to free rows, else uses new connection
		result, err := db.ExecContext(ctx, line)
		if err != nil {
			log.Printf("err = %v, result = %T - %+v", err, result, result)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}


type 