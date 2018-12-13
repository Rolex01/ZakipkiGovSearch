package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func listUsersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	users := fillUsers(listUsers(db), true)
	w.Header().Set("Content-type", "application/json")
	if s, err := json.Marshal(users); err != nil {
		log.Println(err)
		// log.Println(users)
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, "%s", s)
	}
}

// User describes user properties
type User struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Permission     int    `json:"permission"`
	Created        string `json:"created"`
	LifeTime       string `json:"lifetime"`
	SessionsLimit  int    `json:"sessionslimit"`
	RequestsLimit  int    `json:"requestslimit"`
	tSessionsLimit sql.NullInt64
	tRequestsLimit sql.NullInt64
	tCreated       time.Time
	tLifeTime      time.Time
	tIP            sql.NullString
	IP             string `json:"ip"`
	Verified       bool   `json:"verified"`
	Templates      string `json:"templates"`
	tTemplates     sql.NullString
}

func (u *User) Fill() {
	if v, _ := u.tRequestsLimit.Value(); u.tRequestsLimit.Valid {
		u.RequestsLimit = int(v.(int64))
	}
	if v, _ := u.tSessionsLimit.Value(); u.tSessionsLimit.Valid {
		u.SessionsLimit = int(v.(int64))
	}
	if v, _ := u.tIP.Value(); u.tIP.Valid {
		u.IP = v.(string)
	}
	if v, _ := u.tTemplates.Value(); u.tTemplates.Valid {
		u.Templates = v.(string)
	}
}

func fillUsers(u []User, jsFormat bool) []User {
	for i, _ := range u {
		if v, _ := u[i].tRequestsLimit.Value(); u[i].tRequestsLimit.Valid {
			u[i].RequestsLimit = int(v.(int64))
		}
		if v, _ := u[i].tSessionsLimit.Value(); u[i].tSessionsLimit.Valid {
			u[i].SessionsLimit = int(v.(int64))
		}
		if v, _ := u[i].tIP.Value(); u[i].tIP.Valid {
			u[i].IP = v.(string)
		}
		if v, _ := u[i].tTemplates.Value(); u[i].tTemplates.Valid {
			u[i].Templates = v.(string)
		}
		if jsFormat {
			u[i].LifeTime = jsTime(u[i].tLifeTime)
			u[i].Created = jsTime(u[i].tCreated)
			log.Println("jsFormatting!", jsTime(u[i].tCreated), timeToStr(u[i].tCreated))
		}
	}
	return u
}

func addZero(i int) string {
	if i >= 10 {
		return fmt.Sprintf("%d", i)
	}
	return fmt.Sprintf("0%d", i)

}
func timeToStr(t time.Time) string {
	return fmt.Sprintf("%s.%s.%d", addZero(t.Day()), addZero(int(t.Month())), t.Year())
}

func jsTime(t time.Time) string {
	return fmt.Sprintf("%s.%s.%d", addZero(int(t.Month())), addZero(t.Day()), t.Year())
}

func (p *User) FormatTime() {
	p.Created = timeToStr(p.tCreated)
	p.LifeTime = timeToStr(p.tLifeTime)
	log.Println("tCreated:", p.tCreated, ";Created:", p.Created)
}

func (p User) SessionsLimitRule() string {
	if p.SessionsLimit == 0 {
		return "Без ограничений"
	}
	return fmt.Sprintf("%d", p.SessionsLimit)
}

func (p User) PermissionText() string {
	switch p.Permission {
	case -2:
		return "Клиент (Просрочена оплата)"
		break
	case -1:
		return "Демо (Просрочен)"
		break
	case 0:
		return "Демо"
		break
	case 1:
		return "Клиент"
		break
	case 2:
		return "Сотрудник"
	case 90:
		return "Администратор"
		break
	case 100:
		return "Разработчик"
		break
	default:
		return "Не указан"
		break
	}
	return "Не указан"
}

// ControlUsersParams - context params for control-users page
type ControlUsersParams struct {
	Users        []User
	Action       string
	Demo         bool
	NavType      string
	NotPermitted bool
	User         User
	Lvl          int
}

func controlUsersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var context ControlUsersParams
	tmpl := template.Must(template.New("").Delims("{[", "]}").ParseFiles("templates/header.html", "templates/footer.html", "templates/control-users-list.html"))

	params := getParams(r)
	/*
	   switch on params.Get("action").
	   list - list all users
	   add  - add a new user
	   edit - edit selected user. Required: email
	   remove - remove selected user. Required: email. Level check here. User with lower priveleges can't remove user with higher privelege
	*/
	switch params.Get("action") {
	case "list":
		context.Users = listUsers(db)
		context.Action = "list"
		context.Demo = false
		if err := tmpl.ExecuteTemplate(w, "control-users", context); err != nil {
			log.Println(err)
		}
		break
	case "edit":
		if params.Get("email") != "" {
			if params.Get("commit") != "" {

				user := User{}
				user.Name = params.Get("name")
				user.Email = params.Get("email")
				user.Password = params.Get("password")
				t, _ := strconv.ParseInt(params.Get("permission"), 10, 32)
				user.Permission = int(t)
				user.Created = params.Get("created")
				user.LifeTime = params.Get("lifetime")
				t, _ = strconv.ParseInt(params.Get("sessionslimit"), 10, 32)
				user.SessionsLimit = int(t)
				log.Println("SLIMIT:", t)
				t, _ = strconv.ParseInt(params.Get("requestslimit"), 10, 32)
				user.RequestsLimit = int(t)
				log.Println("RLIMIT:", t)
				context.User = user
				log.Println("COMMIT", user)
				editUser(user, db)
				http.Redirect(w, r, "/control/?action=list", http.StatusFound)
			} else {
				user := getUser(params.Get("email"), db)
				context.User = user
			}
			if test, _ := accessLVL(GetEmail(r), context.User.Permission, db); !test {
				context.NotPermitted = true
			}

			context.Action = "edit"
			context.Demo = false

		} else {
			context.Users = listUsers(db)
			context.Action = "list"
		}
		if err := tmpl.ExecuteTemplate(w, "control-users", context); err != nil {
			log.Println(err)
		}
		break
	case "add":
		context.Action = "add"
		context.Demo = false
		if params.Get("email") != "" {
			var u User
			u.Email = params.Get("email")
			t, _ := strconv.ParseInt(params.Get("permission"), 10, 64)
			u.Permission = int(t)
			u.Password = params.Get("password")
			u.Created = params.Get("created")
			u.Name = params.Get("name")
			log.Println("NAME:", u.Name)
			u.LifeTime = params.Get("lifetime")
			t, _ = strconv.ParseInt(params.Get("sessionslimit"), 10, 64)
			u.SessionsLimit = int(t)
			t, _ = strconv.ParseInt(params.Get("requestslimit"), 10, 64)
			u.RequestsLimit = int(t)
			context.User = u
			addUser(u, db)
			http.Redirect(w, r, "/control/?action=list", http.StatusFound)
		}
		if err := tmpl.ExecuteTemplate(w, "control-users", context); err != nil {
			log.Println(err)
		}
		break
	case "remove":
		context.Action = "remove"
		context.Demo = false
		if params.Get("email") != "" {
			// ToDo Permission check.
			rmUser(params.Get("email"), db)
			http.Redirect(w, r, "/control/?action=list", http.StatusFound)
		}
		if err := tmpl.ExecuteTemplate(w, "control-users", context); err != nil {
			log.Println(err)
		}
		break
	}
}

func (p ControlUsersParams) GenPermissionsList() (html string) {
	markup := map[int]string{
		-2:  "Клиент (Просрочена оплата)",
		-1:  "Демо (Просрочен)",
		0:   "Демо",
		1:   "Клиент",
		2:   "Сотрудник",
		90:  "Администратор",
		100: "Разработчик"}
	l := p.User.Permission
	for i, v := range markup {
		if l == i {
			html += fmt.Sprintf("<option value=\"%d\" selected>%s</option>", i, v)
		} else {
			html += fmt.Sprintf("<option value=\"%d\">%s</option>", i, v)
		}
		//log.Println("l", l, "i", i, "html", html)
	}
	/*
		<option value="-1">Демо (Просрочен)</option>
			<option value="-2">Клиент (Просрочена оплата)</option>
			<option value="0">Демо</option>
			<option value="1">Клиент</option>
			<option value="2">Сотрудник</option>
			<option value="11">Администратор</option>
	*/
	return
}
