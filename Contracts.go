package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"

	cbrate "bitbucket.org/company-one/cbrrate"
	shards "bitbucket.org/company-one/tender-one/shards"
	fields "bitbucket.org/company-one/tender-one/document-fields"
	"bitbucket.org/company-one/tender-one/api"
	"bitbucket.org/company-one/tender-one/sphinx-switch"
	uquery "bitbucket.org/company-one/tender-one/query"

	_ "github.com/lib/pq"
	"github.com/yunge/sphinx"
	//	"regexp"
	"strconv"
	"strings"
	"html/template"
	"time"
)

type Product struct {
	Name            string         `json:"name"`
	Price           float64        `json:"price"`
	Quantity        float64        `json:"quantity"`
	Unit            string         `json:"unit"`
	Sum             float64        `json:"sum"`
	Rate            float64        `json:"rate"`
	TotalSum        float64        `json:"sum"`
	OKPD            string         `json:"OKPD"`
	OKPDInfo        sql.NullString `json:"OKPDInfo"`
	Customer        string         `json:"customer"`
	CustomerINN     string         `json:"customerINN"`
	CustomerKPP     sql.NullString `json:"customerKPP"`
	CustomerRegion  uint64         `json:"customerRegion"`
	OrgCode         sql.NullString `json:"orgCode"`
	Signed          time.Time      `json:"signed"`
	Exec            time.Time      `json:"exec"`
	Supplier        sql.NullString `json:"supplier"`
	SupplierINN     sql.NullString `json:"supplierINN"`
	SupplierKPP     sql.NullString `json:"supplierKPP"`
	SupplierAddress sql.NullString `json:"supplierAddress"`
	SupplierPhone   sql.NullString `json:"supplierPhone"`
	Budget          sql.NullString `json:"budget"`
	BudgetSource    sql.NullString `json:"budgetSource"`
	Status          string         `json:"status"`
	RegNum          string         `json:"regnum"`
	Number          sql.NullString `json:"number"`
	PurchaseNumber  sql.NullString `json:"purchaseNumber"`
	Placing         sql.NullString `json:"placing"`
	URL             string         `json:"url"`
	Currency        string         `json:"currency"`
	Paid            sql.NullFloat64        `json:"paid"`
	NMCK            sql.NullFloat64
	ObjectInfo      sql.NullString
	Query           string
}

// ContractsBag - handle shards information
type ContractsBag struct {
	shard string
	id    int
}

// ContractResult - content which goes to template
type ContractResult struct {
	Email              string    `json:"email"`
	Data               []Product `json:"data"`
	Total              int       `json:"total"`
	Count              int       `json:"count"`
	Query              string    `json:"query"`
	FilterDateFrom     string    `json:"DateFrom"`
	FilterDateTo       string    `json:"DateTo"`
	FilterDateExecFrom string    `json:"DateEFrom"`
	FilterDateExecTo   string    `json:"DateETo"`
	FilterPriceFrom    float64   `json:"PriceFrom"`
	FilterPriceTo      float64   `json:"PriceTo"`
	FilterSumFrom      float64   `json:"SumFrom"`
	FilterSumTo        float64   `json:"SumTo"`
	FilterCustomer     []string  `json:"Customer"`
	FilterSupplier     []string  `json:"Supplier"`
	FilterRegion       []string  `json:"Region"`
	FilterCode         []string  `json:"OKPDCode"`
	FilterCode2        []string  `json:"OKPDCode2"`
	FilterPnum         string    `json:"Anum"`
	FilterRegnum       string    `json:"Cnum"`
	FilterPage         int       `json:"filterPage"`
	FilterStrictFlag   bool      `json:"filterStrictFlag"`
	PerPage            int       `json:"perPage"`
	Fields             []fields.Field
	IsSearch           bool
	Demo               bool
	Lvl                int
	NavType            string
}

func (p Product) PurchaseNumberFix() string {
	if p.PurchaseNumber.Valid != true {
		log.Println("Pnum error:", p.RegNum)
		return ""

	}
	s, _ := p.PurchaseNumber.Value()
	//log.Println(p.RegNum, ";", s)
	if s.(string) == "" {
		return "б/н"
	} else {
		return s.(string)
	}
}
func (p Product) FixBudgetSource() string {
	if p.BudgetSource.Valid != true {
		return ""
	}
	s, _ := p.BudgetSource.Value()
	if s.(string) == "00" || s.(string) == "" {
		return "не указан"
	} else {
		return s.(string)
	}
}

func (p Product) FixBudget() string {
	if p.Budget.Valid != true {
		return ""
	}
	s, _ := p.Budget.Value()
	if s.(string) == "" {
		return "не указан"
	} else {
		return s.(string)
	}
}

func (result *ContractResult) setMeta(query string, f ContractFilters) {
	log.Println("Meta Pew-Pew", f.filterPriceFrom)
	result.Query = query
	if f.filterDateFrom == "01.01.2014" {
		f.filterDateFrom = "01.01.2014"
	}

	if f.filterDateTo == "31.12.2018" {
		f.filterDateTo = "31.12.2018" // ToDo - set to current Date
	}
	/*
		if f.filterDateExecFrom == "01.01.2014" {
			f.filterDateExecFrom = "01.01.2014"
		}

		if f.filterDateExecTo == "31.12.2018" {
			f.filterDateExecTo = "31.12.2018" // ToDo - set to current Date
		}*/
	if f.demo {
		log.Println("Demo user!")
		f.filterDateFrom = "01.01.2017"
		f.filterDateTo = "31.12.2017"
	}
	// log.Println(result.Query)
	result.FilterDateFrom = f.filterDateFrom
	// log.Println(result.FilterDateFrom)
	result.FilterDateTo = f.filterDateTo
	// log.Println(result.FilterDateTo)
	result.FilterDateExecFrom = f.filterDateExecFrom
	// log.Println(result.FilterDateExecFrom)
	result.FilterDateExecTo = f.filterDateExecTo
	// log.Println(result.FilterDateExecTo)
	result.FilterPriceFrom = float64(f.filterPriceFrom)

	r := *result
	r.FilterPriceFrom = float64(f.filterPriceFrom)
	log.Println("SMeta Pew-Pew", r.FilterPriceFrom)
	result.FilterPriceTo = float64(f.filterPriceTo)
	result.FilterPriceFrom = float64(f.filterSumFrom)
	result.FilterPriceTo = float64(f.filterSumTo)
	result.FilterCustomer = f.filterCustomer
	result.FilterSupplier = f.filterSupplier
	result.FilterRegion = f.filterRegion
	result.FilterCode = f.filterCode
	result.FilterCode2 = f.filterCode2
	result.FilterPnum = strings.Join(f.filterPnum, ",")
	result.FilterRegnum = strings.Join(f.filterRegnum, ",")
	result.FilterPage = int(f.filterPage)
	result.PerPage = int(f.perPage)
	result.NavType = "contracts"

}
func omitempty(t string, v interface{}) string {

	switch v.(type) {
	case string:
		if v != "" {
			return fmt.Sprintf("&"+t+"=%s", url.QueryEscape(v.(string)))
		} else {
			return ""
		}
	case float64:
		if v.(float64) == 0 {
			return ""
		}
		return fmt.Sprintf("&"+t+"=%2f", v.(float64))
	case float32:
		if v.(float32) == 0 {
			return ""
		}
		return fmt.Sprintf("&"+t+"=%2f", v.(float32))
	case int:
		if v.(int) == 0 {
			return ""
		}
		return fmt.Sprintf("&"+t+"=%d", v.(int))
	case bool:
		return fmt.Sprintf("&"+t+"=%t", v.(bool))
	case []string:
		a := v.([]string)
		//log.Println("Debug 1 omitempty:", a)
		if len(a) != 0 {
			//log.Println("Debug 2 omitempty:", "&"+t+"=%s", urlEscapeStringArray(a, t))
			return fmt.Sprintf("&"+t+"=%s", urlEscapeStringArray(a, t))
		} else {
			return ""
		}
	default:
		return ""
	}
}
func (r ContractResult) ContractsToTemplate() string {
	return "/templates/?save=" + url.QueryEscape(r.genURL(false, 0))
}
func (r *ContractResult) genURL(xls bool, page int) string {
	s := "/contracts/?query=" + url.QueryEscape(r.Query)
	s += omitempty("Supplier", r.FilterSupplier)
	s += omitempty("DateFrom", r.FilterDateFrom)
	s += omitempty("DateTo", r.FilterDateTo)
	s += omitempty("DateEFrom", r.FilterDateExecFrom)
	s += omitempty("DateETo", r.FilterDateExecTo)
	s += omitempty("Customer", r.FilterCustomer)
	s += omitempty("OKPDCode", r.FilterCode)
	s += omitempty("OKPDCode2", r.FilterCode2)
	s += omitempty("Region", r.FilterRegion)
	s += omitempty("PriceFrom", r.FilterPriceFrom)
	//log.Println("OMIT PRICE(", r.FilterPriceFrom, ") FROM:", omitempty("filterPriceFrom", r.FilterPriceFrom))
	s += omitempty("PriceTo", r.FilterPriceTo)
	s += omitempty("SumFrom", r.FilterSumFrom)
	s += omitempty("SumTo", r.FilterSumTo)
	s += omitempty("Anum", r.FilterPnum)
	s += omitempty("Cnum", r.FilterRegnum)

	if r.FilterStrictFlag {
		s += omitempty("FilterStrictFlag", "on")
	}
	if page > 0 {
		s += omitempty("filterPage", page)
	}
	s += omitempty("perPage", r.PerPage)
	if xls {
		s += omitempty("xls", "short")
	}
	return s
}

// ContractFilters describes available filters for Contracts
type ContractFilters struct {
	demo               bool
	filterDateFrom     string
	filterDateTo       string
	filterDateExecFrom string
	filterDateExecTo   string
	filterCustomer     []string
	filterSupplier     []string
	filterCode         []string
	filterCode2        []string
	filterRegion       []string
	filterPnum         []string
	filterRegnum       []string
	filterCustomerN    []uint64
	filterCustomerINNn []uint64
	filterCustomerKPPn []uint64
	filterCodeN        []uint64
	filterCode2N       []uint64
	filterRegionN      []uint64
	filterPnumN        []uint64
	filterRegnumN      []uint64
	filterPriceFrom    float64
	filterPriceTo      float64
	filterSumFrom      float64
	filterSumTo        float64
	filterStrictFlag   bool
	filterPage         int
	perPage            int
}

func (f *ContractFilters) setFilters(params url.Values, db *sql.DB) {
	f.filterStrictFlag = params.Get("FilterStrictFlag") == "on"

	f.filterDateFrom = params.Get("DateFrom")
	f.filterDateTo = params.Get("DateTo")

	f.filterDateExecFrom = params.Get("DateEFrom")
	f.filterDateExecTo = params.Get("DateETo")

	f.filterPriceFrom, _ = strconv.ParseFloat(params.Get("PriceFrom"), 32)
	f.filterPriceTo, _ = strconv.ParseFloat(params.Get("PriceTo"), 32)
	f.filterSumFrom, _ = strconv.ParseFloat(params.Get("SumFrom"), 32)
	f.filterSumTo, _ = strconv.ParseFloat(params.Get("SumTo"), 32)

	if params["Customer"] != nil {
		f.filterCustomer = params["Customer"]
	}
	if params["Supplier"] != nil {
		f.filterSupplier = params["Supplier"]
	}
	if params["OKPDCode"] != nil {
		f.filterCode = params["OKPDCode"]
	}
	if params["OKPDCode2"] != nil {
		f.filterCode2 = params["OKPDCode2"]
	}
	if params["Region"] != nil {
		f.filterRegion = params["Region"]
	}
	tmp := params.Get("Cnum")
	tmp = strings.Replace(tmp, ", ", ",", -1)
	f.filterRegnum = strings.Split(tmp, ",")
	tmp = params.Get("Anum")
	tmp = strings.Replace(tmp, ", ", ",", -1)
	f.filterPnum = strings.Split(tmp, ",")
	f.filterPage, _ = strconv.Atoi(params.Get("filterPage"))
	f.filterPage = abs(f.filterPage)
	f.perPage, _ = strconv.Atoi(params.Get("perPage"))
	f.perPage = abs(f.perPage)
	if f.filterPage != 0 {
		f.filterPage -= 1
	}
	if f.perPage > 100000 {
		f.perPage = 100000
	}
	if f.perPage == 0 {
		f.perPage = perPageDefault
	}
	if params.Get("xls") != "" {
		f.filterPage = -1
	}
	if f.filterDateFrom == "" {
		f.filterDateFrom = "01.01.2014"
	}

	if f.filterDateTo == "" {
		f.filterDateTo = "31.12.2018" // ToDo - set to current Date
	}
	/*if f.filterDateExecFrom == "" {
		f.filterDateExecFrom = "01.01.2014"
	}

	if f.filterDateExecTo == "" {
		f.filterDateExecTo = "31.12.2018" // ToDo - set to current Date
	}*/

	if f.demo {
		f.filterDateFrom = "01.01.2017"
		f.filterDateTo = "31.12.2017"
		/*	f.filterDateExecFrom = "01.07.2014"
			f.filterDateExecTo = "30.09.2014"
		*/
	}
	if f.filterRegion != nil && len(f.filterRegion) != 0 {
		f.filterRegionN = numFilterN(f.filterRegion)
	}
	if f.filterCode != nil && len(f.filterCode) != 0 {
		okpds := api.StringIdsToOKPD1(f.filterCode, db)
		for _, c := range append([]api.OKPD1(nil), okpds...) {
			okpds = append(okpds, c.GetChildren(db)...)
		}
		f.filterCodeN = api.OKPD1toIds(okpds)
		log.Println("OKPD1 was set:", f.filterCodeN)
	}
	if f.filterCode2 != nil && len(f.filterCode2) != 0 {
		log.Println("")
		okpds := api.StringIdsToOKPD2(f.filterCode2, db)
		for _, c := range append([]api.OKPD2(nil), okpds...) {
			okpds = append(okpds, c.GetChildren(db)...)
		}
		f.filterCode2N = api.OKPD2toIds(okpds)
		log.Println("OKPD2 was set:", f.filterCode2N)
	}
	if f.filterCustomer != nil && len(f.filterCustomer) != 0 {
		//f.filterCustomerINNn, f.filterCustomerKPPn = numFilterNCustomers(f.filterCustomer)
		f.filterCustomerN = numFilterN(f.filterCustomer)
	}
	if f.filterRegnum != nil && len(f.filterRegnum) != 0 {
		f.filterRegnumN = numFilterN(f.filterRegnum)
	}
	if f.filterPnum != nil && len(f.filterPnum) != 0 {
		f.filterPnumN = numFilterN(f.filterPnum)
	}
}

func (f *ContractFilters) ApplyFilters(sc *sphinx.Client, db *sql.DB) {
	sc.SetFilterRange("Signed", filterDateTimestamp(f.filterDateFrom), filterToDateTimestamp(f.filterDateTo), false)
	log.Println("Signed: ", f.filterDateFrom, filterDateTimestamp(f.filterDateFrom), " to ", f.filterDateTo, filterToDateTimestamp(f.filterDateTo))

	log.Println("ExecDate tests:", filterDateTimestamp(f.filterDateExecFrom), " - ", filterToDateTimestamp(f.filterDateExecTo))

	if f.filterDateExecFrom != "" && f.filterDateExecTo == "" {
		sc.SetFilterRange("Exec", filterDateTimestamp(f.filterDateExecFrom), filterToDateTimestamp("31.12.2019"), false)
		log.Println("Exec 1:", filterDateTimestamp(f.filterDateExecFrom), " to ", filterToDateTimestamp("31.12.2019"))
	}
	if f.filterDateExecFrom == "" && f.filterDateExecTo != "" {
		sc.SetFilterRange("Exec", filterDateTimestamp("01.01.2014"), filterToDateTimestamp(f.filterDateExecTo), false)
		log.Println("Exec 2:", filterDateTimestamp("01.01.2014"), " to ", filterToDateTimestamp(f.filterDateExecTo))
	}
	if f.filterDateExecFrom != "" && f.filterDateExecTo != "" {
		sc.SetFilterRange("Exec", filterDateTimestamp(f.filterDateExecFrom), filterToDateTimestamp(f.filterDateExecTo), false)
		log.Println("Exec 3:", filterDateTimestamp(f.filterDateExecFrom), " to ", filterToDateTimestamp(f.filterDateExecTo))
	}

	if f.filterPriceTo != 0 {
		sc.SetFilterFloatRange("Price", float32(f.filterPriceFrom), float32(f.filterPriceTo), false)
		log.Println("Filter - Price:", float32(f.filterPriceFrom), float32(f.filterPriceTo), sc.GetLastError())
	} else {
		if f.filterPriceFrom != 0 {
			sc.SetFilterFloatRange("Price", 0.0, float32(f.filterPriceFrom), true)
		}
	}
	if f.filterSumTo != 0 {
		sc.SetFilterFloatRange("Sum", float32(f.filterSumFrom), float32(f.filterSumTo), false)
		log.Println("Filter - Sum:", float32(f.filterSumFrom), "-", float32(f.filterSumTo), sc.GetLastError())
	} else {
		if f.filterSumFrom != 0 {
			sc.SetFilterFloatRange("Sum", 0.0, float32(f.filterSumFrom), true)
		}
	}

	/*if f.filterCustomerINNn != nil && len(f.filterCustomerINNn) != 0 && f.filterCustomerKPPn != nil && len(f.filterCustomerKPPn) != 0 {
		sc.SetFilter("CINNn", f.filterCustomerINNn, false)
		sc.SetFilter("CKPPn", f.filterCustomerKPPn, false)
		log.Println("Filter - Customers:", f.filterCustomer, f.filterCustomerINNn, f.filterCustomerKPPn, sc.GetLastError())
	}*/
	if f.filterCustomerN != nil && len(f.filterCustomerN) != 0 {
		f.filterCustomerN = hackMKDC(f.filterCustomerN)
		sc.SetFilter("orgCode", f.filterCustomerN, false)
		log.Println("Filter - Customers:", f.filterCustomer, sc.GetLastError())
	} /*
		if f.filterCustomerN != nil && len(f.filterCustomerN) != 0 {
			f.filterCustomerINNn, f.filterCustomerKPPn = getOrgsInfo(f.filterCustomerN, db)
			sc.SetFilter("cinnn", f.filterCustomerINNn, false)
			sc.SetFilter("ckppn", f.filterCustomerKPPn, false)
			log.Println("Filter - Customers inn:", f.filterCustomer, sc.GetLastError())
		}*/

	if f.filterCodeN != nil && len(f.filterCodeN) != 0 {
		sc.SetFilter("okpdcode", f.filterCodeN, false)
		log.Println("Filter - OKPD:", f.filterCodeN, sc.GetLastError())
	}
	if f.filterCode2N != nil && len(f.filterCode2N) != 0 {
		sc.SetFilter("okpdcode", f.filterCode2N, false)
		log.Println("Filter - OKPD2:", f.filterCode2N, sc.GetLastError())
	}

	if f.filterRegionN != nil && len(f.filterRegionN) != 0 {
		sc.SetFilter("Region", f.filterRegionN, false)
		log.Println("Filter - region:", f.filterRegionN, sc.GetLastError())
	}

	if f.filterRegnumN != nil && len(f.filterRegnumN) != 0 {
		sc.SetFilter("rnum", f.filterRegnumN, false)
		log.Println("Filter - regnum:", f.filterRegnumN, sc.GetLastError())
	}

	if f.filterPnumN != nil && len(f.filterPnumN) != 0 {
		sc.SetFilter("pnum", f.filterPnumN, false)
		log.Println("Filter - pnum:", f.filterPnumN, sc.GetLastError())
	}

}

func numFilter(filter string) []uint64 {
	var array []uint64
	log.Println("Num FILTER EBA", len(filter))
	if filter != "" {
		s := strings.Split(filter, " ")
		log.Println("NUMF#1", s)
		for _, v := range s {
			tval, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				log.Println(err)
			}
			array = append(array, tval)
		}
		log.Println("NUM#2", array)
	}
	return array
}

/*(func numFilterNCustomers(filter []string) (inns, kpps []uint64) {
	var array []uint64
	log.Println("NUM#1", filter)
	for _, v := range filter {
		if v != "" {
			a := strings.Split(v, ";")
			log.Println("FC:", a)
			if len(a) == 2 {
				inn, err := strconv.ParseUint(a[0], 10, 64)
				if err != nil {
					log.Println(err)
				}
				kpp, err := strconv.ParseUint(a[1], 10, 64)

				if err != nil {
					log.Println(err)
				}
				inns = append(inns, inn)
				kpps = append(kpps, kpp)
			}
		}
	}
	log.Println("NUM#2", array)
	return
}*/
func numFilterN(filter []string) []uint64 {
	var array []uint64
	log.Println("NUM#1", filter)
	for _, v := range filter {
		if v != "" {
			tval, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				log.Println(err)
			}
			array = append(array, tval)
		}
	}
	log.Println("NUM#2", array)
	return array
}
func urlEscapeStringArray(array []string, prefix string) string {
	var s string
	for i, v := range array {
		if i != 0 {
			s += "&" + prefix + "=" + url.QueryEscape(v)
		} else {
			s += url.QueryEscape(v)
		}
	}
	return s
}

// contractsHandler - describes main processing for contracts
func contractsHandler(w http.ResponseWriter, r *http.Request, db, dbs *sql.DB, shd string, testmode bool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, lvl := accessLVL(GetEmail(r), 0, dbs)
	log.Println("IS TEST?", testmode)
	//params := map[string][]string{}
	params := url.Values{}
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	if r.Method == "GET" {
		//params = r.URL.Query()
		r.ParseForm()
		params = r.Form
	}
	log.Println("Params:", params)
	var query string
	var filters ContractFilters
	var result ContractResult
	result.Demo = lvl == 0
	result.Lvl = lvl

	filters.demo = result.Demo
	if result.Demo {
		result.FilterDateFrom = "01.01.2017"
		result.FilterDateTo = "31.12.2017"
	}
	result.NavType = "contracts"
	if len(params) != 0 {
		log.Println("Params Debug:", params)
		query = params.Get("query")
		filters.setFilters(params, db)
		result = contractsSearch(db, shd, query, &filters)
		result.IsSearch = true
		_, result.Lvl = accessLVL(GetEmail(r), 0, dbs)
		result.Fields = fields.FilterByLvl(fields.ContractFields, result.Lvl)
		result.Email = GetEmail(r)
		result.FilterStrictFlag = filters.filterStrictFlag
		//log.Println("?! FileStrictFlag:", result.FilterStrictFlag)
		decreaseRequestsCount(GetEmail(r), dbs)
		LogQuery(GetEmail(r), getIP(r), result.Query, GenURLContracts(result, params.Get("xls")), result.Total, params.Get("xls") == "short", dbs)
		if r.Method == "POST" {
			w.Header().Set("Content-type", "application/json")
			t, _ := json.Marshal(&result)
			fmt.Fprintf(w, "%s", string(t))
		}

		if r.Method == "GET" {
			if params.Get("xls") == "short" {
				_, result.Lvl = accessLVL(GetEmail(r), 0, dbs) // нахер второй раз?

				inputFields := params.Get("fields")

				if len(inputFields) != 0 {
					result.Fields = fields.FilterByIds(result.Fields, strings.Split(inputFields, ","))
				}

				log.Println("XLS LVL TEST:", result.Lvl)
				w.Header().Set("Content-type", "application/vnd.ms-excel")
				w.Header().Set("Content-Disposition", "attachment; filename=Контракты.xls")
				tmpl, err := template.New("").Delims("{[", "]}").ParseFiles("templates/contracts.xls.html")
				if err != nil {
					log.Println(err)
				}
				//tmpl := template.Must(template.New("cxls").Parse(ContractExcel))
				err = tmpl.ExecuteTemplate(w, "contracts.xls.html", result)
				if err != nil {
					log.Println(err)
				}
			} else {
				//tmpl := template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/loading-info.html","templates/contracts.html"))
				//tmpl.ExecuteTemplate(w, "contracts", result)
				tmpl, _ := template.ParseFiles("templates/contracts.html")
				tmpl.Execute(w, result)
			}
		}
	} else {
		if r.Method == "POST" {

			w.Header().Set("Content-type", "application/json")
			t, _ := json.Marshal(&result)
			fmt.Fprintf(w, "%s", string(t))
		}
		if r.Method == "GET" {
			//var tmpl *template.Template
			var err error
			//tmpl = template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html","templates/loading-info.html", "templates/contracts.html"))

			tmpl, err := template.ParseFiles("templates/contracts.html")
			tmpl.Execute(w, result)
			//err = tmpl.ExecuteTemplate(w, "contracts", result)
			if err != nil {
				log.Println(err)
			}
		}
	}

	//	fmt.Fprintf(w, "%s \n %s", query, parseQuery(query))
}

func sphinxCSearch(shard, q string, f ContractFilters, db *sql.DB) ([]int, int) {

	sc := sphinx.NewClient().SetServer(sphinxSwitch.GetHost(), 0).SetConnectTimeout(1000).SetLimits(0, 1000000, 1000000, 0)
	
	if err := sc.Error(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected with Sphinx! Shard:", shard)
	log.Println("SPew:", f.filterPriceFrom)
	f.ApplyFilters(sc, db)

	// log.Println(q)
	log.Println("Sphinx trying... ", shard)
	
	benchStart := time.Now()
	res, err := sc.Query(q, shard, shard)
	
	if err != nil {
		log.Println("Query error:", err)
	}

	//	log.Println(shd, " found: ", len(res.Matches), "/", res.Total, "|", res.Matches)
	var reslen int
	if res == nil {
		log.Println("MEMTEST:", res == nil)
		reslen = 0
		ids := make([]int, reslen)
		return ids, 0
	} else {
		reslen = len(res.Matches)
	}
	ids := make([]int, reslen)
	//log.Println(res)
	log.Println(len(res.Matches), "/", res.Total)
	for i, match := range res.Matches {
		ids[i] = int(match.DocId)
		//log.Println(match.DocId)
	}
	//log.Println(ids)
	log.Println(time.Since(benchStart))
	return ids, res.Total
}

func contractsSearch(db *sql.DB, shdSTR, query string, f *ContractFilters) ContractResult {

	// Prepare shards

	//shard := shards(f.filterDateFrom, f.filterDateTo, shdSTR)
	shard := shards.GenRangeByPattern("star_products_%s")
	//shard := shards(filterDateFrom, filterDateTo, "products")
	q := query

	if f.filterSupplier != nil && len(f.filterSupplier) != 0 {
		q += " " + strings.Join(f.filterSupplier, " |")
	}
	
	log.Println("Shards:", shard)

	// log.Println("Query before parse: ", q)
	//q = parseQuery(q)
	log.Println("!! isStrict?", f.filterStrictFlag)
	q = uquery.ParseQuery(q, f.filterStrictFlag)

	var ans ContractResult
	var pages []Bag
	//var ids []int

	for _, shd := range shard {
		ids, _ := sphinxCSearch(shd, q, *f, db)
		for _, id := range ids {
			var bag Bag
			bag.id = id
			bag.shard = shd
			pages = append(pages, bag)
		}
	}
	//var sh string
	if f.filterPage == -1 {
		f.filterPage = 0
		f.perPage = MaxDocs
	}
	log.Println("Pages debug:", f.filterPage, ";", f.perPage, ";", f.perPage*(f.filterPage+1), ";", len(pages))
	ans.Total = len(pages)

	for i := f.perPage * f.filterPage; i < ans.Total; i++ {
		if i >= f.perPage * (f.filterPage + 1) {
			break
		}
		log.Println("Overflow: ", i, "/", len(pages))
		r, err := retrieveContracts2(pages[i], db, query)
		if err != nil {
			log.Println(err)
		}
		ans.Data = concatContract(ans.Data, r)
	}
	/*
	for i := f.perPage * f.filterPage; i <= f.perPage*(f.filterPage+1); i++ {
		if i >= len(pages) {
			log.Println("Overflow: ", i, "/", len(pages))
			r, err := retrieveContracts(ids, db, sh, query)
			if err != nil {
				log.Println(err)
			}
			ans.Data = concatContract(ans.Data, r)
			break
		}
		//log.Println("Item:", i, i >= perPage*(filterPage+1))
		if i >= f.perPage*(f.filterPage+1) {
			r, err := retrieveContracts(ids, db, sh, query)
			if err != nil {
				log.Println(err)
			}
			ans.Data = concatContract(ans.Data, r)
			sh = pages[i].shard
			ids = nil
		}
		if sh == "" {
			sh = pages[i].shard
		} else {
			if sh != pages[i].shard {
				r, err := retrieveContracts(ids, db, sh, query)
				if err != nil {
					log.Println(err)
				}
				ans.Data = concatContract(ans.Data, r)
				sh = pages[i].shard
				ids = nil
			}
		}
		ids = append(ids, pages[i].id)
	}
	/**/

	ans.Count = len(ans.Data)
	// log.Println(ans.Data)
	log.Println(ans.Count, "/", ans.Total)
	log.Println("Pew-pew:", f.filterPriceFrom)
	ans.setMeta(query, *f)
	ans.FilterPriceFrom = f.filterPriceFrom
	ans.FilterSumFrom = f.filterSumFrom
	ans.FilterPriceTo = f.filterPriceTo
	ans.FilterSumTo = f.filterSumTo
	log.Println("Pew-n-Pew", ans.FilterPriceFrom)
	return ans
}

func (p *Product) HightlightQuery() string {
	/*
		var wordsArray []string
		words := p.Name
		specialChars, _ := regexp.Compile("[+*\\'\"«»]+")
		OR, _ := regexp.Compile("[|]")
		var hasOR bool
		if hasOR = OR.MatchString(p.Query); hasOR {
			log.Println("Hightlight - hasOR")
			wordsArray = strings.Split(p.Query, "|")
		} else {
			wordsArray = strings.Split(p.Query, " ")
		}
		//log.Println("EBATAS", hasOR, p.Query, wordsArray)
		for _, w := range wordsArray {
			w := specialChars.ReplaceAllLiteralString(w, "")
			re := regexp.MustCompile("(?i)" + w)
			words = re.ReplaceAllLiteralString(words, "<span style='background-color: #FFFF88;'>"+w+"</span>")
		}
		//return words
	*/
	return p.Name
}

func (p ContractResult) Pages() template.HTML {
	if p.PerPage == 0 {
		p.PerPage = int(perPageDefault)
	}

	log.Println("PAGES():", p.PerPage, p.FilterPage)
	totalPage := int(math.Ceil(float64(p.Total) / float64(p.PerPage)))

	startPage := p.FilterPage - 1;
	if startPage < 1 { startPage = 1 }
	endPage := startPage + 4
	if endPage > totalPage {
		endPage = totalPage
		startPage = endPage - 4
		if startPage < 1 { startPage = 1}
	}
	prevPage := p.FilterPage - 1
	if prevPage < 1 { prevPage = 1 }
	nextPage := p.FilterPage + 1
	if nextPage > totalPage { nextPage = totalPage }

	var templateD string
	if p.Total > p.PerPage {
		templateD += fmt.Sprintf(`
			<div class="row justify-content-md-center">
				<nav>
					<ul class="pagination">
						<li class="page-item"><a class="page-link" href="%s"><<</a></li>
						<li class="page-item"><a class="page-link" href="%s"><</a></li>`, p.genURL(false, 0), p.genURL(false, prevPage))

		for i := startPage; i <= endPage; i++ {
			if i == p.FilterPage + 1 {
				templateD += fmt.Sprintf(`
						<li class="page-item"><a class="page-link" style="color: #000;" href="%s">%d</a></li>`, p.genURL(false, i), i)
			} else {
				templateD += fmt.Sprintf(`
						<li class="page-item"><a class="page-link" href="%s">%d</a></li>`, p.genURL(false, i), i)
			}
		}
		templateD += fmt.Sprintf(`
						<li class="page-item"><a class="page-link" href="%s">></a></li>
						<li class="page-item"><a class="page-link" href="%s">>></a></li>
					</ul>
				</nav>
			<div class="row justify-content-md-center">`, p.genURL(false, nextPage), p.genURL(false, totalPage))
	}

	return template.HTML(templateD)
}

func (p ContractResult) ExcelLink() string {
	return p.genURL(true, 0)
	//return fmt.Sprintf("/contracts/?query=%s&filterDateFrom=%s&filterDateTo=%s&filterPriceFrom=%2f&filterPriceTo=%2f&filterSumFrom=%2f&filterSumTo=%2f&filterCode=%s&filterRegion=%s&filterCustomer=%s&filterSupplier=%s&xls=short", url.QueryEscape(p.Query), url.QueryEscape(p.FilterDateFrom), url.QueryEscape(p.FilterDateTo), p.FilterPriceFrom, p.FilterPriceTo, p.FilterSumFrom, p.FilterSumTo, urlEscapeStringArray(p.FilterCode, "filterCode"), urlEscapeStringArray(p.FilterRegion, "filterRegion"), urlEscapeStringArray(p.FilterCustomer, "filterCustomer"), url.QueryEscape(p.FilterSupplier))
}

func concatContract(old1, old2 []Product) []Product {
	newslice := make([]Product, len(old1)+len(old2))
	copy(newslice, old1)
	copy(newslice[len(old1):], old2)
	return newslice
}

func (p *Product) Check() bool {
	if p.RegNum == "" {
		return false
	}
	return true
}

func (p *Product) ExecDate() string {
	return fmt.Sprintf("%02d.%02d.%d", p.Exec.Day(), p.Exec.Month(), p.Exec.Year())
}

func (p *Product) SignedDate() string {
	return fmt.Sprintf("%02d.%02d.%d", p.Signed.Day(), p.Signed.Month(), p.Signed.Year())
}

func (p *Product) SignedXlsDate() string {
	return fmt.Sprintf("%02d.%02d.%d", p.Signed.Day(), p.Signed.Month(), p.Signed.Year())
}

func (p *Product) PriceFormat() string {
	price := fmt.Sprintf("%.2f", p.Price)
	price = strings.Replace(price, ".", ",", -1)
	return price
}

func (p *Product) PaidFormat() string {
	price := fmt.Sprintf("%.2f", p.Paid.Float64)
	price = strings.Replace(price, ".", ",", -1)
	return price
}

func (p *Product) PriceUSDFormat() string {
	rate := cbrate.GetCurrencyRate(p.SignedDate(), USDRate)
	if rate == 0 {
		return ""
	}
	price := fmt.Sprintf("%.2f", p.Price/rate)
	price = strings.Replace(price, ".", ",", -1)
	return price
}

func (p *Product) NMCKFormat() string {
	if !p.NMCK.Valid {
		return ""
	}
	t, _ := p.NMCK.Value()
	price := fmt.Sprintf("%.2f", t.(float64))
	price = strings.Replace(price, ".", ",", -1)
	//log.Println("NMCKTEST:", price)
	if price == "0,00" {
		return ""
	}
	return price
}

func (p *Product) NMCKUSDFormat() string {
	if !p.NMCK.Valid {
		log.Println("Rate is not Valid!")
		return ""
	}
	//log.Println("USDRate len:", len(USDRate), p.SignedDate(), cbrate.GetCurrencyRate(p.SignedDate(), USDRate))
	rate := cbrate.GetCurrencyRate(p.SignedDate(), USDRate)
	if rate == 0 {
		log.Println("Rate is zero!")
		return ""
	}
	t, _ := p.NMCK.Value()
	//log.Println("Rate:", rate, "RUR price:", t.(float64))

	price := fmt.Sprintf("%.2f", t.(float64)/rate)
	price = strings.Replace(price, ".", ",", -1)
	//log.Println("NMCKTEST:", price)
	if price == "0,00" {
		return ""
	}
	return price
}

func (p *Product) SumFormat() string {
	price := fmt.Sprintf("%.2f", p.Sum)
	price = strings.Replace(price, ".", ",", -1)
	return price
}

func (p *Product) SumUSDFormat() string {
	rate := cbrate.GetCurrencyRate(p.SignedDate(), USDRate)
	price := fmt.Sprintf("%.2f", p.Sum/rate)
	price = strings.Replace(price, ".", ",", -1)
	return price
}

func (p *Product) TotalSumFormat() string {

	price := fmt.Sprintf("%.3f", p.TotalSum)
	price = strings.Replace(price, ".", ",", -1)
	return price
}
func (p *Product) TotalSumUSDFormat() string {
	rate := cbrate.GetCurrencyRate(p.SignedDate(), USDRate)
	price := fmt.Sprintf("%.3f", p.TotalSum/rate)
	price = strings.Replace(price, ".", ",", -1)
	return price
}
func (p *Product) QuantityHint() string {
	hint := fmt.Sprintf("Сумма: %.2f руб.", p.Price * p.Quantity)
	hint = strings.Replace(hint, ".", ",", -1)
	return hint
}
func (p *Product) CustomerRegionGet() string {
	//log.Println("Trying region:", p.RegNum, p.CustomerINN[:2])
	return GetRegion(p.CustomerINN)
}

func retrieveContracts2(page Bag, db *sql.DB, queryst string) ([]Product, error) {

	shard := strings.Replace(page.shard, "star_", "", -1)
	data := new([]Product)
	table_num := strings.Split(shard, "_")[1]

	log.Println("SHARD", shard)
	query := fmt.Sprintf(`
		SELECT regexp_replace(regexp_replace(p.name,'\r','','g'),'\n','','g'),
			p.price as Price,
			p.quantity as Quantity,
			p.Units as Unit,
			p.sum as Sum,
			p.rate as Rate,
			c.price as TotalSum,
			p.okpd as OKPD,
			p.okpdinfo as OKPDInfo,
			p.customer as Customer,
			p.customerINN as CustomerINN,
			p.customerKPP as CustomerKPP,
			p.customerRegion as CustomerRegion,
			p.signed as Signed,
			p.exec as Exec,
			p.supplier as Supplier,
			p.supplierINN as SupplierINN,
			p.supplierKPP as SupplierKPP,
			p.supplierAddress as SupplierAddress,
			p.supplierPhone as SupplierPhone,
			p.Budget as Budget,
			p.BudgetSource as BudgetSource,
			c.Status as Status,
			p.regnum as RegNum,
			p.number as Number,
			p.purchasenumber as PurchaseNumber,
			p.placingway as Placing,
			c.Paid as Paid,
			c.NMCK as NMCK,
			c.ntname as ObjectInfo
		FROM products_%s as p, contracts_%s as c
		WHERE p.regnum = c.regnum and p.version in (select max(version) from products_%s where regnum=p.regnum)
			AND p.row in (%d) ORDER BY RegNum;`, table_num, table_num, table_num, page.id,
	)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	j := 0
	for rows.Next() {
		var row Product
		err := rows.Scan(
			&row.Name,
			&row.Price,
			&row.Quantity,
			&row.Unit,
			&row.Sum,
			&row.Rate,
			&row.TotalSum,
			&row.OKPD,
			&row.OKPDInfo,
			&row.Customer,
			&row.CustomerINN,
			&row.CustomerKPP,
			&row.CustomerRegion,
			&row.Signed,
			&row.Exec,
			&row.Supplier,
			&row.SupplierINN,
			&row.SupplierKPP,
			&row.SupplierAddress,
			&row.SupplierPhone,
			&row.Budget,
			&row.BudgetSource,
			&row.Status,
			&row.RegNum,
			&row.Number,
			&row.PurchaseNumber,
			&row.Placing,
			&row.Paid,
			&row.NMCK,
			&row.ObjectInfo,
		)
		row.Query = queryst
		*data = append(*data, row)
		if err != nil {
			log.Fatal(err)
		}
		if row.Name != "" {
			j++
		} else {
			log.Println("Error!")
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("!shard:", *data)
	return *data, err
}

func retrieveContracts(ids []int, db *sql.DB, shard, queryst string) ([]Product, error) {

	shard = strings.Replace(shard, "star_", "", -1)

	if len(ids) == 0 {
		return nil, nil
	}
	//data := make([]Product, len(ids))
	query := strings.Replace(`
								SELECT regexp_replace(regexp_replace(p.name,'\r','','g'),'\n','','g'),
												p.price as Price,
												p.quantity as Quantity,
												p.Units as Unit,
												p.sum as Sum,
												p.rate as Rate,
												c.price as TotalSum,
												p.okpd as OKPD,
												p.okpdinfo as OKPDInfo,
												p.customer as Customer,
												p.customerINN as CustomerINN,
												p.customerKPP as CustomerKPP,
												p.customerRegion as CustomerRegion,
												p.signed as Signed,
												p.exec as Exec,
												p.supplier as Supplier,
												p.supplierINN as SupplierINN,
												p.supplierKPP as SupplierKPP,
												p.supplierAddress as SupplierAddress,
												p.supplierPhone as SupplierPhone,
												p.Budget as Budget,
												p.BudgetSource as BudgetSource,
												c.Status as Status,
												p.regnum as RegNum,
												p.number as Number,
												p.purchasenumber as PurchaseNumber,
												p.placingway as Placing,
												c.Paid as Paid,
												c.NMCK as NMCK,
												c.ntname as ObjectInfo
												FROM products_%s as p, contracts_%s as c
												WHERE p.regnum = c.regnum and p.version in (select max(version) from products_%s where regnum=p.regnum)
												and p.row in (`, "%s", strings.Split(shard, "_")[1], -1)

	//WHERE c.regnum = p.regnum %r query = strings.Replace(query, "%r", whereRegion, -1)
	lenQ := strings.Replace(`SELECT         count(DISTINCT p.regnum)
                                                FROM products_%s as p, contracts_%s as c
                                                WHERE p.regnum = c.regnum and p.version in (select max(version) from products_%s where regnum=p.regnum)
                                                and p.row in (`, "%s", strings.Split(shard, "_")[1], -1)

	for l, i := range ids {
		if l != len(ids)-1 {
			query += strconv.Itoa(i) + ", "
			lenQ += strconv.Itoa(i) + ", "
		} else {
			//query += strconv.Itoa(i) + ") limit " + strconv.Itoa(len(ids)) + ";"
			query += strconv.Itoa(i) + ")"
			lenQ += strconv.Itoa(i) + ");"
		}
	}

	// log.Println(query)

	my_per := strings.Replace(`SELECT COUNT(*) FROM (%s) t;`, "%s", query, -1)
	query += " Order by RegNum;"
	lenQ = my_per

	var tl int
	err := db.QueryRow(lenQ).Scan(&tl)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("real count:", tl)
	if tl == 0 {
		return nil, nil
	}
	data := make([]Product, tl)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	j := 0

	for rows.Next() {
		err := rows.Scan(&data[j].Name,
			&data[j].Price,
			&data[j].Quantity,
			&data[j].Unit,
			&data[j].Sum,
			&data[j].Rate,
			&data[j].TotalSum,
			&data[j].OKPD,
			&data[j].OKPDInfo,
			&data[j].Customer,
			&data[j].CustomerINN,
			&data[j].CustomerKPP,
			&data[j].CustomerRegion,
			&data[j].Signed,
			&data[j].Exec,
			&data[j].Supplier,
			&data[j].SupplierINN,
			&data[j].SupplierKPP,
			&data[j].SupplierAddress,
			&data[j].SupplierPhone,
			&data[j].Budget,
			&data[j].BudgetSource,
			&data[j].Status,
			&data[j].RegNum,
			&data[j].Number,
			&data[j].PurchaseNumber,
			&data[j].Placing,
			&data[j].Paid,
			&data[j].NMCK,
			&data[j].ObjectInfo,
		)
		data[j].Query = queryst

		if err != nil {
			log.Fatal(err)
		}
		//log.Println(j, ids[j])
		if data[j].Name != "" {
			j++
		} else {
			log.Println("Error!")
		}

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(data)
	return data, err
}

func (p ContractResult) FilterDateFromFormat() string {
	return strings.Replace(p.FilterDateFrom, ".", "/", -1)
}
func (p ContractResult) FilterDateToFormat() string {
	return strings.Replace(p.FilterDateTo, ".", "/", -1)
}

func (p ContractResult) FieldWithIdExist(id string) bool {
	return fields.FieldWithIdExist(p.Fields, id)
}
