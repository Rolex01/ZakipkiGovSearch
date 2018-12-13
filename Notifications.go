package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"html/template"
	"time"

	shards "bitbucket.org/company-one/tender-one/shards"
	fields "bitbucket.org/company-one/tender-one/document-fields"
	uquery "bitbucket.org/company-one/tender-one/query"
	"bitbucket.org/company-one/tender-one/sphinx-switch"
	_ "github.com/lib/pq"
	"github.com/yunge/sphinx"
)

type NotificationResult struct {
	Email              string         `json:"email"`
	Data               []Notification `json:"data"`
	Total              int            `json:"total"`
	Count              int            `json:"count"`
	Query              string         `json:"query"`
	FilterDateFrom     string         `json:"DateFrom"`
	FilterDateTo       string         `json:"DateTo"`
	FilterMaxPriceFrom float64        `json:"MaxPriceFrom"`
	FilterMaxPriceTo   float64        `json:"MaxPriceTo"`
	FilterCustomer     []string       `json:"Customer"`
	FilterOrg          []string       `json:"Org"`
	FilterRegion       []string       `json:"Region"`
	FilterCode         []string       `json:"OKPDCode"`
	FilterPnum         string         `json:"Anum"`
	FilterRegnum       string         `json:"filterRegnum"`
	FilterPage         int            `json:"filterPage"`
	PerPage            int            `json:"perPage"`
	Fields             []fields.Field
	IsSearch           bool
	Demo               bool
	Lvl                int
	NavType            string
	FilterDocsFlag     bool
	FilterStrictFlag   bool
}

/*func WriteToExcel(r NotificationResult, w *io.Writer) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Выгрузка")
	// writing header
	r0 := sheet.AddRow()
	r0.AddCell().Value = "Номер извещения о проведении торгов"
	r1 := sheet.AddRow()
	r1.AddCell().Value = "Наименование аукциона"
	r2 := sheet.AddRow()
	r2.AddCell().Value = "Дата размещения"
	r3 := sheet.AddRow()
	r3.AddCell().Value = "НМЦК, руб."
	r4 := sheet.AddRow()
	r4.AddCell().Value = "Организация, осуществляющая закупку"
	r5 := sheet.AddRow()
	r5.AddCell().Value = "ИНН организации, осуществляющей заказ"
	r6 := sheet.AddRow()
	r6.AddCell().Value = "Заказчик"
	r7 := sheet.AddRow()
	r7.AddCell().Value = "ИНН заказчик"
	if r.Lvl >= 1 {
		r8 := sheet.AddRow()
		r8.AddCell().Value = "Номер реестровой записи"
		r9 := sheet.AddRow()
		r9.AddCell().Value = "Сумма контракта, руб."
	}
	if r.Lvl >= 90 {
		r10 := sheet.AddRow()
		r10.AddCell().Value = "Наименование контракта (for Яна with <3)"
	}
	return err
}*/

func (p *NotificationResult) isNotEmpty() bool {
	return p.Count != 0
}
func (r NotificationResult) NotificationsToTemplate() string {
	return "/templates/?save=" + url.QueryEscape(r.genURL(false, 0))
}
func (r *NotificationResult) genURL(xls bool, page int) string {
	//filterDateFrom=17.09.2018&filterDateTo=&query=&filterMaxPriceFrom=&filterMaxPriceTo=&filterPnum=&filterPage=2&perPage=25
	//http://195.201.86.48:96/notifications/?filterDateFrom=17.09.2018&filterDateTo=&query=&filterMaxPriceFrom=&filterMaxPriceTo=&filterPnum=&filterPage=3&perPage=25
	//filterDateFrom=17.09.2018&filterDateTo=31.12.2018&query=&filterPage=3&perPage=25
	s := "/auctions/?query=" + url.QueryEscape(r.Query)
	s += omitempty("Org", r.FilterOrg)
	s += omitempty("DateFrom", r.FilterDateFrom)
	s += omitempty("DateTo", r.FilterDateTo)
	s += omitempty("Customer", r.FilterCustomer)
	s += omitempty("OKPDCode", r.FilterCode)
	s += omitempty("Region", r.FilterRegion)
	s += omitempty("Anum", r.FilterPnum)
	//s += omitempty("filterRegnum", r.FilterRegnum)
	s += omitempty("MaxPriceFrom", r.FilterMaxPriceFrom)
	s += omitempty("MaxPriceTo", r.FilterMaxPriceTo)
	if r.FilterDocsFlag {
		s += omitempty("FilterDocsFlag", "on")
	}
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
func (p NotificationResult) Pages() template.HTML {
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
	pageNum := int(math.Ceil(float64(p.Total) / float64(p.PerPage)))
	var templateD string
	for i := 1; i <= pageNum; i++ {
		templateD += fmt.Sprintf("<a href=\"%s\">%d</a> &nbsp;", p.genURL(false, i), i)
		//template += fmt.Sprintf("<a href='/notifications/?query=%s&filterDateFrom=%s&filterDateTo=%s&filterPriceFrom=%2f&filterPriceTo=%2f&filterCustomer=%s&filterCode=%s&filterRegion=%s&perPage=%d&filterPage=%d'>%d</a> &nbsp;", url.QueryEscape(p.Query), p.FilterDateFrom, p.FilterDateTo, p.FilterPriceFrom, p.FilterPriceTo, url.QueryEscape(p.FilterCustomer), url.QueryEscape(p.FilterCode), url.QueryEscape(p.FilterRegion), p.PerPage, i, i)
	}
	*/
	return template.HTML(templateD)
}
func (p NotificationResult) ExcelLink() string {
	return p.genURL(true, 0)
	//return fmt.Sprintf("/contracts/?query=%s&filterDateFrom=%s&filterDateTo=%s&filterPriceFrom=%2f&filterPriceTo=%2f&filterSumFrom=%2f&filterSumTo=%2f&filterCode=%s&filterRegion=%s&filterCustomer=%s&filterSupplier=%s&xls=short", url.QueryEscape(p.Query), url.QueryEscape(p.FilterDateFrom), url.QueryEscape(p.FilterDateTo), p.FilterPriceFrom, p.FilterPriceTo, p.FilterSumFrom, p.FilterSumTo, urlEscapeStringArray(p.FilterCode, "filterCode"), urlEscapeStringArray(p.FilterRegion, "filterRegion"), urlEscapeStringArray(p.FilterCustomer, "filterCustomer"), url.QueryEscape(p.FilterSupplier))
}

// Notification - general description for result
type Notification struct {
	Row            int            `json:"row"`
	ObjectInfo     string         `json:"ObjectInfo"`
	PurchaseNumber string         `json:"PurchaseNumber"`
	Published      time.Time      `json:"published"`
	URL            string         `json:"url"`
	Org            sql.NullString `json:"org"`
	OrgINN         sql.NullString `json:"org_inn"`
	OrgKPP         sql.NullString `json:"org_kpp"`
	OrgRole        string         `json:"org_role"`
	OrgRegion      sql.NullInt64  `json:"org_region"`
	Code           string         `json:"code"`
	Query          string
	Type           string 		   `json:"type"`
	Canceled       bool
	Finished       bool
	Currency       sql.NullString
	FinanceSource  sql.NullString
	GMaxPrice      float64
	CMaxPrice      sql.NullFloat64
	Customer       sql.NullString
	CustomerINN    sql.NullString
	CustomerKPP    sql.NullString
	CustomerRegion sql.NullInt64
	CtPrice        sql.NullFloat64
	Regnum         sql.NullString
	CtPName        sql.NullString
}

func (p *Notification) PublishedDate() string {
	return fmt.Sprintf("%02d.%02d.%d", p.Published.Day(), p.Published.Month(), p.Published.Year())
}

func (p *Notification) PublishedXlsDate() string {
	return fmt.Sprintf("%02d.%02d.%d", p.Published.Day(), p.Published.Month(), p.Published.Year())
}

func (p *Notification) GetMaxPrice() string {
	var maxprice float64
	t, err := p.Customer.Value()
	if t == nil {
		maxprice = p.GMaxPrice
	} else {
		if err != nil {
			log.Println(err)
		}
		to, err := p.Org.Value()
		if err != nil {
			log.Println(err)
		}
		tp, err := p.CMaxPrice.Value()
		if to == t.(string) {
			maxprice = p.GMaxPrice
		} else {
			maxprice = tp.(float64)
		}
	}
	price := fmt.Sprintf("%.2f", maxprice)
	price = strings.Replace(price, ".", ",", -1)
	return price
}
func (p *Notification) PriceExcel() string {
	return p.GetMaxPrice()
}
func (p *Notification) HightlightQuery() string {
	words := p.ObjectInfo
	for _, w := range strings.Split(p.Query, " ") {
		re := regexp.MustCompile("(?i)" + w)
		words = re.ReplaceAllLiteralString(words, "<span style='background-color: #FFFF88;'>"+w+"</span>")
	}
	//return words
	return p.ObjectInfo
}
func (p Notification) FormatCtPrice() string {
	if p.CtPrice.Valid {
		t, _ := p.CtPrice.Value()
		price := fmt.Sprintf("%.2f", t)
		price = strings.Replace(price, ".", ",", -1)
		return price
	}
	return ""
}
func (p Notification) GetRegnum() string {
	if p.Regnum.Valid {
		t, _ := p.Regnum.Value()
		return t.(string)
	}
	return ""
}
func (p *NotificationResult) FormatSumTo() string {
	if p.FilterMaxPriceTo == 0 {
		return "0"
	} else {
		return fmt.Sprintf("%f", p.FilterMaxPriceTo)
	}
}
func (p *NotificationResult) FormatSumFrom() string {
	if p.FilterMaxPriceFrom == 0 {
		return "0"
	} else {
		return fmt.Sprintf("%f", p.FilterMaxPriceFrom)
	}
}

type NotificationFilters struct {
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

func (f *NotificationFilters) setFilters(params url.Values) {
	f.filterDocsFlag = params.Get("FilterDocsFlag") == "on"
	f.filterStrictFlag = params.Get("FilterStrictFlag") == "on"
	log.Println("Param FilterStrictFlag", params.Get("FilterStrictFlag"))

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
	if params["OKPDCode"] != nil {
		f.filterCode = params["OKPDCode"]
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

	f.filterCanceled, _ = strconv.ParseBool(params.Get("filterCanceled"))
	f.filterFinished, _ = strconv.ParseBool(params.Get("filterFinished"))

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

func (f *NotificationFilters) ApplyFilters(sc *sphinx.Client) {
	log.Println("Date filters: ", f.filterDateFrom, "-", f.filterDateTo)
	sc.SetFilterRange("published", filterDateTimestamp(f.filterDateFrom), filterToDateTimestamp(f.filterDateTo), false)
	if f.filterCodeN != nil && len(f.filterCodeN) != 0 {
		sc.SetFilter("okpd", f.filterCodeN, false)
		log.Println("Filter - OKPD:", f.filterCodeN, sc.GetLastError())
	}
	if f.filterPnumN != nil && len(f.filterPnumN) != 0 {
		sc.SetFilter("pnum", f.filterPnumN, false)
		log.Println("Filter - pnum:", f.filterPnumN, sc.GetLastError())
	}
	if f.filterRegnumN != nil && len(f.filterRegnumN) != 0 {
		sc.SetFilter("rnum", f.filterRegnumN, false)
		log.Println("Filter - regnum:", f.filterRegnumN, sc.GetLastError())
	}
	if f.filterRegionN != nil && len(f.filterRegionN) != 0 {
		sc.SetFilter("org_region", f.filterRegionN, false)
		log.Println("Filter - region:", f.filterRegionN, sc.GetLastError())
	}
	if f.filterMaxPriceTo != 0 {
		sc.SetFilterFloatRange("gMaxPrice", float32(f.filterMaxPriceFrom), float32(f.filterMaxPriceTo), false)
		log.Println("Filter - Price:", float32(f.filterMaxPriceFrom), float32(f.filterMaxPriceTo), sc.GetLastError())
	}
	if f.filterOrgN != nil && len(f.filterOrgN) != 0 {
		f.filterOrgN = hackMKDC(f.filterOrgN)
		sc.SetFilter("org", f.filterOrgN, false)
		log.Println("Filter - Orgs:", f.filterOrg, sc.GetLastError())
	}
	if f.filterCustomerN != nil && len(f.filterCustomerN) != 0 {
		f.filterCustomerN = hackMKDC(f.filterCustomerN)
		sc.SetFilter("customer", f.filterCustomerN, false)
		log.Println("Filter - Customers:", f.filterCustomer, sc.GetLastError())
	}

}

func (result *NotificationResult) setMeta(query string, f NotificationFilters) {
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
	result.FilterCode = f.filterCode
	result.FilterPnum = strings.Join(f.filterPnum, ",")
	result.FilterRegnum = strings.Join(f.filterRegnum, ",")
	result.FilterPage = int(f.filterPage)
	result.PerPage = f.perPage
	result.NavType = "notifications"
}

type plutoCMD struct {
	Query      string `json:"query"`
	FilterPage int    `json:"filterPage"`
	PerPage    int    `json:"perPage"`
}

func flyToPluto(query string) (ids []uint64) {
	conn := "http://tender-one.ru:7431?signal="
	signal := plutoCMD{}
	if query == "" {
		return
	} else {
		signal.Query = query
	}
	var sids []string
	t, _ := json.Marshal(signal)
	if resp, err := http.Get(conn + url.QueryEscape(string(t))); err != nil {
		log.Println(err)
		resp.Body.Close()
		return
	} else {
		defer resp.Body.Close()
		jres := json.NewDecoder(resp.Body)
		jres.Decode(&sids)
		for _, v := range sids {
			t, _ := strconv.ParseUint(v, 10, 64)
			ids = append(ids, t)
		}
	}
	return
}

func searchInDocs(q string, f NotificationFilters) (ids []uint64) {
	log.Println("Looking for docs...")
	// log.Println("docs query: ", q)
	for i := 1; i < 10; i++ {
		
		shard := fmt.Sprintf("star_ndocs_%d", i)

		sc := sphinx.NewClient().SetServer(sphinxSwitch.GetHost(), 0).SetConnectTimeout(1000).SetLimits(0, 10000000, 1000000, 0)
		sc.SetFilterRange("published", filterDateTimestamp(f.filterDateFrom), filterToDateTimestamp(f.filterDateTo), false)

		if err := sc.Error(); err != nil {
			log.Fatal(err)
		}

		res, err := sc.Query(q, shard, shard)
		
		if err != nil {
			log.Println("Query error:", err)
		}

		if res != nil {
			//log.Println("Docs-Shard: ", shard, "; Res: ", len(res.Matches))
			//log.Println(res.Matches)
			for _, match := range res.Matches {

				ids = append(ids, uint64(match.DocId))
			}
		}

	}
	return
}

func notificationsHandler(w http.ResponseWriter, r *http.Request, db, dbs *sql.DB) {
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
	var result NotificationResult
	result.Demo = lvl == 0
	result.Lvl = lvl
	var filters NotificationFilters
	filters.demo = result.Demo
	if result.Demo {
		result.FilterDateFrom = "01.01.2017"
		result.FilterDateTo = "31.12.2017"
	}
	result.NavType = "notifications"
	if len(params) != 0 {
		log.Println("Params Debug:", params)
		query = params.Get("query")
		filters.setFilters(params)

		result = notificationsSearch(db, query, filters)
		result.IsSearch = true
		result.Email = GetEmail(r)
		_, result.Lvl = accessLVL(GetEmail(r), 0, dbs)
		result.Fields = fields.FilterByLvl(fields.NotificationFields, result.Lvl)
		decreaseRequestsCount(GetEmail(r), dbs)
		LogQuery(GetEmail(r), getIP(r), result.Query, GenURLNotifications(result, params.Get("xls")), result.Total, params.Get("xls") == "short", dbs)

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
				w.Header().Set("Content-Disposition", "attachment; filename=Аукционы.xls")
				tmpl, err := template.New("").Delims("{[", "]}").ParseFiles("templates/notifications.xls.html")
				if err != nil {
					log.Fatal(err)
				}
				tmpl.ExecuteTemplate(w, "notifications.xls.html", result)
			}
			//tmpl := template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/loading-info.html","templates/notifications.html"))
			result.FilterDocsFlag = filters.filterDocsFlag
			result.FilterStrictFlag = filters.filterStrictFlag
			log.Println("TEST LVL:", result.Lvl)

			tmpl, err := template.ParseFiles("templates/notifications.html")
			tmpl.Execute(w, result)
			//err := tmpl.ExecuteTemplate(w, "notifications", result)
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
			//tmpl := template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/loading-info.html","templates/notifications.html"))
			//err := tmpl.ExecuteTemplate(w, "notifications", result)
			tmpl, err := template.ParseFiles("templates/notifications.html")
			tmpl.Execute(w, result)
			if err != nil {
				log.Println(err)
			}
			// log.Println(result.Data)
		}
	}

	//	fmt.Fprintf(w, "%s \n %s", query, parseQuery(query))
}

func sphinxNSearch(shard, q string, f NotificationFilters) ([]uint64, int) {

	var ids []uint64
	
	sc := sphinx.NewClient().SetServer(sphinxSwitch.GetHost(), 0).SetConnectTimeout(1000).SetLimits(0, 10000000, 1000000, 0)
	log.Println("Connected with Sphinx! Shard:", shard, shard[:4])
	//f.ApplyFilters(sc)

	log.Println("Date filters: ", f.filterDateFrom, "-", f.filterDateTo)
	sc.SetFilterRange("published", filterDateTimestamp(f.filterDateFrom), filterToDateTimestamp(f.filterDateTo), false)
	if f.filterCodeN != nil && len(f.filterCodeN) != 0 {
		sc.SetFilter("okpd", f.filterCodeN, false)
		log.Println("Filter - OKPD:", f.filterCodeN, sc.GetLastError())
	}
	if f.filterDocsData != nil && len(f.filterDocsData) != 0 && f.filterDocsFlag != false {
		sc.SetFilter("pnum", f.filterDocsData, false)
		q = ""
		log.Println("Filter - pnum (Docs's data!):", len(f.filterDocsData), sc.GetLastError())
	} else {
		if f.filterPnumN != nil && len(f.filterPnumN) != 0 {
			sc.SetFilter("pnum", f.filterPnumN, false)
			log.Println("Filter - pnum:", f.filterPnumN, sc.GetLastError())
		}
	}
	if f.filterRegnumN != nil && len(f.filterRegnumN) != 0 {
		sc.SetFilter("rnum", f.filterRegnumN, false)
		log.Println("Filter - regnum:", f.filterRegnumN, sc.GetLastError())
	}
	if f.filterRegionN != nil && len(f.filterRegionN) != 0 {
		sc.SetFilter("org_region", f.filterRegionN, false)
		log.Println("Filter - region:", f.filterRegionN, sc.GetLastError())
	}
	if f.filterMaxPriceTo != 0 {
		sc.SetFilterFloatRange("gMaxPrice", float32(f.filterMaxPriceFrom), float32(f.filterMaxPriceTo), false)
		log.Println("Filter - Price:", float32(f.filterMaxPriceFrom), float32(f.filterMaxPriceTo), sc.GetLastError())
	}

	if f.filterCustomerN != nil && len(f.filterCustomerN) != 0 {
		sc.SetFilter("customer", f.filterCustomerN, false)
		if strings.Split(shard, "_")[1] != "ncustomers" {
			sc.SetFilter("org", f.filterCustomerN, false)
			log.Println("Filter - Orgs:", f.filterOrg, sc.GetLastError())
		}
		log.Println("Filter - Customers:", f.filterCustomer, sc.GetLastError())
	} else {
		if f.filterOrgN != nil && len(f.filterOrgN) != 0 {
			sc.SetFilter("org", f.filterOrgN, false)
			log.Println("Filter - Orgs:", f.filterOrg, sc.GetLastError())
		}
	}

	// log.Println("The query:", q)
	if err := sc.Error(); err != nil {
		log.Fatal(err)
	}

	res, err := sc.Query(q, shard, shard)
	
	if err != nil {
		log.Println("Query error:", err)
	}

	//	log.Println(shd, " found: ", len(res.Matches), "/", res.Total, "|", res.Matches)
	if res == nil {
		//log.Println("Panic!", f.filterDocsData)
		log.Println("Panic! res is nil! Len of filter:", len(f.filterDocsData), "\n The query:", q, " and shard is ", shard, "\n and the filters are:", f)
		return ids, 0
	} else {
		ids = make([]uint64, len(res.Matches))
		log.Println(len(res.Matches), "/", res.Total)
		for i, match := range res.Matches {
			ids[i] = match.DocId

			//log.Println(match.DocId)
		}
		// log.Println(ids)
	}

	return ids, res.Total
}

type Bag struct {
	shard string
	id    int
}

type Bag_u64 struct {
	shard string
	id    uint64
}

func abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}

func notificationsSearch(db *sql.DB, query string, f NotificationFilters) NotificationResult {

	q := query

	// Prepare shards
	log.Println("Trying to generate shards", f.filterDateFrom, f.filterDateTo)
	//shard := shards(f.filterDateFrom, f.filterDateTo, "test_notifications")
	shard := shards.GenRangeByPattern("star_notifications_%s")[1:]
	// Above - temporary hack, because zero shard is off for now.
	if f.filterCustomerN != nil {
		//shard = AppendShards(shard, shards(f.filterDateFrom, f.filterDateTo, "cust_notifications"))
		shard = AppendShards(shard, shards.GenRangeByPattern("star_ncustomers_%s"))
		//shard = AppendShards(shard, shards(f.filterDateFrom, f.filterDateTo, "test_c1_notifications"))
		log.Println("Experimental searching...")
	}

	log.Println("Shards:", shard)

	//log.Println("Query before parse: ", q)
	if f.filterDocsFlag {
		// Old or new?
		if StrToDate(f.filterDateFrom).Before(StrToDate("31.12.2015")) {
			q = uquery.PlutoFix(q, f.filterStrictFlag)
		} else {
			// new here
			//log.Println("!! isStrict?", f.filterStrictFlag)
			q = uquery.ParseQuery(q, f.filterStrictFlag)
		}
	} else {
		//log.Println("!! isStrict?", f.filterStrictFlag)
		q = uquery.ParseQuery(q, f.filterStrictFlag)
	}

	var ans NotificationResult
	var pages []Bag_u64
	//var ids []int
	//var notifications []Notification
	//var total int
	if f.filterDocsFlag {
		if StrToDate(f.filterDateFrom).Before(StrToDate("31.12.2015")) {
			log.Println("Flying to Pluto...")
			f.filterDocsData = flyToPluto(q)
			log.Println("Pluto's results is:", len(f.filterDocsData))
		} else {
			f.filterDocsData = searchInDocs(q, f)
		}
		log.Println("We found", len(f.filterDocsData), "entries in docs; Looking in general data now...")
		foundInTheDocs := len(f.filterDocsData)
		/* Fake filters to get non docs flags too */
		f.filterDocsFlag = false
		for _, shd := range shard {
			var idsfrt []uint64
			qfke := uquery.ParseQuery(query, f.filterStrictFlag)
			//log.Println("qfke:", qfke)
			idsf, _ := sphinxNSearch(shd, qfke, f)
			for _, v := range idsf {
				idsfrt = append(idsfrt, uint64(v))
			}
			//log.Println("Faking docs (", q, "): +", len(idsfrt))
			f.filterDocsData = append(f.filterDocsData, idsfrt...)
		}
		f.filterDocsFlag = true
		log.Println("Additional", len(f.filterDocsData)-foundInTheDocs, " entries were found in general data.")
	}
	log.Println("Let's do real index search now")
	for _, shd := range shard {
		ids, _ := sphinxNSearch(shd, q, f)
		for _, id := range ids {
			var bag Bag_u64
			bag.id = id
			bag.shard = shd
			pages = append(pages, bag)
		}
	}
	//log.Println("Huinya test:", len(f.filterDocsData) == len(ids))

	//var sh string
	if f.filterPage == -1 {
		f.filterPage = 0
		f.perPage = MaxDocs
	}
	log.Println("Pages debug:", f.filterPage, ";", f.perPage, ";", f.perPage*(f.filterPage+1), ";", len(pages))
	ans.Total = len(pages)
	//id
	for i := f.perPage * f.filterPage; i < ans.Total; i++ {
		if i >= f.perPage * (f.filterPage + 1) {
			break
		}
		log.Println("Overflow: ", i, "/", len(pages))
		r, err := retrieveNotifications2(pages[i], db, query, f)
		if err != nil {
			log.Println(err)
		}
		ans.Data = concatNotifications(ans.Data, r)
	}
	/*
	for i := f.perPage * f.filterPage; i < f.perPage*(f.filterPage+1); i++ {
		if i >= len(pages) {
			log.Println("Overflow: ", i, "/", len(pages))
			r, err := retrieveNotifications(ids, db, sh, query, f)
			if err != nil {
				log.Println(err)
			}
			ans.Data = concatNotifications(ans.Data, r)
			break
		}
		//log.Println("Item:", i, i+1 >= perPage*(filterPage+1))
		if i+1 >= f.perPage*(f.filterPage+1) {
			r, err := retrieveNotifications(ids, db, sh, query, f)
			if err != nil {
				log.Println(err)
			}
			ans.Data = concatNotifications(ans.Data, r)
			sh = pages[i].shard
			ids = nil
		}
		if sh == "" {
			sh = pages[i].shard
		} else {
			if sh != pages[i].shard {
				r, err := retrieveNotifications(ids, db, sh, query, f)
				if err != nil {
					log.Println(err)
				}
				ans.Data = concatNotifications(ans.Data, r)
				sh = pages[i].shard
				ids = nil
			}
		}
		//log.Printf("SH: %s; ids_no_file: %d\n", pages[i].shard, len(ids_no_file))
		ids = append(ids, pages[i].id)
	}

	/**/
	ans.Count = len(ans.Data)
	ans.Query = query
	ans.setMeta(query, f)
	log.Println(ans.Count, "/", ans.Total)

	return ans
}

func concatNotifications(old1, old2 []Notification) []Notification {
	newslice := make([]Notification, len(old1)+len(old2))
	copy(newslice, old1)
	copy(newslice[len(old1):], old2)
	return newslice
}

func retrieveNotifications2(page Bag_u64, db *sql.DB, queryst string, f NotificationFilters) ([]Notification, error) {

	shard := strings.Replace(page.shard, "star_", "", -1)
	data := new([]Notification)
	table_num := strings.Split(shard, "_")[1]

	log.Println("SHARD", shard)
	query := fmt.Sprintf(`
		SELECT c.customer_name,
			c.maxprice,
			n.purchasenumber,
			n.objectInfo,
			timestamptz(n.published) as published,
			n.org_name, n.org_inn, n.org_kpp, n.org_role, n.org_region, n.url, n.type, n.maxprice, n.canceled, n.finished,
			c.customer_name, c.customer_inn, c.customer_kpp, c.customer_region, c.ctprice, c.regnum, c.ctobject
		FROM notifications_%s as n
			JOIN notifications_customers_%s as c on n.pnum=c.pnum
				and c.pnum in (%s) ORDER BY published;`, table_num, table_num, strconv.FormatUint(page.id, 10),
	)
	// log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	j := 0

	for rows.Next() {
		var row Notification
		var sr sql.NullString
		err := rows.Scan(&sr, &row.CMaxPrice, &row.PurchaseNumber, &row.ObjectInfo, &row.Published, &row.Org, &row.OrgINN, &row.OrgKPP, &row.OrgRole, &row.OrgRegion, &row.URL, &row.Type, &row.GMaxPrice, &row.Canceled, &row.Finished, &row.Customer, &row.CustomerINN, &row.CustomerKPP, &row.CustomerRegion, &row.CtPrice, &row.Regnum, &row.CtPName)
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
	//log.Println("!shard:", *data)
	// log.Println("Retrieving notifications...:", shard, ";", ids, ";Retrievied:", len(*data))
	return *data, err
}

func retrieveNotifications(ids []int, db *sql.DB, shard, queryst string, f NotificationFilters) ([]Notification, error) {

	shard = strings.Replace(shard, "star_", "", -1)
	//shard = strings.Replace(shard, "ncustomers", "notifications_customers", -1)
	//shard = strings.Replace(shard, "test_c1_", "", -1)

	data := new([]Notification)

	if len(ids) == 0 {
		return *data, nil
	}

	log.Println("SHARD", shard)
	query := strings.Replace(`
		SELECT c.customer_name,
			c.maxprice,
			n.purchasenumber,
			n.objectInfo,
			timestamptz(n.published) as published,
			n.org_name, n.org_inn, n.org_kpp, n.org_role, n.org_region, n.url, n.type, n.maxprice, n.canceled, n.finished,
			c.customer_name, c.customer_inn, c.customer_kpp, c.customer_region, c.ctprice, c.regnum, c.ctobject
		FROM notifications_%s as n
			JOIN notifications_customers_%s as c on n.pnum=c.pnum
				and c.pnum in (
	`, "%s", strings.Split(shard, "_")[1], -1)

	for l, i := range ids {
		if l != len(ids)-1 {
			query += strconv.Itoa(i) + ", "
		} else {
			query += strconv.Itoa(i) + ") " + " ORDER BY published;"
		}
	}
	// log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	j := 0

	for rows.Next() {
		var row Notification
		var sr sql.NullString
		err := rows.Scan(&sr, &row.CMaxPrice, &row.PurchaseNumber, &row.ObjectInfo, &row.Published, &row.Org, &row.OrgINN, &row.OrgKPP, &row.OrgRole, &row.OrgRegion, &row.URL, &row.Type, &row.GMaxPrice, &row.Canceled, &row.Finished, &row.Customer, &row.CustomerINN, &row.CustomerKPP, &row.CustomerRegion, &row.CtPrice, &row.Regnum, &row.CtPName)
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
	//log.Println("!shard:", *data)
	// log.Println("Retrieving notifications...:", shard, ";", ids, ";Retrievied:", len(*data))
	return *data, err
}
func (p Notification) FixPurchaseNumber() string {
	s := ""
	if len(p.PurchaseNumber) < 19 {
		for i := 0; i < 19-len(p.PurchaseNumber); i++ {
			s += "0"
		}

	}
	s += p.PurchaseNumber
	return s
}

func (p Notification) RegionGet() string {
	if p.CustomerRegion.Valid {
		tmp, _ := p.CustomerRegion.Value()
		return GetRegionFi(tmp.(int64))
	}
	return ""
}

func (p NotificationResult) FieldWithIdExist(id string) bool {
	return fields.FieldWithIdExist(p.Fields, id)
}
