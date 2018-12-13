package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type CabinetPage struct {
	Templates  []tmpl
	Requests   int
	Lifetime   string
	NavType    string
	Demo       bool
	Lvl        int
	Username   string
}

func (p CabinetPage) TemplatesIsEmpty() bool {
	return len(p.Templates) == 0
}

func (p CabinetPage) ContractTemplatesIsEmpty() bool {

	counter := 0

	for _, t := range p.Templates {

		if p.IsContractTemplate(t.Link) {
			counter++
		}
	}

	return counter == 0
}

func (p CabinetPage) NotificationTemplatesIsEmpty() bool {

	counter := 0

	for _, t := range p.Templates {

		if p.IsNotificationTemplate(t.Link) {
			counter++
		}
	}

	return counter == 0
}

func (p CabinetPage) IsContractTemplate(templateUrl string) bool {
	return strings.Contains(templateUrl, "/contracts/")
}

func (p CabinetPage) IsNotificationTemplate(templateUrl string) bool {
	return strings.Contains(templateUrl, "/notifications/")
}

func cabinetHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var result CabinetPage

	email := GetEmail(r)

	_, lvl := accessLVL(email, 0, db)
	result.Demo = lvl == 0
	result.Lvl = lvl
	result.NavType = "cabinet"
	result.Username = getUser(email, db).Name
	result.Templates = loadTemplates(email, db)

	var lifetime time.Time
	var requestslimit int

	err := db.QueryRow(`
		SELECT lifetime, requestslimit FROM users WHERE email=$1
	`, email).Scan(&lifetime, &requestslimit)
	
	if err != nil {
		log.Println(err)
	}

	result.Requests = requestslimit
	result.Lifetime = lifetime.Format("01.02.2006")

	t, err := template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/cabinet.html")

	tmpl := template.Must(t, err)
	err = tmpl.ExecuteTemplate(w, "cabinet", result)
	if err != nil {
		log.Println(err)
	}
}
