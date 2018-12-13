package main

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"gopkg.in/gomail.v2"
)

const permissionDemo int = 0
const permissionUser int = 1
const permissionEmployee int = 2

/**
 * simple mail func w/sendmail
 */
func mail(to, from, subject, message string) error {
	/*
	var err error

	cmd := exec.Command("/usr/sbin/sendmail", "-t", "-i")
	pipe, err := cmd.StdinPipe()

	if err != nil {
		return err
	}

	//data := []byte("Мониторинг цен госзакупок")
	data := []byte("")
	str := "=?utf-8?b?" + base64.StdEncoding.EncodeToString(data) + "?= <help@tender-one.ru>\r"
	from = str

	if err = cmd.Start(); err != nil {
		return err
	}

	_, err = fmt.Fprintf(pipe, "Content-Type: text/html; charset=\"utf-8\"\n")
	_, err = fmt.Fprintf(pipe, "To: %s\n", to)
	_, err = fmt.Fprintf(pipe, "From: %s\n", from)
	_, err = fmt.Fprintf(pipe, "Reply-to: %s\n", from)
	//_, err = fmt.Fprintf(pipe, "_, err = fmt.Fprintf(pipe, "From: %s\n", from) %s\n", from)
	_, err = fmt.Fprintf(pipe, "Subject: %s\n", subject)
	_, err = fmt.Fprintf(pipe, "\n%s\n", message)

	err = pipe.Close()
	if err != nil {
		return err
	}
*/
	return nil //cmd.Wait()
}
func sendMail(to, from, subject, message string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "Tender One")
	m.SetAddressHeader("To", to, "")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := gomail.NewPlainDialer("smtp.mail.ru", 465, "info@tender-one.ru", "xxx")

	err := d.DialAndSend(m);
	if  err != nil {
		log.Fatal(err)
	}

	//str := "=?utf-8?b?" + base64.StdEncoding.EncodeToString(data) + "?= <help@tender-one.ru>\r"

	return err//cmd.Wait()
}

func UsersTable(db *sql.DB) {
	// -1 - Unpaid; 0 - demo; 1 - user; 99 - admin; 100 - master.
	const q string = `CREATE table IF NOT EXISTS users (
	row					SERIAL		PRIMARY KEY,
	Email				text		UNIQUE,
	Name				text		,
	Created		       	timestamp	,
	Permission			numeric		DEFAULT 0,
	Password			text		,
	IP					text		,
	Verified			boolean		DEFAULT FALSE,
	Templates			text
 );`
	_, err := db.Exec(q)
	if err != nil {
		log.Println(err)
	}
}

func SessionsTable(db *sql.DB) {

	const q string = `CREATE table IF NOT EXISTS sessions (
	row					SERIAL		PRIMARY KEY,
	Email				text		,
	Time		       	timestamp	,
	IP					text
 );`
	_, err := db.Exec(q)
	if err != nil {
		log.Println(err)
	}
}

func RemindTable(db *sql.DB) {
	const q string = `CREATE table IF NOT EXISTS remind (
	row					SERIAL		PRIMARY KEY,
	Email				text		,
	Time		       	timestamp	,
	IP					text,
	ObsoluteP			text,
	Password			text,
	Key					text
 );	`
	_, err := db.Exec(q)
	if err != nil {
		log.Println(err)
	}
}

func QueriesTable(db *sql.DB) {

	const q string = `CREATE table IF NOT EXISTS queries (
	row					SERIAL		PRIMARY KEY,
	Email				text		,
	Time		       	timestamp	,
	IP					text		,
	Query				text		,
	Filters				text		,
	Url					text		,
	Results				numeric		,
	Xls					boolean		DEFAULT FALSE
 );`
	_, err := db.Exec(q)
	if err != nil {
		log.Println(err)
	}
}

func pwdCheck(db *sql.DB, email, password string) bool {
	//h := md5.New()
	//h.Write([]byte(password))
	//password = hex.EncodeToString(h.Sum(nil))
	var controlSum string
	db.QueryRow("SELECT password from users WHERE email=$1", email).Scan(&controlSum)
	log.Println("controlSum: ", controlSum, "password: ", password)
	if controlSum == password {
		return true
	}
	return false
}
func isAuthorized(db *sql.DB, session *sessions.Session) bool {
	//var email string
	//var pwd string
	// 1. is session empty?
	if v, ok := session.Values["email"]; ok {
		email, _ := v.(string)
		if v, ok := session.Values["pwd"]; ok {
			pwd, _ := v.(string)
			// 2. Checking control sum
			return pwdCheck(db, email, pwd)
		} else {
			return false
		}
	} else {
		return false
	}
}
func accessLVL(email string, lvl int, db *sql.DB) (bool, int) {
	var resLVL int
	//var createdStr string
	var created time.Time
	err := db.QueryRow("SELECT permission,created from users WHERE email=$1", email).Scan(&resLVL, &created)
	log.Println("AccessLVL for ", email, " is ", resLVL, ";Created:", created, ";Since:", time.Since(created).Hours())
	if time.Since(created).Hours() > 24*5 && resLVL == 0 {
		if _, err := db.Exec("update users set permission=-1 where email=$1", email); err != nil {
			log.Println(err)
		}
		return false, -1
	}
	if err != nil {
		log.Println(err)
	}
	return resLVL >= lvl, resLVL
}

func checkLVL(lvl int, db *sql.DB, session *sessions.Session) bool {
	if v, ok := session.Values["email"]; ok {
		s, _ := v.(string)
		so, _ := accessLVL(s, lvl, db)
		return so
	} else {
		return false
	}

}

func sessionsCount(email string, db *sql.DB) (count int) {
	if err := db.QueryRow("SELECT COUNT(*) from sessions where email=$1 and now() - interval '2 day';", email).Scan(&count); err != nil {
		log.Println(err)
	}
	return
}

func rulersChk(r *http.Request, db *sql.DB) bool {
	//  lifetime, sessionslimit, requestslimit
	var lifetime time.Time
	var sessionslimit, requestslimit int
	if err := db.QueryRow("SELECT sessionslimit, lifetime, requestslimit from users where email=$1", GetEmail(r)).Scan(&sessionslimit, &lifetime, &requestslimit); err != nil {
		log.Println(err)
	}

	//return requestslimit>=0 && sessionsCount(GetEmail(r),db) <= sessionslimit && time.Since(lifetime).Hours() >= time.Since(time.Now()).Hours()
	log.Println("Requests rule:", requestslimit >= 0)
	log.Println("Lifetime rule:", time.Since(lifetime).Hours() <= time.Since(time.Now()).Hours())
	return requestslimit >= 0 && time.Since(lifetime).Hours() <= time.Since(time.Now()).Hours()
}

func decreaseRequestsCount(email string, db *sql.DB) (count int) {
	if err := db.QueryRow("SELECT requestslimit from users where email=$1;", email).Scan(&count); err != nil {
		log.Println(err)
	}
	log.Println("Debug decreaseRequestsCount: selected rlimit:", count)
	if count <= 0 {
		log.Println("Debug decreaseRequestsCount: <= 0 case;")
		return
	} else {
		if count-1 == 0 {
			count = count - 2
			log.Println("Debug decreaseRequestsCount: count-1 == 0 case; Result count:", count)
		} else {
			count = count - 1
			log.Println("Debug decreaseRequestsCount: count-1 case; Result count:", count)
		}
		if _, err := db.Exec("UPDATE users set requestslimit=$2 where email=$1", email, count); err != nil {
			log.Println(err)
		}
	}
	return
}

func login(email, password, ip string, db *sql.DB, session *sessions.Session) bool {
	fpassword := password
	h := md5.New()
	h.Write([]byte(password))

	password = hex.EncodeToString(h.Sum(nil))
	if !pwdCheck(db, email, password) {
		log.Println("Login Error: ", email, password, ip)
		return false // incorrect user data
	}
	// Registering session
	timestamp := time.Now().UTC()
	_, err := db.Exec("INSERT INTO sessions (email,ip,time) VALUES ($1,$2,to_timestamp(translate($3, 'T', ' '), 'YYYY-MM-DD HH24:MI:SS'))", email, ip, timestamp)
	if err != nil {
		log.Println("Login error: ", err)
	}
	_, session.Values["lvl"] = accessLVL(email, 0, db)
	if session.Values["lvl"].(int) < 0 {
		return false
	}
	session.Values["email"] = email
	session.Values["pwd"] = password
	session.Values["f"] = fpassword
	log.Println("Saving to session:", password, fpassword)
	session.Values["timestamp"] = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d\n",
		timestamp.Year(), timestamp.Month(), timestamp.Day(),
		timestamp.Hour(), timestamp.Minute(), timestamp.Second())
	return true
}
func logout(session *sessions.CookieStore) {
	session.MaxAge(-1)
}
func signUp(email, password, ip, name string, permission int, db *sql.DB) bool {
	var count int
	db.QueryRow("SELECT count(*) FROM users where email = $1", email).Scan(&count)
	if count != 0 {
		log.Println("User with email ", email, "already exists")
		return false // Username already exists
	}
	tmpl, err := template.ParseFiles("templates/reg.mail")
	if err != nil {
		log.Println(err)
	}
	var body bytes.Buffer
	var params RemindParams
	params.Email = email
	params.Password = password
	t := time.Now()
	params.Today = fmt.Sprintf("%02d.%02d.%d\n", t.Day(), t.Month(), t.Year())
	err = tmpl.Execute(&body, params)
	if err != nil {
		log.Println(err)
	}
	//err = mail(email, "БД Мониторинг цен госзакупок <from@monitoring-crm.ru>", "Демо доступ к БД \"Мониторинг цен госзакупок\"", body.String())
	err = mail(email, "", "", body.String())
	if err != nil {
		log.Println(err)
	}
	timestamp := time.Now().UTC()
	h := md5.New()
	h.Write([]byte(password))
	password = hex.EncodeToString(h.Sum(nil))
	db.Exec("INSERT INTO users (email,name,permission,password,ip,created,lifetime,sessionslimit,requestslimit) VALUES ($1,$2,$3,$4,$5,to_timestamp(translate($6, 'T', ' '), 'YYYY-MM-DD HH24:MI:SS'),'01.01.2021',0,0)", strings.TrimSpace(email), name, permission, password, ip, timestamp)

	return true
}

func getIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func getForgotKey(email, ip string, db *sql.DB) bool {
	var controlSum, obsoluteP string
	var params RemindParams
	db.QueryRow("SELECT password from users WHERE email=$1", email).Scan(&obsoluteP)
	if obsoluteP == "" {
		return false
	}
	controlSum = email + obsoluteP
	h := md5.New()
	h.Write([]byte(controlSum))
	controlSum = hex.EncodeToString(h.Sum(nil))

	tmpl, err := template.ParseFiles("templates/forgot.mail")
	if err != nil {
		log.Println(err)
	}
	var body bytes.Buffer
	params.Email = email
	params.Password = randStr(8, "alphanum")
	params.Key = controlSum
	timestamp := time.Now()
	// log.Println(timestamp)
	_, err = db.Exec("INSERT INTO remind (email,time,ip,obsolutep,password,key) VALUES ($1,$2,$3,$4,$5,$6)", email, timestamp, ip, obsoluteP, params.Password, controlSum)
	if err != nil {
		log.Println(err)
		return false
	}
	err = tmpl.Execute(&body, params)
	if err != nil {
		log.Println(err)
		return false
	}
	//err = mail(email, "БД Мониторинг цен госзакупок <from@monitoring-crm.ru>", "Восстановление доступа к БД \"Мониторинг цен госзакупок\"", body.String())
	err = mail(email, "", "", body.String())
	if err != nil {
		return false
	}

	return true
}

func updPassword(email, key string, db *sql.DB) (string, bool) {
	var password string
	db.QueryRow("SELECT password from remind WHERE email=$1 and key=$2 order by time desc limit 1", email, key).Scan(&password)
	if password == "" {
		log.Println("Updating password: didn't found key+email pair (", email, ",", key, ")")
		return "", false
	}
	h := md5.New()
	h.Write([]byte(password))
	hpassword := hex.EncodeToString(h.Sum(nil))
	db.Exec("UPDATE users SET password=$1 WHERE email=$2;", hpassword, email)
	log.Println("Updated to a new password: ", hpassword, " for user: ", email)
	return password, true
}

type NPasswordParams struct {
	Email    string
	Password string
	Today    string
}

func (p NPasswordParams) newPassword(db *sql.DB) error {
	h := md5.New()
	h.Write([]byte(p.Password))
	hpassword := hex.EncodeToString(h.Sum(nil))
	_, err := db.Exec("UPDATE users set password=$1 where email=$2", hpassword, p.Email)
	if err == nil {
		tmpl, err := template.ParseFiles("templates/new-password.mail")
		if err != nil {
			log.Println(err)
		}
		var body bytes.Buffer
		t := time.Now()
		p.Today = fmt.Sprintf("%02d.%02d.%d\n", t.Day(), t.Month(), t.Year())
		err = tmpl.Execute(&body, p)
		if err != nil {
			log.Println(err)
			return err
		}
		//mail(p.Email, "БД Мониторинг цен госзакупок <from@monitoring-crm.ru>", "Новый пароль доступа к БД \"Мониторинг цен госзакупок\"", body.String())
		mail(p.Email, "", "", body.String())
	}
	return err
}

func GetSession(r *http.Request) *sessions.Session {
	var store = sessions.NewCookieStore([]byte("something-very-secret"))
	session, err := store.Get(r, "user-settings")
	if err != nil {
		log.Println("GetSession ERROR: ", err)
	}
	return session
}

func GetEmail(r *http.Request) string {
	session := GetSession(r)
	if v, ok := session.Values["email"]; ok {
		str, _ := v.(string)
		return str
	} else {
		return ""
	}
}

func listUsers(db *sql.DB) (users []User) {
	rows, err := db.Query("SELECT name,email,permission,created,lifetime,sessionslimit, requestslimit FROM users order by created desc;")
	if err != nil {
		log.Println("listUsers error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Name, &user.Email, &user.Permission, &user.tCreated, &user.tLifeTime, &user.SessionsLimit, &user.tRequestsLimit); err != nil {
			log.Println("scan error on listUsers: ", err)
		} else {
			user.FormatTime()
			users = append(users, user)
		}
	}
	return
}

func rmUser(user string, db *sql.DB) {
	_, err := db.Exec("DELETE FROM users WHERE email=$1;", user)
	if err != nil {
		log.Println("remove user error: ", err)
	}
}

func addUser(user User, db *sql.DB) {
	var query, columns, values string
	if user.Email != "" {
		query = "INSERT INTO users "
		columns = "(email,permission,"
		values = fmt.Sprintf(" VALUES('%s','%d',", user.Email, user.Permission)
	} else {
		log.Println("editUser error: email is null!", query)
		return
	}
	if user.Password != "" {
		columns += "password,"
		h := md5.New()
		h.Write([]byte(user.Password))
		hpassword := hex.EncodeToString(h.Sum(nil))
		values += fmt.Sprintf("'%s',", hpassword)
	}
	columns += "sessionslimit,"
	values += fmt.Sprintf("%d,", user.SessionsLimit)
	// log.Println("SESSIONSLIMIT:", user.SessionsLimit)
	columns += "requestslimit,"
	values += fmt.Sprintf("%d,", user.RequestsLimit)
	// log.Println("REQUESTSLIMIT:", user.RequestsLimit)

	if user.Name != "" {
		columns += "name,"
		values += fmt.Sprintf("'%s',", user.Name)
	}
	if user.LifeTime != "" {
		// log.Println("LIFETIME:", user.LifeTime)
		columns += "lifetime,"
		values += fmt.Sprintf("to_timestamp('%s','DD.MM.YYYY'),", user.LifeTime)
	}
	if user.Created != "" {
		// log.Println("CREATED:", user.Created)
		columns += "created,"
		values += fmt.Sprintf("to_timestamp('%s','DD.MM.YYYY'),", user.Created)
	}
	columns = columns[:len(columns)-1] + ")"
	values = values[:len(values)-1] + ")"
	query += columns + values + ";"

	if _, err := db.Exec(query); err == nil {
		log.Println(err)
	}

	return
}

func getUser(email string, db *sql.DB) (u User) {
	if err := db.QueryRow("SELECT email,password,permission, created,lifetime,name,sessionslimit,requestslimit from users where email=$1", email).Scan(&u.Email, &u.Password, &u.Permission, &u.tCreated, &u.tLifeTime, &u.Name, &u.tSessionsLimit, &u.tRequestsLimit); err != nil {
		log.Println(err)
	}
	u.FormatTime()
	u.Fill()
	return
}

func editUser(user User, db *sql.DB) {
	var query, values string
	if user.Email != "" {
		query = "UPDATE users SET"
		values = fmt.Sprintf(" email='%s', permission=%d,", user.Email, user.Permission)
	} else {
		// log.Println("editUser error: email is null!", query)
		return
	}
	if user.Password != "" {
		h := md5.New()
		h.Write([]byte(user.Password))
		hpassword := hex.EncodeToString(h.Sum(nil))
		values += fmt.Sprintf("password='%s',", hpassword)
	}
	if user.SessionsLimit != 0 {
		values += fmt.Sprintf("sessionslimit=%d,", user.SessionsLimit)
	}
	if user.RequestsLimit != 0 {
		values += fmt.Sprintf("requestslimit=%d,", user.RequestsLimit)
	}
	if user.Name != "" {
		values += fmt.Sprintf("name='%s',", user.Name)
	}
	if user.LifeTime != "" {
		values += fmt.Sprintf("lifetime=to_timestamp('%s','DD.MM.YYYY'),", user.LifeTime)
	}
	if user.Created != "" {
		values += fmt.Sprintf("created=to_timestamp('%s','DD.MM.YYYY'),", user.Created)
	}
	values = values[:len(values)-2] + ")"
	query += values + " WHERE email='" + user.Email + "';"
	if _, err := db.Exec(query); err == nil {
		log.Println(err)
	}
	return
}
