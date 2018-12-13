package main

import (
	"fmt"
	"log"
	"strings"
	"html/template"
)

func genFilterCustomer(customers []string) (html string) {
	log.Println("genFilterCustomer:", customers, ";", len(customers))
	if len(customers) != 0 {
		//var orgs []Organization
		q := strings.Join(customers, "|")
		log.Println("q", q)
		if orgs, err := retrieveOrganizations(sphinxOrgSearch(q), dbM); err != nil {
			log.Println(err)
			return
		} else {

			for _, org := range orgs {
				html += fmt.Sprintf("<option value=\"%s\" selected=\"selected\" title='%s'>%s</option>", org.Regnum, org.Shortname, org.Shortname)
				// log.Println(fmt.Sprintf("<option value=\"%s\" selected=\"selected\" title='%s'>%s</option>", org.Regnum, org.Shortname, org.Shortname))
			}

		}
	} else {
		html = ""
	}
	return
}

func genFilterSupplier(customers []string) (html string) {
	log.Println("genFilterCustomer:", customers, ";", len(customers))
	if len(customers) != 0 {
		//var orgs []Organization
		q := strings.Join(customers, "|")
		log.Println("q", q)
		if orgs, err := retrieveSuppliers(sphinxSuppliersSearch(q), dbM); err != nil {
			log.Println(err)
			return
		} else {

			var text string

			for _, org := range orgs {

				if len(org.Name) > 15 {
					text = org.Name[0:15] + "..."
				} else {
					text = org.Name
				}

				html += fmt.Sprintf(
					"<option value=\"%s\" selected=\"selected\" title='%s'>%s</option>",
					org.INN,
					org.Name,
					text,
                )
			}

		}
	} else {
		html = ""
	}
	return
}

func (p ContractResult) GenFilterCustomerContracts() template.HTML {

	res := genFilterCustomer(p.FilterCustomer)
	log.Println("Generating filter customer", p.FilterCustomer, res)
	return template.HTML(res)
}

func (p ContractResult) GenFilterSupplierContracts() template.HTML {

	res := genFilterSupplier(p.FilterSupplier)
	log.Println("Generating filter customer", p.FilterSupplier, res)
	return template.HTML(res)
}

func (p NotificationResult) GenFilterCustomerNotifications() template.HTML {

	res := genFilterCustomer(p.FilterCustomer)
	log.Println("Generating filter customer", p.FilterCustomer, res)
	return template.HTML(res)
}
func (p Notification223Result) GenFilterCustomerNotifications() template.HTML {

	res := genFilterCustomer(p.FilterCustomer)
	log.Println("Generating filter customer", p.FilterCustomer, res)
	return template.HTML(res)
}


func (p NotificationResult) GenFilterOrgNotifications() template.HTML {

	res := genFilterCustomer(p.FilterOrg)
	log.Println("Generating filter customer", p.FilterOrg, res)
	return template.HTML(res)
}
func (p Notification223Result) GenFilterOrgNotifications() template.HTML {

	res := genFilterCustomer(p.FilterOrg)
	log.Println("Generating filter customer", p.FilterOrg, res)
	return template.HTML(res)
}