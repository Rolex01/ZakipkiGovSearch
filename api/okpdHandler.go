package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// OKPDHandler provide search in OKPD1 & OKPD2 dictionaries. Returns JSON array.
func OKPDHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := url.Values{}
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	if r.Method == "GET" {
		params = r.URL.Query()
	}
	w.Header().Set("Content-type", "application/json")
	query := params.Get("query")
	version := params.Get("v")
	log.Println("OKPD API:", version, query)
	if len(query) != 0 && len(version) == 1 {
		index := "star_nsiokpd" + version
		switch version {
		case "1":
			if b, err := json.Marshal(IdsToOKPD1(sphinxSearch(query, index), db)); err != nil {
				log.Println(err)
			} else {
				io.WriteString(w, string(b))
			}
		case "2":
			if b, err := json.Marshal(IdsToOKPD2(sphinxSearch(query, index), db)); err != nil {
				log.Println(err)
			} else {
				io.WriteString(w, string(b))
			}
		}
	}
}

// OKPDChildrenHandler provides list of children for selected OKPD row.
func OKPDChildrenHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	params := url.Values{}
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	if r.Method == "GET" {
		params = r.URL.Query()
	}
	w.Header().Set("Content-type", "application/json")
	row, _ := strconv.Atoi(params.Get("row"))
	version := params.Get("v")
	log.Println("Children test:", row, version)
	switch version {
	case "1":
		if b, err := json.Marshal(OKPD1toIds(IdsToOKPD1([]int{row}, db)[0].GetChildren(db))); err != nil {
			log.Println(err)
		} else {
			io.WriteString(w, string(b))
		}
	case "2":
		if b, err := json.Marshal(OKPD2toIds(IdsToOKPD2([]int{row}, db)[0].GetChildren(db))); err != nil {
			log.Println(err)
		} else {
			io.WriteString(w, string(b))
		}
	}
}
