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

type Organization struct {
	row               int
	Regnum            string `json:"Regnum"`
	Shortname         string `json:"Shortname"`
	Fullname          string `json:"Fullname"`
	Region            int    `json:"Region"`
	OKVED             string `json:"Okved"`
	URL               string `json:"Url",omitempty`
	INN               string `json:"Inn"`
	OGRN              string `json:"Ogrn"`
	KPP               string `json:"Kpp"`
	OKTMO             string `json:"Oktmo"`
	PostAddr          string `json:"PostAddr"`
	ContactPerson     string `json:"Contact"`
	SubordinationType string `json:"Type"`
}

func customerSearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	params := url.Values{}

	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}

	if r.Method == "GET" {
		params = r.URL.Query()
	}

	w.Header().Set("Content-type", "application/json")

	var orgs []Organization
	var err error
	var body []byte

	if query := params.Get("query"); len(query) != 0 {

		orgs, err = retrieveOrganizations(sphinxOrgSearch(query), db)
		
		if err != nil {
			log.Println("Orgs Search Error: ", err)
		}
		
		body, err = json.Marshal(orgs)

	} else {
		log.Println("else")
		body, err = json.Marshal(orgs)
	}

	if err != nil {
		log.Println("Orgs Search Error: ", err)
	}

	io.WriteString(w, string(body))
}

func sphinxOrgSearch(q string) []int {
	
	var ids []int
	
	sc := sphinx.NewClient().SetServer(sphinxSwitch.GetHost(), 0).SetConnectTimeout(1000).SetLimits(0, 200, MaxDocs, 0)
	
	if err := sc.Error(); err != nil {
		log.Fatal(err)
	}

	index := "star_nsiorg"
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

func retrieveOrganizations(ids []int, db *sql.DB) ([]Organization, error) {

	data := make([]Organization, len(ids))
	
	if len(ids) == 0 {
		return data, nil
	}

	query := `SELECT regnum, shortname, substring(kpp,1,2)::int as region, inn, kpp, fullname  FROM nsiorg WHERE actual = TRUE and row in (`
	
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
	
		err := rows.Scan(&data[j].Regnum, &data[j].Shortname, &data[j].Region, &data[j].INN, &data[j].KPP, &data[j].Fullname)
	
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
