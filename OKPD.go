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

type OKPD struct {
	Row  int    `json:"Row"`
	Code string `json:"Code"`
	Name string `json:"Name"`
}

func okpdSearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := url.Values{}

	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	if r.Method == "GET" {
		params = r.URL.Query()
	}

	w.Header().Set("Content-type", "application/json")

	var okpds []OKPD
	var err error
	var body []byte

	//OKPD := new([]OKPD)
	if query := params.Get("query"); len(query) != 0 {
		log.Println("OKPD search: ", query)
		okpds, err = retrieveOKPD(sphinxokpdsearch(query), db)
		if err != nil {
			log.Println("okpds Search Error: ", err)
		}
		body, err = json.Marshal(okpds)
	} else {
		log.Println("else")
		body, err = json.Marshal(okpds)
	}

	if err != nil {
		log.Println("okpds Search Error: ", err)
	}
	//io.WriterTo(w, body)
	io.WriteString(w, string(body))
	//w.Write(body)
}

func sphinxokpdsearch(q string) []int {

	var ids []int
	
	sc := sphinx.NewClient().SetServer(sphinxSwitch.GetHost(), 0).SetConnectTimeout(1000).SetLimits(0, 1000, MaxDocs, 0)
	
	if err := sc.Error(); err != nil {
		log.Fatal(err)
	}
	
	index := "star_nsiokpd"
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
		// log.Println("okpds found: ", res.Total)
	}
	return ids
}

func retrieveOKPD(ids []int, db *sql.DB) ([]OKPD, error) {

	data := make([]OKPD, len(ids))
	if len(ids) == 0 {
		return data, nil
	}

	query := `SELECT row,code,name from nsiokpd where actual = TRUE and row in (`
	for l, i := range ids {
		if l != len(ids) - 1 {
			query += strconv.Itoa(i) + ", "
		} else {
			query += strconv.Itoa(i) + ");"
		}
	}
	// log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	j := 0
	// log.Println(rows)
	for rows.Next() {
		err := rows.Scan(&data[j].Row, &data[j].Code, &data[j].Name)
		if err != nil {
			log.Fatal(err)
		}
		//log.Println(data[j])
		j++
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(data)
	return data, err
}
