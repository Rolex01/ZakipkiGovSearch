package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	"bitbucket.org/company-one/tender-one/sphinx-switch"
	"github.com/yunge/sphinx"
)

//OKPD interface for OKPD/OKPD2 dictionaries
type OKPD interface {
	GetRow() int
	GetParentRow() int
	GetCode() string
	GetParent() string
	GetName() string
	JSON() string
	GetChildren(*sql.DB) []OKPD
}

//OKPD1 implementation of OKPD dict
type OKPD1 struct {
	Row        int    `json:"Row"`
	ParentRow  int    `json:"ParentRow"`
	Code       string `json:"Code"`
	ParentCode string `json:"ParentCode"`
	Name       string `json:"Name"`
}

//GetRow returns primary id
func (c OKPD1) GetRow() int {
	return c.Row
}

//GetParentRow returns primary id of a parent element. Equals zero if there is no parent element.
func (c OKPD1) GetParentRow() int {
	return c.ParentRow
}

//GetCode returns okpd code as a string.
func (c OKPD1) GetCode() string {
	return c.Code
}

//GetParentCode returns okpd code of a parent element as a string. Null if there is no parent element.
func (c OKPD1) GetParentCode() string {
	return c.ParentCode
}

//GetName returns description of an element.
func (c OKPD1) GetName() string {
	return c.Name
}

//JSON returns JSON encrypted string
func (c OKPD1) JSON() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Println("Error on marshaling OKPD1: ", err)
	}
	return string(b)
}

//GetChildren returns all children of an element
func (c OKPD1) GetChildren(db *sql.DB) (children []OKPD1) {
	//log.Println("GetChildren of ", c.GetCode(), c.GetRow())
	rows, err := db.Query("select row,prow,code,parent,name from nsiOKPD1 where prow=$1;", c.GetRow())
	if err != nil {
		log.Println("Error on GetChildren:", c.GetRow(), ";", err)
	}
	//log.Println("What the fuck? ")
	for rows.Next() {
		var child OKPD1
		var prow sql.NullInt64
		var parentCode sql.NullString
		if err := rows.Scan(&child.Row, &prow, &child.Code, &parentCode, &child.Name); err != nil {
			log.Println("Error on GetChildren:", c.GetRow(), ";", err)
		}
		//log.Println("Scanning: ", child, prow, parentCode, parentCode.Valid, prow.Valid)
		if parentCode.Valid {
			child.ParentCode = parentCode.String
		}
		if prow.Valid {
			child.ParentRow = int(prow.Int64)
			//chlds := child.GetChildren(db)
			//	children = append(children, chlds...)
		}
		children = append(children, child)
	}
	rows.Close()
	tchildren := append([]OKPD1(nil), children...)
	for _, ch := range tchildren {
		children = append(children, ch.GetChildren(db)...)
	}
	//log.Println("Finish children func for: ", c.GetRow(), ";", children)
	return
}

//OKPD2 implementation of OKPD dict
type OKPD2 struct {
	Row        int    `json:"Row"`
	ParentRow  int    `json:"ParentRow"`
	Code       string `json:"Code"`
	ParentCode string `json:"ParentCode"`
	Name       string `json:"Name"`
}

//GetRow returns primary id
func (c OKPD2) GetRow() int {
	return c.Row
}

//GetParentRow returns primary id of a parent element. Equals zero if there is no parent element.
func (c OKPD2) GetParentRow() int {
	return c.ParentRow
}

//GetCode returns okpd code as a string.
func (c OKPD2) GetCode() string {
	return c.Code
}

//GetParentCode returns okpd code of a parent element as a string. Null if there is no parent element.
func (c OKPD2) GetParentCode() string {
	return c.ParentCode
}

//GetName returns description of an element.
func (c OKPD2) GetName() string {
	return c.Name
}

//JSON returns JSON encrypted string
func (c OKPD2) JSON() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Println("Error on marshaling OKPD2: ", err)
	}
	return string(b)
}

//GetChildren returns all children of an element
func (c OKPD2) GetChildren(db *sql.DB) (children []OKPD2) {
	//log.Println("GetChildren of ", c.GetCode(), c.GetRow())
	rows, err := db.Query("select row,prow,code,parent,name from nsiokpd2 where prow=$1;", c.GetRow())
	if err != nil {
		log.Println("Error on GetChildren:", c.GetRow(), ";", err)
	}
	//log.Println("What the fuck? ")
	for rows.Next() {
		var child OKPD2
		var prow sql.NullInt64
		var parentCode sql.NullString
		if err := rows.Scan(&child.Row, &prow, &child.Code, &parentCode, &child.Name); err != nil {
			log.Println("Error on GetChildren:", c.GetRow(), ";", err)
		}
		//log.Println("Scanning: ", child, prow, parentCode, parentCode.Valid, prow.Valid)
		if parentCode.Valid {
			child.ParentCode = parentCode.String
		}
		if prow.Valid {
			child.ParentRow = int(prow.Int64)
			//chlds := child.GetChildren(db)
			//	children = append(children, chlds...)
		}
		children = append(children, child)
	}
	rows.Close()
	tchildren := append([]OKPD2(nil), children...)
	for _, ch := range tchildren {
		children = append(children, ch.GetChildren(db)...)
	}
	//log.Println("Finish children func for: ", c.GetRow(), ";", children)
	return
}

func sphinxSearch(q string, index string) (ids []int) {
	// log.Println(q, index)
	sc := sphinx.NewClient().SetServer(sphinxSwitch.GetHost(), 0).SetConnectTimeout(1000).SetLimits(0, 100000, 100000, 0)
	
	if err := sc.Error(); err != nil {
		log.Fatal(err)
	}
	
	res, err := sc.Query(q, index, index)
	
	if err != nil {
		log.Println("Query error:", err)
	}

	if res == nil {
		log.Println("Panic!")
		return
	}

	ids = make([]int, len(res.Matches))
	log.Println(len(res.Matches), "/", res.Total)

	for i, match := range res.Matches {
		ids[i] = int(match.DocId)
	}

	log.Println("okpds found: ", res.Total)

	return
}

//IdsToOKPD1 retrieve OKPD1 by row slice
func IdsToOKPD1(ids []int, db *sql.DB) (res []OKPD1) {
	for _, id := range ids {
		var c OKPD1
		var prow sql.NullInt64
		var parentCode sql.NullString
		err := db.QueryRow("select row,prow,code,parent,name from nsiokpd1 where row = $1;", id).Scan(&c.Row, &prow, &c.Code, &parentCode, &c.Name)
		if err != nil {
			log.Println("Error on IdsToOKPD1:", id, ";", err)
		}
		if parentCode.Valid {
			c.ParentCode = parentCode.String
		}
		if prow.Valid {
			c.ParentRow = int(prow.Int64)
		}
		res = append(res, c)
	}
	return
}

//StringIdsToOKPD1 returns OKPD1
func StringIdsToOKPD1(sids []string, db *sql.DB) (res []OKPD1) {
	for _, id := range sids {
		var c OKPD1
		var prow sql.NullInt64
		var parentCode sql.NullString
		err := db.QueryRow("select row,prow,code,parent,name from nsiokpd1 where row = "+id+";").Scan(&c.Row, &prow, &c.Code, &parentCode, &c.Name)
		if err != nil {
			log.Println("Error on IdsToOKPD1:", id, ";", err)
		}
		if parentCode.Valid {
			c.ParentCode = parentCode.String
		}
		if prow.Valid {
			c.ParentRow = int(prow.Int64)
		}
		res = append(res, c)
	}
	return
}

//StringIdsToOKPD2 returns OKPD2
func StringIdsToOKPD2(sids []string, db *sql.DB) (res []OKPD2) {

	rows, err := db.Query("select row,prow,code,parent,name from nsiokpd2 where row in (" + strings.Join(sids, ",") + ") order by row;")
	// log.Println("select row,prow,code,parent,name from nsiokpd2 where row in (" + strings.Join(sids, ",") + ") order by row;")
	defer rows.Close()
	for rows.Next() {
		var c OKPD2
		var prow sql.NullInt64
		var parentCode sql.NullString
		rows.Scan(&c.Row, &prow, &c.Code, &parentCode, &c.Name)
		if err != nil {
			log.Println("Error on IdsToOKPD2:", err)
		}

		if parentCode.Valid {
			c.ParentCode = parentCode.String
		}
		if prow.Valid {
			c.ParentRow = int(prow.Int64)
		}
		res = append(res, c)
	}
	return
}

//IdsToOKPD2 retrieve OKPD1 by row slice
func IdsToOKPD2(ids []int, db *sql.DB) (res []OKPD2) {
	for _, id := range ids {
		var c OKPD2
		var prow sql.NullInt64
		var parentCode sql.NullString
		err := db.QueryRow("select row,prow,code,parent,name from nsiokpd2 where row = $1;", id).Scan(&c.Row, &prow, &c.Code, &parentCode, &c.Name)
		if err != nil {
			log.Println("Error on IdsToOKPD1:", id, ";", err)
		}
		if parentCode.Valid {
			c.ParentCode = parentCode.String
		}
		if prow.Valid {
			c.ParentRow = int(prow.Int64)
		}
		res = append(res, c)
	}
	return
}

//OKPD1toIds collect primary ids from OKPD structures
func OKPD1toIds(cs []OKPD1) (ids []uint64) {
	for _, c := range cs {
		ids = append(ids, uint64(c.GetRow()))
	}
	return
}

//OKPD2toIds collect primary ids from OKPD structures
func OKPD2toIds(cs []OKPD2) (ids []uint64) {
	for _, c := range cs {
		ids = append(ids, uint64(c.GetRow()))
	}
	return
}
