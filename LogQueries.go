package main

import (
	"database/sql"
	//"fmt"
	"log"
	//"net/url"
	"time"
)

// LogQuery - saves users query information to site's database
func LogQuery(email, ip, query, url string, results int, xls bool, db *sql.DB) error {
	moment := time.Now()
	if _, err := db.Exec("insert into queries (email, query, url, results, xls, ip, time) values($1,$2,$3,$4,$5,$6,to_timestamp(translate($7, 'T', ' '), 'YYYY-MM-DD HH24:MI:SS'))", email, query, url, results, xls, ip, moment); err != nil {
		log.Println("LogQuery ERROR: ", err)
		return err
	}
	return nil
}

// GenURLContracts generates url for LogQuery func using setted filters
func GenURLContracts(p ContractResult, xls string) string {
	xlsFlag := xls != ""
	return p.genURL(xlsFlag, p.FilterPage)
}

// GenURLNotifications generates url for LogQuery func using setted filters
func GenURLNotifications(p NotificationResult, xls string) string {
	xlsFlag := xls != ""
	return p.genURL(xlsFlag, p.FilterPage)
	//return fmt.Sprintf("http://new.monitoring-crm.ru/notifications/?query=%s&filterDateFrom=%s&filterDateTo=%s&filterPriceFrom=%2f&filterPriceTo=%2f&filterCustomer=%s&filterCode=%s&filterRegion=%s&perPage=%d&filterPage=%dxls=%s", url.QueryEscape(p.Query), p.FilterDateFrom, p.FilterDateTo, p.FilterPriceFrom, p.FilterPriceTo, url.QueryEscape(p.FilterCustomer), url.QueryEscape(p.FilterCode), url.QueryEscape(p.FilterRegion), p.PerPage, p.FilterPage, xls)
}

func GenURLNotifications223(p Notification223Result, xls string) string {
	xlsFlag := xls != ""
	return p.genURL(xlsFlag, p.FilterPage)
	//return fmt.Sprintf("http://new.monitoring-crm.ru/notifications/?query=%s&filterDateFrom=%s&filterDateTo=%s&filterPriceFrom=%2f&filterPriceTo=%2f&filterCustomer=%s&filterCode=%s&filterRegion=%s&perPage=%d&filterPage=%dxls=%s", url.QueryEscape(p.Query), p.FilterDateFrom, p.FilterDateTo, p.FilterPriceFrom, p.FilterPriceTo, url.QueryEscape(p.FilterCustomer), url.QueryEscape(p.FilterCode), url.QueryEscape(p.FilterRegion), p.PerPage, p.FilterPage, xls)
}
