package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"html/template"
	fields "bitbucket.org/company-one/tender-one/document-fields"
	"bitbucket.org/company-one/tender-one/sphinx-switch"
	_ "github.com/lib/pq"
	"time"
	"regexp"
	"github.com/yunge/sphinx"
	"net/http"
)

type Notification223Result struct {
	Email              string         `json:"email"`
	Data               []Notification223 `json:"data"`
	Total              int            `json:"total"`
	Count              int            `json:"count"`
	Query              string         `json:"query"`
	FilterDateFrom     string         `json:"DateFrom"`
	FilterDateTo       string         `json:"DateTo"`
	FilterCustomer     []string       `json:"Customer"`
	FilterMaxPriceFrom float64        `json:"MaxPriceFrom"`
	FilterMaxPriceTo   float64        `json:"MaxPriceTo"`
	FilterOrg          []string       `json:"Org"`
	FilterRegion       []string       `json:"Region"`
	FilterPnum         string         `json:"Anum"`
	FilterRegnum       string         `json:"filterRegnum"`
	FilterPage         int            `json:"PageNum"`
	PerPage            int            `json:"perPage"`
	Fields             []fields.Field
	IsSearch           bool
	Demo               bool
	Lvl                int
	NavType            string
}



func (p *Notification223Result) isNotEmpty() bool {
	return p.Count != 0
}
func (r Notification223Result) NotificationsToTemplate() string {
	return "/templates/?save=" + url.QueryEscape(r.genURL(false, 0))
}
func (r *Notification223Result) genURL(xls bool, page int) string {
	s := "/auctions223/?query=" + url.QueryEscape(r.Query)
	s += omitempty("Org", r.FilterOrg)
	s += omitempty("DateFrom", r.FilterDateFrom)
	s += omitempty("DateTo", r.FilterDateTo)
	s += omitempty("MaxPriceFrom", r.FilterMaxPriceFrom)
	s += omitempty("MaxPriceTo", r.FilterMaxPriceTo)
	s += omitempty("Customer", r.FilterCustomer)
	s += omitempty("Region", r.FilterRegion)
	s += omitempty("Anum", r.FilterPnum)
	//s += omitempty("filterRegnum", r.FilterRegnum)
	if page > 0 {
		s += omitempty("PageNum", page)
	}
	s += omitempty("perPage", r.PerPage)
	if xls {
		s += omitempty("xls", "short")
	}
	return s
}
func (p Notification223Result) Pages() template.HTML {
	if p.PerPage == 0 {
		p.PerPage = perPageDefault
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

	/*
	pageNum := int(math.Ceil(float64(p.Total / p.PerPage)))
	for i := 1; i <= pageNum; i++ {
		templateD += fmt.Sprintf("<a href=\"%s\">%d</a> &nbsp;", p.genURL(false, i), i)
		//template += fmt.Sprintf("<a href='/notifications/?query=%s&filterDateFrom=%s&filterDateTo=%s&filterPriceFrom=%2f&filterPriceTo=%2f&filterCustomer=%s&filterCode=%s&filterRegion=%s&perPage=%d&filterPage=%d'>%d</a> &nbsp;", url.QueryEscape(p.Query), p.FilterDateFrom, p.FilterDateTo, p.FilterPriceFrom, p.FilterPriceTo, url.QueryEscape(p.FilterCustomer), url.QueryEscape(p.FilterCode), url.QueryEscape(p.FilterRegion), p.PerPage, i, i)
	}
	*/
	return template.HTML(templateD)
}
func (p Notification223Result) ExcelLink() string {
	return p.genURL(true, 0)
	//return fmt.Sprintf("/contracts/?query=%s&filterDateFrom=%s&filterDateTo=%s&filterPriceFrom=%2f&filterPriceTo=%2f&filterSumFrom=%2f&filterSumTo=%2f&filterCode=%s&filterRegion=%s&filterCustomer=%s&filterSupplier=%s&xls=short", url.QueryEscape(p.Query), url.QueryEscape(p.FilterDateFrom), url.QueryEscape(p.FilterDateTo), p.FilterPriceFrom, p.FilterPriceTo, p.FilterSumFrom, p.FilterSumTo, urlEscapeStringArray(p.FilterCode, "filterCode"), urlEscapeStringArray(p.FilterRegion, "filterRegion"), urlEscapeStringArray(p.FilterCustomer, "filterCustomer"), url.QueryEscape(p.FilterSupplier))
}

// Notification - general description for result
type Notification223 struct {
	Row            int				`json:"row"`
	RegNum		   string			`json:"RegNum"`
	Guid		   string			`json:"guid"`
	Name		   string
	Published      time.Time		`json:"published"`
	Created        time.Time		`json:"published"`
	Query          string
	Type           string			`json:"type"`
	Sum		       float64
	Lotname		   string
	Customer       sql.NullString
	CustomerINN    sql.NullString
	CustomerKPP    sql.NullString
	Placer	       sql.NullString
	PlacerINN      sql.NullString
	PlacerKPP      sql.NullString
	Status 		   sql.NullString
	Cancelled	   bool
	Region		   int64
	MethodName	   sql.NullString
	PlaceName	   sql.NullString
	PlaceUrl	   sql.NullString
	TotalSum	   float64
}

func (p *Notification223) PublishedDate() string {
	return fmt.Sprintf("%02d.%02d.%d", p.Published.Day(), p.Published.Month(), p.Published.Year())
}

func (p *Notification223) PublishedXlsDate() string {
	return fmt.Sprintf("%02d.%02d.%d", p.Published.Day(), p.Published.Month(), p.Published.Year())
}

func (p *Notification223) GetMaxPrice() string {
		return fmt.Sprintf("%.2f",p.Sum)
}
func (p *Notification223) PriceExcel() string {
	return p.GetMaxPrice()
}

func (p *Notification223) GetTotalPrice() string {
	return fmt.Sprintf("%.2f",p.TotalSum)
}
func (p *Notification223) TotalPriceExcel() string {
	return p.GetTotalPrice()
}

func (p *Notification223) HightlightQuery() string {
	words := p.Name
	for _, w := range strings.Split(p.Query, " ") {
		re := regexp.MustCompile("(?i)" + w)
		words = re.ReplaceAllLiteralString(words, "<span style='background-color: #FFFF88;'>"+w+"</span>")
	}
	//return words
	return p.Name
}


type Notification223Filters struct {
	filterDateFrom     string
	filterDateTo       string
	filterCustomer     []string
	filterOrg          []string
	filterCode         []string
	filterRegion       []string
	filterPnum         []string
	filterRegnum       []string
	filterCustomerN    []uint64
	filterOrgN         []uint64
	filterCodeN        []uint64
	filterRegionN      []uint64
	filterPnumN        []uint64
	filterDocsData     []uint64
	filterRegnumN      []uint64
	filterMaxPriceFrom float64
	filterMaxPriceTo   float64
	filterCanceled     bool
	filterFinished     bool
	filterPage         int
	filterDocsFlag     bool
	filterStrictFlag   bool
	perPage            int
	demo               bool
}

func (f *Notification223Filters) setFilters(params url.Values) {
	f.filterDateFrom = params.Get("DateFrom")
	f.filterDateTo = params.Get("DateTo")

	f.filterMaxPriceFrom, _ = strconv.ParseFloat(params.Get("MaxPriceFrom"), 32)
	f.filterMaxPriceTo, _ = strconv.ParseFloat(params.Get("MaxPriceTo"), 32)

	if params["Customer"] != nil {
		f.filterCustomer = params["Customer"]
	}
	if params["Org"] != nil {
		f.filterOrg = params["Org"]
	}
	if params["filterCode"] != nil {
		f.filterCode = params["filterCode"]
	}
	if params["Region"] != nil {
		f.filterRegion = params["Region"]
	}
	tmp := params.Get("filterRegnum")
	tmp = strings.Replace(tmp, ", ", ",", -1)
	f.filterRegnum = strings.Split(tmp, ",")
	tmp = params.Get("Anum")
	tmp = strings.Replace(tmp, ", ", ",", -1)
	f.filterPnum = strings.Split(tmp, ",")

	//f.filterPage, _ = strconv.ParseInt(params.Get("filterPage"), 10, 0)
	f.filterPage, _ =  strconv.Atoi(params.Get("PageNum"))
	f.filterPage = abs(f.filterPage)
	//f.perPage, _ = strconv.ParseInt(params.Get("perPage"), 10, 0)
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
	if f.demo {
		f.filterDateFrom = "01.01.2017"
		f.filterDateTo = "31.12.2017"
	}
	if f.filterRegion != nil && len(f.filterRegion) != 0 {
		f.filterRegionN = numFilterN(f.filterRegion)
	}
	if f.filterCode != nil && len(f.filterCode) != 0 {
		f.filterCodeN = numFilterN(f.filterCode)
	}
	if f.filterCustomer != nil && len(f.filterCustomer) != 0 {
		f.filterCustomerN = numFilterN(f.filterCustomer)
	}
	if f.filterOrg != nil && len(f.filterOrg) != 0 {
		f.filterOrgN = numFilterN(f.filterOrg)
	}
	if f.filterRegnum != nil && len(f.filterRegnum) != 0 {
		f.filterRegnumN = numFilterN(f.filterRegnum)
	}
	if f.filterPnum != nil && len(f.filterPnum) != 0 {
		f.filterPnumN = numFilterN(f.filterPnum)
	}
}

func (f *Notification223Filters) ApplyFilters(sc *sphinx.Client) {
	log.Println("Date filters: ", f.filterDateFrom, "-", f.filterDateTo)

	sc.SetFilterRange("published", filterDateTimestamp(f.filterDateFrom), filterToDateTimestamp(f.filterDateTo), false)

	if f.filterPnumN != nil && len(f.filterPnumN) != 0 {
		sc.SetFilter("regnum", f.filterPnumN, false) // regnum <=> pnum
		log.Println("Filter - pnum:", f.filterPnumN, sc.GetLastError())
	}
	if f.filterRegnumN != nil && len(f.filterRegnumN) != 0 {
		sc.SetFilter("rnum", f.filterRegnumN, false)
		log.Println("Filter - regnum:", f.filterRegnumN, sc.GetLastError())
	}
	if f.filterRegionN != nil && len(f.filterRegionN) != 0 {
		sc.SetFilter("region", f.filterRegionN, false)
		log.Println("Filter - region:", f.filterRegionN, sc.GetLastError())
	}
	if f.filterMaxPriceTo != 0 {
		sc.SetFilterFloatRange("sum", float32(f.filterMaxPriceFrom), float32(f.filterMaxPriceTo), false)
		log.Println("Filter - Price:", float32(f.filterMaxPriceFrom), float32(f.filterMaxPriceTo), sc.GetLastError())
	}
	if f.filterOrgN != nil && len(f.filterOrgN) != 0 {
		f.filterOrgN = hackMKDC(f.filterOrgN)
		sc.SetFilter("pinn", f.filterOrgN, false)
		log.Println("Filter - Orgs:", f.filterOrg, sc.GetLastError())
	}
	if f.filterCustomerN != nil && len(f.filterCustomerN) != 0 {
		f.filterCustomerN = hackMKDC(f.filterCustomerN)
		sc.SetFilter("cinn", f.filterCustomerN, false)
		log.Println("Filter - Customers:", f.filterCustomer, sc.GetLastError())
	}

}

func (result *Notification223Result) setMeta(query string, f Notification223Filters) {
	result.Query = query
	if f.demo {
		f.filterDateFrom = "01.01.2017"
		f.filterDateTo = "31.12.2017"
	}
	// log.Println(result.Query)
	result.FilterDateFrom = f.filterDateFrom
	// log.Println(result.FilterDateFrom)
	result.FilterDateTo = f.filterDateTo
	// log.Println(result.FilterDateTo)
	result.FilterMaxPriceFrom = float64(f.filterMaxPriceFrom)
	result.FilterMaxPriceTo = float64(f.filterMaxPriceTo)
	result.FilterCustomer = f.filterCustomer
	result.FilterOrg = f.filterOrg
	result.FilterRegion = f.filterRegion
	result.FilterPnum = strings.Join(f.filterPnum, ",")
	result.FilterRegnum = strings.Join(f.filterRegnum, ",")
	result.FilterPage = f.filterPage
	result.PerPage = int(f.perPage)
	result.NavType = "notifications223"
}



func notifications223Handler(w http.ResponseWriter, r *http.Request, db, dbs *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//params := map[string][]string{}
	_, lvl := accessLVL(GetEmail(r), 0, dbs)
	params := url.Values{}
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	if r.Method == "GET" {
		params = r.URL.Query()
	}
	var query string
	var result Notification223Result
	result.Demo = lvl == 0
	result.Lvl = lvl
	var filters Notification223Filters
	filters.demo = result.Demo
	if result.Demo {
		result.FilterDateFrom = "01.01.2017"
		result.FilterDateTo = "31.12.2017"
	}
	result.NavType = "notifications223"
	if len(params) != 0 {
		log.Println("Params Debug:", params)
		query = params.Get("query")
		filters.setFilters(params)

		result = notifications223Search(db, query, filters)
		result.IsSearch = true
		result.Email = GetEmail(r)
		_, result.Lvl = accessLVL(GetEmail(r), 0, dbs)
		result.Fields = fields.FilterByLvl(fields.Notification223Fields, result.Lvl)
		decreaseRequestsCount(GetEmail(r), dbs)
		LogQuery(GetEmail(r), getIP(r), result.Query, GenURLNotifications223(result, params.Get("xls")), result.Total, params.Get("xls") == "short", dbs)

		if r.Method == "POST" {
			w.Header().Set("Content-type", "application/json")
			t, _ := json.Marshal(&result)

			fmt.Fprintf(w, "%s", string(t))
		}
		if r.Method == "GET" {
			if params.Get("xls") == "short" {

				inputFields := params.Get("fields")

				if len(inputFields) != 0 {
					result.Fields = fields.FilterByIds(result.Fields, strings.Split(inputFields, ","))
				}

				w.Header().Set("Content-type", "application/vnd.ms-excel")
				w.Header().Set("Content-Disposition", "attachment; filename=223ФЗ.xls")
				tmpl, err := template.New("").Delims("{[", "]}").ParseFiles("templates/notifications223.xls.html")
				if err != nil {
					log.Fatal(err)
				}
				tmpl.ExecuteTemplate(w, "notifications223.xls.html", result)
			}
			//tmpl := template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/loading-info.html","templates/notifications223.html"))
			//tmpl := template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/tests.html"))
			tmpl, err := template.ParseFiles("templates/notifications223.html")
			tmpl.Execute(w, result)
			log.Println("RESULT_1: ", result, "DATE:", result.FilterDateFrom, result.FilterDateTo)

			log.Println("TEST LVL:", result.Lvl)
			//err := tmpl.ExecuteTemplate(w, "notifications223", result)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		if r.Method == "POST" {
			w.Header().Set("Content-type", "application/json")
			t, _ := json.Marshal(&result)
			fmt.Fprintf(w, "%s", string(t))
		}
		if r.Method == "GET" {
			//tmpl := template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/loading-info.html","templates/notifications223.html"))
			//tmpl := template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/tests.html"))
			tmpl, err := template.ParseFiles("templates/notifications223.html")
			tmpl.Execute(w, result)
			log.Println("RESULT_2: ", result, "DATE:", result.FilterDateFrom, result.FilterDateTo)

			//err := tmpl.ExecuteTemplate(w, "notifications223", result)
			if err != nil {
				log.Println(err)
			}
			// log.Println(result.Data)
		}
	}

	log.Println("RESULT_FIN: ", result, "DATE:", result.FilterDateFrom, result.FilterDateTo, result.FilterMaxPriceTo)
	//	fmt.Fprintf(w, "%s \n %s", query, parseQuery(query))
}

func sphinxN223Search(shard, q string, f Notification223Filters) ([]int, int) {

	var ids []int

	sc := sphinx.NewClient().SetServer(sphinxSwitch.GetHost(), 0).SetConnectTimeout(1000).SetLimits(0, 10000000, 1000000, 0)

	log.Println("Connected with Sphinx! Shard:", shard, shard[:4])
	log.Println("Date filters: ", f.filterDateFrom, "-", f.filterDateTo)
	log.Println(filterDateTimestamp(f.filterDateFrom),"-",filterToDateTimestamp(f.filterDateTo))

	sc.ResetFilters()
	sc.SetFilterRange("published", filterDateTimestamp(f.filterDateFrom), filterToDateTimestamp(f.filterDateTo), false)
	//if f.filterCodeN != nil && len(f.filterCodeN) != 0 {
	//	sc.SetFilter("okpd", f.filterCodeN, false)
	//	log.Println("Filter - OKPD:", f.filterCodeN, sc.GetLastError())
	//}
	//if f.filterDocsData != nil && len(f.filterDocsData) != 0 && f.filterDocsFlag != false {
	//	sc.SetFilter("pnum", f.filterDocsData, false)
	//	q = ""
	//	log.Println("Filter - pnum (Docs's data!):", len(f.filterDocsData), sc.GetLastError())
	//} else {
	log.Println("Checking pnumN filter:",f.filterPnumN)
	if f.filterPnumN != nil && len(f.filterPnumN) != 0 {
		sc.SetFilter("regnumN", f.filterPnumN, false)
		log.Println("Filter - pnum:", f.filterPnumN, sc.GetLastError())
	}
	//}
	if f.filterRegnumN != nil && len(f.filterRegnumN) != 0 {
		sc.SetFilter("rnum", f.filterRegnumN, false)
		log.Println("Filter - regnum:", f.filterRegnumN, sc.GetLastError())
	}
	if f.filterRegionN != nil && len(f.filterRegionN) != 0 {
		sc.SetFilter("region", f.filterRegionN, false)
		log.Println("Filter - region:", f.filterRegionN, sc.GetLastError())
	}
	if f.filterMaxPriceTo != 0 {
		sc.SetFilterFloatRange("sum", float32(f.filterMaxPriceFrom), float32(f.filterMaxPriceTo), false)
		log.Println("Filter - Price:", float32(f.filterMaxPriceFrom), float32(f.filterMaxPriceTo), sc.GetLastError())
	}
	//
	if f.filterCustomerN != nil && len(f.filterCustomerN) != 0 {
		sc.SetFilter("cinn", f.filterCustomerN, false)
	}

	//	if strings.Split(shard, "_")[1] != "ncustomers" {
	//		sc.SetFilter("org", f.filterCustomerN, false)
	//		log.Println("Filter - Orgs:", f.filterOrg, sc.GetLastError())
	//	}
	//	log.Println("Filter - Customers:", f.filterCustomer, sc.GetLastError())
	//} else {
	if f.filterOrgN != nil && len(f.filterOrgN) != 0 {
			sc.SetFilter("pinn", f.filterOrgN, false)
			log.Println("Filter - Orgs:", f.filterOrg, sc.GetLastError())
	}
	//}

	log.Println("The query: <", q, ">")
	res, err := sc.Query(q, "star_notifications223", "star_notifications223")
	//log.Println("Res::",res.Matches,"\n",res.TotalFound)

	if err != nil {
		log.Println("Query error:", err)
	}

	if err = sc.Error(); err != nil {
		log.Fatal(err)
	}
	//	log.Println(shd, " found: ", len(res.Matches), "/", res.Total, "|", res.Matches)
	if res == nil {
		//log.Println("Panic!", f.filterDocsData)
		log.Println("Panic! res is nil! Len of filter:", len(f.filterDocsData), "\n The query:", q, " and shard is ", shard, "\n and the filters are:", f)
		return ids, 0
	} else {
		ids = make([]int, len(res.Matches))
		log.Println("Post query: ", len(res.Matches), "/", res.Total)
		for i, match := range res.Matches {
			ids[i] = int(match.DocId)

			//log.Println("machet:", match.DocId, int(match.DocId), "; i:", i)
		}

		//log.Println(ids)
	}

	return ids, res.Total
}


func notifications223Search(db *sql.DB, query string, f Notification223Filters) Notification223Result {

	q := query

	// Prepare shards
	log.Println("Trying to generate shards", f.filterDateFrom, f.filterDateTo)
	//shard := shards(f.filterDateFrom, f.filterDateTo, "test_notifications")
	shard := []string{"star_notifications223"}
	// Above - temporary hack, because zero shard is off for now.

	var ans Notification223Result
	var pages []Bag
	//var ids []int
	//var notifications []Notification
	//var total int

	log.Println("Let's do real index search now")
	shd := "star_notifications223"
	//for _, shd := range shard {
	ids, _ := sphinxN223Search(shd, q, f)

	for _, id := range ids {
		var bag Bag
		bag.id = id
		bag.shard = shd
		pages = append(pages, bag)
	}
	//}
	log.Println("Huinya test:", len(f.filterDocsData) == len(ids), len(ids))
	//var sh string
	if f.filterPage == -1 {
		f.filterPage = 0
		f.perPage = MaxDocs
	}
	log.Println("Pages debug:", f.filterPage, ";", f.perPage, ";", f.perPage*(f.filterPage+1), ";", len(pages))
	ans.Total = len(pages)
	////////////////////////////////////////////////////////////////////

	data := new([]Notification223)

	if len(ids) > 0 && f.filterPage >= 0 {// && int64(ans.Total) >= f.perPage*(f.filterPage) {
		log.Println("SHARD", shard)
		get_page := `
			SELECT n.row
			  , n.regnum
			  , n.name
			  , n.customerName AS cName
			  , n.customerINN AS cINN
			  , n.customerKPP AS cKPP
			  , n.placerName AS pName
			  , n.placerINN AS pINN
			  , n.placerKPP AS pKPP
			  , SUBSTRING(n.customerKPP,1,2)::integer AS region
			  , timestamptz(n.created) AS created
			  , timestamptz(n.published) AS published
			  , n.methodname AS mName
			  , n.placename
			  , n.placeurl
			  , l.subject AS lotname
			  , l.initSum AS sum
			  , l.cancelled
			  , n.status
			  , SUM(l.initSum) OVER (PARTITION BY n.regnum) AS TotalSum
			FROM notifications223 AS n, lots223 AS l
			WHERE n.regnum = l.regnum
			  AND n.version = l.version
			  AND n.row IN (`

		for i, val := range ids[f.perPage * f.filterPage:] {
			//log.Println("forrrrr: ", i, val)
			if i == f.perPage { break; }

			if i == 0 {
				get_page += strconv.Itoa(val)
			} else {
				get_page += ", " + strconv.Itoa(val)
			}
		}
		get_page +=  ") ORDER BY published;" //fmt.Sprintf("OFFSET %d LIMIT %d;", f.perPage * f.filterPage,  f.perPage)

		//log.Println("BEGIN Querry: ", get_page, "END Querry")

		rows, err := db.Query(get_page)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		j := 0

		for rows.Next() {
			var row Notification223
			err := rows.Scan(
				&row.Row,
				&row.RegNum,
				&row.Name,
				&row.Customer,
				&row.CustomerINN,
				&row.CustomerKPP,
				&row.Placer,
				&row.PlacerINN,
				&row.PlacerKPP,
				&row.Region,
				&row.Created,
				&row.Published,
				&row.MethodName,
				&row.PlaceName,
				&row.PlaceUrl,
				&row.Lotname,
				&row.Sum,
				&row.Cancelled,
				&row.Status,
				&row.TotalSum,
			)
			row.Query = query
			*data = append(*data, row)
			if err != nil {
				log.Fatal(err)
			}
			j++
		}
		log.Println("PerPages_______JJJ:", j)
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	ans.Data = concatNotifications223(ans.Data, *data)

	////////////////////////////////////////////////////////////////////
	/*
	for i := f.perPage * f.filterPage; i <= f.perPage*(f.filterPage+1); i++ {
		if i >= int64(len(pages)) {
			log.Println("Overflow: ", i, "/", len(pages))
			r, err := retrieveNotifications223(ids, db, sh, query, f)
			if err != nil {
				log.Println(err)
			}
			ans.Data = concatNotifications223(ans.Data, r)
			break
		}
		//log.Println("Item:", i, i+1 >= perPage*(filterPage+1))
		if i+1 >= f.perPage*(f.filterPage+1) {
			r, err := retrieveNotifications223(ids, db, sh, query, f)
			if err != nil {
				log.Println(err)
			}
			ans.Data = concatNotifications223(ans.Data, r)
			sh = pages[i].shard
			ids = nil
		}
		if sh == "" {
			sh = pages[i].shard
		} else {
			if sh != pages[i].shard {
				r, err := retrieveNotifications223(ids, db, sh, query, f)
				if err != nil {
					log.Println(err)
				}
				ans.Data = concatNotifications223(ans.Data,r)
				sh = pages[i].shard
				ids = nil
			}
		}
		//log.Printf("SH: %s; ids: %d\n", pages[i].shard, len(ids))
		ids = append(ids, pages[i].id)
	}
	//*/
	ans.Count = len(ans.Data)
	ans.Query = query
	ans.FilterPage = f.filterPage + 1
	ans.setMeta(query, f)
	log.Println(ans.Count, "/", ans.Total)

	return ans
}

func concatNotifications223(old1, old2 []Notification223) []Notification223 {
	newslice := make([]Notification223, len(old1)+len(old2))
	copy(newslice, old1)
	copy(newslice[len(old1):], old2)
	return newslice
}

func retrieveNotifications223(ids []int, db *sql.DB, shard, queryst string, f Notification223Filters) ([]Notification223, error) {

	shard = strings.Replace(shard, "star_", "", -1)
	//shard = strings.Replace(shard, "ncustomers", "notifications_customers", -1)
	//shard = strings.Replace(shard, "test_c1_", "", -1)

	data := new([]Notification223)

	if len(ids) == 0 {
		return *data, nil
	}

	log.Println("SHARD", shard)
	query := `
		SELECT n.row
		  , n.regnum
		  , n.name
		  , n.customerName AS cName
		  , n.customerINN AS cINN
		  , n.customerKPP AS cKPP
		  , n.placerName AS pName
		  , n.placerINN AS pINN
		  , n.placerKPP AS pKPP
		  , SUBSTRING(n.customerKPP,1,2)::integer AS region
		  , timestamptz(n.created) AS created
		  , timestamptz(n.published) AS published
		  , n.methodname AS mName
		  , n.placename
		  , n.placeurl
		  , l.subject AS lotname
		  , l.initSum AS sum
		  , l.cancelled
		  , n.status
		  , SUM(l.initSum) OVER (PARTITION BY n.regnum) AS TotalSum
		FROM notifications223 AS n, lots223 AS l
		WHERE n.regnum = l.regnum
		  AND n.version = l.version
		  AND n.row IN (
	`

	for l, i := range ids {
		if l != len(ids)-1 {
			query += strconv.Itoa(i) + ", "
		} else {
			query += strconv.Itoa(i) + ") " + " ORDER BY published;"
		}
	}
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	j := 0

	for rows.Next() {
		var row Notification223
		err := rows.Scan(&row.Row,
			   			 &row.RegNum,
						 &row.Name,
			   			 &row.Customer,
						 &row.CustomerINN,
						 &row.CustomerKPP,
						 &row.Placer,
						 &row.PlacerINN,
			 			 &row.PlacerKPP,
						 &row.Region,
			             &row.Created,
						 &row.Published,
						 &row.MethodName,
						 &row.PlaceName,
						 &row.PlaceUrl,
						 &row.Lotname,
						 &row.Sum,
						 &row.Cancelled,
						 &row.Status,
						 &row.TotalSum,
		)
		row.Query = queryst
		*data = append(*data, row)
		if err != nil {
			log.Fatal(err)
		}
		j++
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return *data, err
}
func (p Notification223) FixPurchaseNumber() string {
	s := ""
	if len(p.RegNum) < 11 {
		for i := 0; i < 11-len(p.RegNum); i++ {
			s += "0"
		}

	}
	s += p.RegNum
	return s
}

func (p Notification223) RegionGet() string {
		return GetRegionFi(p.Region)
}

func (p Notification223Result) FieldWithIdExist(id string) bool {
	return fields.FieldWithIdExist(p.Fields, id)
}
