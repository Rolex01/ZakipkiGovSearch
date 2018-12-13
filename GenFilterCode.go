package main

import (
	"fmt"
	"log"
	"strconv"

	"bitbucket.org/company-one/tender-one/api"
	"html/template"
)

func toInt(c []string) (n []int) {
	for _, v := range c {
		l, _ := strconv.ParseInt(v, 10, 64)
		n = append(n, int(l))
	}
	return
}

func genFilterCode(codes []string) (html string) {
	if len(codes) != 0 {

		if okpds, err := retrieveOKPD(toInt(codes), dbM); err != nil {
			log.Println(err)
			return
		} else {
			for _, code := range okpds {
				html += fmt.Sprintf("<option value=\"%d\" title=\"%s\" selected=\"selected\">%s</option>", code.Row, code.Code, code.Code)
			}
		}
	}

	return
}

func (p ContractResult) GenFilterCodeContracts() template.HTML {
	html := ""
	if p.FilterCode == nil {
		return ""
	}
	for _, code := range api.StringIdsToOKPD1(p.FilterCode, dbM) {
		html += fmt.Sprintf("<option value=\"%d\" title=\"%s\" selected=\"selected\">%s</option>", code.Row, code.Code, code.Code)
		log.Println("GTest:", code)
	}
	return template.HTML(html)
}
func (p ContractResult) GenFilterCode2Contracts() template.HTML {
	if p.FilterCode2 == nil {
		return ""
	}
	html := ""
	for _, code := range api.StringIdsToOKPD2(p.FilterCode2, dbM) {
		html += fmt.Sprintf("<option value=\"%d\" title=\"%s\" selected=\"selected\">%s</option>", code.Row, code.Code, code.Code)
		log.Println("GTest:", code)
	}
	return template.HTML(html)
}
func (p NotificationResult) GenFilterCodeNotifications() template.HTML {
	res := genFilterCode(p.FilterCode)
	log.Println("Generating filter OKPD Code", p.FilterCode, res)
	return template.HTML(res)
}
