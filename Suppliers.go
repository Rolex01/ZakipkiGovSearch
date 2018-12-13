package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/yunge/sphinx"
	"bitbucket.org/company-one/tender-one/sphinx-switch"
)

type Supplier struct {
	row  int
	Name string `json:"Name"`
	INN  string `json:"Inn"`
	KPP  string `json:"Kpp"`
}

func supplierSearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	params := url.Values{}
	
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}

	if r.Method == "GET" {
		params = r.URL.Query()
	}

	w.Header().Set("Content-type", "application/json")
	
	var suppliers []Supplier
	var err error
	var body []byte

	if query := params.Get("query"); len(query) != 0 {

		suppliers, err = retrieveSuppliers(sphinxSuppliersSearch(query), db)
		
		if err != nil {
			log.Println("Suppliers Search Error: ", err)
		}
		
		body, err = json.Marshal(suppliers)
	
	} else {
		log.Println("else")
		body, err = json.Marshal(suppliers)
	}

	if err != nil {
		log.Println("Suppliers Search Error: ", err)
	}

	io.WriteString(w, string(body))
}

func sphinxSuppliersSearch(q string) []int {

	var ids []int
	
	sc := sphinx.NewClient().SetServer(sphinxSwitch.GetHost(), 0).SetConnectTimeout(1000).SetLimits(0, 25, MaxDocs, 0)
	
	if err := sc.Error(); err != nil {
		log.Fatal(err)
	}

	index := "star_nsisuppliers"
	res, err := sc.Query(q, index, index)
	
	if err != nil {
		log.Println("Query error:", err)
	}
	
	if res == nil {

		log.Println("Panic!")

		return ids

	} else {
		
		ids = make([]int, len(res.Matches))
		log.Println(len(res.Matches), "/", res.Total)

		for i, match := range res.Matches {
			ids[i] = int(match.DocId)
		}
	}

	return ids
}

func retrieveSuppliers(ids []int, db *sql.DB) ([]Supplier, error) {

	data := make([]Supplier, len(ids))

	if len(ids) == 0 {
		return data, nil
	}

	query := `SELECT name, inn, kpp FROM nsisuppliers WHERE row in (`
	
	for l, i := range ids {
	
		if l != len(ids) - 1 {
			query += strconv.Itoa(i) + ", "
		} else {
			query += strconv.Itoa(i) + ");"
		}
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	j := 0

	for rows.Next() {

		if err != nil {
			log.Fatal(err)
		}
		
		err := rows.Scan(&data[j].Name, &data[j].INN, &data[j].KPP)
		
		if err != nil {
			log.Println(err)
		}
		
		j++
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
	}

	return data, err
}
