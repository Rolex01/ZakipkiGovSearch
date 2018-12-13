package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	_ "github.com/lib/pq"
)

type TemplatePage struct {
	Data       []tmpl
	SaveParams string
	NavType    string
	Demo       bool
	Lvl        int
}

func (p TemplatePage) IsEmpty() bool {
	return len(p.Data) == 0
}

type tmpl struct {
	Name string `json:"Name"`
	Link string `json:"Link"`
}

func loadTemplates(email string, db *sql.DB) (templates []tmpl) {
	var tmpS sql.NullString
	if err := db.QueryRow("select templates from users where email=$1;", email).Scan(&tmpS); err != nil {
		log.Println("While loading templates from user: ", email, "; ", err)
	}
	t, err := tmpS.Value()
	if err != nil {
		log.Println(err)
	}
	var tmp string
	if t == nil {
		tmp = "[]"
	} else {
		tmp = t.(string)
		if tmp == "" {
			tmp = "[]"
		}
	}
	if err := json.Unmarshal([]byte(tmp), &templates); err != nil {
		log.Println("While parsing templates from user: ", email, "; ", err)
	}
	return
}
func saveTemplates(email string, tmp tmpl, db *sql.DB) {
	tmps := loadTemplates(email, db)
	tmps = append(tmps, tmp)
	var pack string
	var packs []byte
	var err error
	if packs, err = json.Marshal(tmps); err != nil {
		log.Println("While serializing templates for user: ", email, "; ", err, " data:", tmps)
	}
	pack = string(packs)
	if _, err := db.Exec("update users set templates=$1 where email=$2;", pack, email); err != nil {
		log.Println("While saving templates for user: ", email, "; data:", pack)
	}
}
func deleteTemplate(email string, tmp tmpl, db *sql.DB) {
	tmpls := loadTemplates(email, db)
	var i int
	var v tmpl
	for i, v = range tmpls {
		if v.Name == tmp.Name && v.Link == tmp.Link {
			break
		}
	}
	tmpls = append(tmpls[:i], tmpls[i+1:]...)
	var pack string
	var packs []byte
	var err error
	if packs, err = json.Marshal(tmpls); err != nil {
		log.Println("While serializing templates for user: ", email, "; ", err, " data:", tmpls)
	}
	pack = string(packs)
	if _, err := db.Exec("update users set templates=$1 where email=$2;", pack, email); err != nil {
		log.Println("While saving templates for user: ", email, "; data:", pack)
	}
}

func (p tmpl) DelLink() string {
	s, _ := json.Marshal(p)
	return "/templates/?delete=" + url.QueryEscape(string(s))
}
func templatesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var result TemplatePage
	params := url.Values{}
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	if r.Method == "GET" {
		params = r.URL.Query()
	}
	//log.Println("Params Debug:", params, r)

	if len(params) != 0 {
		//		log.Println("Params Debug:", params)

		if params["save"] != nil && params["name"] == nil {
			s := params.Get("save")
			result.SaveParams = s
			t, err := template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/templates.html")
			// log.Println("TESTING", err, t)
			tmpl := template.Must(t, err)
			err = tmpl.ExecuteTemplate(w, "templates", result)
			if err != nil {
				log.Println(err)
			}
			return
		}
		if params["save"] != nil && params["name"] != nil {
			var tmp tmpl
			tmp.Link = params.Get("save")
			tmp.Name = params.Get("name")
			saveTemplates(GetEmail(r), tmp, db)
			//list
			result.Data = loadTemplates(GetEmail(r), db)
			t, err := template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/templates.html")
			// log.Println("TESTING", err, t)
			tmpl := template.Must(t, err)
			err = tmpl.ExecuteTemplate(w, "templates", result)
			if err != nil {
				log.Println(err)
			}
			return
		}
		if params["delete"] != nil {
			var tmp tmpl
			//log.Println(params["delete"])
			if err := json.Unmarshal([]byte(strings.Join(params["delete"], "")), &tmp); err != nil {
				log.Println("On delete template: ", err)
			}
			deleteTemplate(GetEmail(r), tmp, db)
			// list here
			result.Data = loadTemplates(GetEmail(r), db)
			t, err := template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/templates.html")
			// log.Println("TESTING", err, t)
			tmpl := template.Must(t, err)
			err = tmpl.ExecuteTemplate(w, "templates", result)
			if err != nil {
				log.Println(err)
			}

			return
		}
	} else {
		// list here
		result.Data = loadTemplates(GetEmail(r), db)
		t, err := template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/templates.html")
		// log.Println("TESTING", err, t)
		tmpl := template.Must(t, err)
		err = tmpl.ExecuteTemplate(w, "templates", result)
		if err != nil {
			log.Println(err)
		}
	}

}
