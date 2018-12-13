package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	"strings"
	"bytes"
	"io"
	"os/exec"

	cbrate "bitbucket.org/company-one/cbrrate"
	"bitbucket.org/company-one/tender-one/cbr"
	"bitbucket.org/company-one/tender-one/user"
	"bitbucket.org/company-one/tender-one/postgres-driver"
	"bitbucket.org/company-one/tender-one/postgres-switch"
	"bitbucket.org/company-one/tender-one/api"
	"bitbucket.org/company-one/tender-one/tasks"
	"bitbucket.org/company-one/tender-one/state"
	"bitbucket.org/company-one/tender-one/sphinx-switch"
)

const MaxDocs int = 100000
const perPageDefault int = 25

var Page int = 1000
var dbM *sql.DB
var USDRate []cbrate.CurrencyRate
var flagHookIP string

type JSON_request struct {
	success	string	`json:"success"`
	message	string	`json:"message"`
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Error    string `json:"error"`
}

type RemindParams struct {
	Email     string `json:"email"`
	Key       string `json:"key"`
	Password  string `json:"password"`
	Error     string `json:"error"`
	Today     string
	EmailSent bool
}

type RegParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Error    string `json:"error"`
	Name     string `json:"name"`
}

func canI(w http.ResponseWriter, r *http.Request, db *sql.DB, lvl int) bool {
	session := GetSession(r)

	return isAuthorized(db, session) && checkLVL(lvl, db, session) && rulersChk(r, db)
}

func canHook(r *http.Request) bool {
	return getIP(r) == flagHookIP
}

func getParams(r *http.Request) (params url.Values) {
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	if r.Method == "GET" {
		params = r.URL.Query()
	}
	return
}

func getParamsPostOnly(r *http.Request) (params url.Values) {
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	return
}

func getParamsGetOnly(r *http.Request) (params url.Values) {
	if r.Method == "GET" {
		params = r.URL.Query()
	}
	return
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	params := url.Values{}
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
	}
	if r.Method == "GET" {
		params = r.URL.Query()
	}
	tmpl, err := template.ParseFiles("templates/new-password.html")
	if err != nil {
		log.Println(err)
	}
	var paramsData NPasswordParams
	if paramsData.Password = params.Get("new"); paramsData.Password == "" {
		err := tmpl.Execute(w, paramsData)
		log.Println(paramsData, err)
		return
	}
	paramsData.Email = GetEmail(r)
	if err := paramsData.newPassword(db); err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/signin", http.StatusFound)

}
func requestHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := url.Values{}
	//var js []byte
	var my_Info JSON_request
	w.Header().Set("Content-Type", "application/json")

	log.Println("Метод: ", r.Method)
	log.Println("ALL: ", r.Form)

	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form

		log.Println("Организация: ", params.Get("company"))
		log.Println("Имя: ", params.Get("name"))
		log.Println("Инфо: ", params.Get("info"))
		log.Println("Тариф: ", params.Get("tariff"))

		message := fmt.Sprintf(
		`Организация: %s
			Имя: %s
			Контактная информация: %s
			Тариф: %s`, params.Get("company"), params.Get("name"), params.Get("info"), params.Get("tariff"))

		err := sendMail("vasya0107@mail.ru", "info@tender-one.ru", "Запрос на доступ к сервису", message)

		if  err != nil {
			my_Info = JSON_request{"false", err.Error()}
		} else {
			my_Info = JSON_request{"true", "Все хорошо!"}
		}

		log.Println("my_Info: ", my_Info)
		js := []byte(fmt.Sprintf(`{"success":%s}`, my_Info.success))
		log.Println("js: ", string(js))

		if err != nil {
			log.Println("ERROR Marshal")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	} else {
		js := []byte(`{"success":false}`)
		w.Write(js)
	}
}
func remindHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := url.Values{}
	var paramsData RemindParams
	tmpl, err := template.ParseFiles("templates/restore.html")
	if err != nil {
		log.Println(err)
	}

	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form

		paramsData.Error = params.Get("error")
		if paramsData.Email = params.Get("email"); paramsData.Email == "" {
			err := tmpl.Execute(w, paramsData)
			log.Println(paramsData, err)
			return
		}
		log.Println("Remind:", paramsData)
		if paramsData.Key = params.Get("key"); paramsData.Key == "" {
			if ok := getForgotKey(paramsData.Email, getIP(r), db); !ok {
				log.Println("Generating key:", paramsData, ok)
				paramsData.Error = "Произошла ошибка при восстановлении пароля."
			} else {
				paramsData.Error = "Письмо отправлено! Следуйте инструкциям в письме."
				paramsData.EmailSent = true
			}
			err := tmpl.Execute(w, paramsData)
			log.Println(paramsData, err)
		} else {
			var ok bool
			if paramsData.Password, ok = updPassword(paramsData.Email, paramsData.Key, db); !ok {
				paramsData.Error = "Произошла ошибка при восстановлении пароля."
				tmpl.Execute(w, paramsData)
			} else {
				session := GetSession(r)
				if login(paramsData.Email, paramsData.Password, getIP(r), db, session) {
					err := session.Save(r, w)
					log.Println(err)
					http.Redirect(w, r, "/auctions223", http.StatusFound)
				} else {
					err := tmpl.Execute(w, paramsData)
					log.Println(paramsData, err)
				}
			}
		}
	} else {
		tmpl.Execute(w, paramsData)
	}
}

func oldDBHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//http://online.monitoring-crm.ru/admin/index.php?Login=l.kalneus@crm54.com&Password=gdtPpDs
	/*
	var store = sessions.NewCookieStore([]byte("something-very-secret"))
	session, _ := store.Get(r, "user-settings")
	if !isAuthorized(db, session) {
		http.Redirect(w, r, "/signin", http.StatusFound)
	} else {
		log.Println(session.Values["f"])
		s := fmt.Sprintf("http://online.monitoring-crm.ru/admin/index.php?Login=%s&Password=%s", session.Values["email"], session.Values["f"])
		http.Redirect(w, r, s, http.StatusFound)
	}
	*/
}

type LandingResult struct {
	IsLogin	bool
}

func LandingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/landing.html")
	if err != nil {
		log.Println(err)
	}

	var result LandingResult
	result.IsLogin = canI(w, r, db, 0)

	tmpl.Execute(w, result)
}

func hiHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := url.Values{}
	var paramsData LoginParams
	tmpl, err := template.ParseFiles("templates/signin.html")
	if err != nil {
		log.Println(err)
	}

	if r.Method == "POST" {
		//tmpl := template.Must(template.New("login").Parse(HiPage))
		//tmpl.Execute(w, result)
		r.ParseForm()
		params = r.Form
		// log.Println(params)
		paramsData.Email = params.Get("email")
		paramsData.Password = params.Get("password")
		var store = sessions.NewCookieStore([]byte("something-very-secret"))
		session, _ := store.Get(r, "user-settings")
		if isAuthorized(db, session) && canI(w, r, db, 0) {
			http.Redirect(w, r, "/auctions223", http.StatusFound)
		}
		if paramsData.Email != "" && paramsData.Password != "" {
			if login(paramsData.Email, paramsData.Password, getIP(r), db, session) {
				log.Println("Testing session:", session.Values)
				err := session.Save(r, w)
				log.Println(err)
				http.Redirect(w, r, "/auctions223", http.StatusFound)
			} else {
				paramsData.Error = "Неверный Email или Пароль!"
				tmpl.Execute(w, paramsData)
			}
		} else {
			var store = sessions.NewCookieStore([]byte("something-very-secret"))
			session, _ := store.Get(r, "user-settings")
			if login(paramsData.Email, paramsData.Password, getIP(r), db, session) {
				err := session.Save(r, w)
				log.Println(err)
				http.Redirect(w, r, "/auctions223", http.StatusFound)
			} else {
				tmpl.Execute(w, paramsData)
			}
		}
		//tmpl.Execute(w, paramsData)
		//fmt.Fprintf(w, "%s", HiPage)

	} else {
		tmpl.Execute(w, paramsData)
	}
	//	http.Redirect(w, r, "/notifications", http.StatusOK) // Temporary moving to notifications page
}

func byeHandler(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) {
	session, _ := store.Get(r, "user-settings")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/signin/", http.StatusFound) // Temporary moving to notifications page
}
func testHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var paramsData Notification223Result
	if 1>2 {
		paramsData.NavType = "notifications223"
	}
	//tmpl := template.Must(template.New("").Delims("{{", "}}").ParseFiles("templates/demo.html"))
	tmpl, err := template.ParseFiles("templates/tests.html")
	//err := tmpl.ExecuteTemplate(w, "notifications223", nil)
	tmpl.Execute(w, paramsData)
	//*
	if err != nil {
		log.Println(err)
	}
	/**/
}
func regHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := url.Values{}
	var paramsData RegParams
	tmpl, err := template.ParseFiles("templates/demo.html")
	if err != nil {
		log.Println(err)
	}
	if r.Method == "POST" {
		r.ParseForm()
		params = r.Form
		if paramsData.Email = params.Get("email"); paramsData.Email == "" {
			paramsData.Error = "Необходимо указать свой email."
			log.Println(paramsData.Error)
			tmpl.Execute(w, paramsData)
			return
		}
		if paramsData.Password = params.Get("password"); paramsData.Password == "" {
			paramsData.Error = "Необходимо указать пароль."
			tmpl.Execute(w, paramsData)
			log.Println(paramsData.Error)
			return
		}
		if paramsData.Name = params.Get("name"); paramsData.Name == "" {
			paramsData.Error = "Необходимо указать Ваше имя."
			tmpl.Execute(w, paramsData)
			log.Println(paramsData.Error)
			return
		}
		if ok := signUp(paramsData.Email, paramsData.Password, getIP(r), paramsData.Name, permissionDemo, db); !ok {
			paramsData.Error = "Произошла ошибка при регистрации. Попробуйте ещё раз."
			tmpl.Execute(w, paramsData)
			log.Println(paramsData.Error)
			return
		} else {
			var store = sessions.NewCookieStore([]byte("something-very-secret"))
			session, _ := store.Get(r, "user-settings")
			if login(paramsData.Email, paramsData.Password, getIP(r), db, session) {
				err := session.Save(r, w)
				log.Println(err)
				http.Redirect(w, r, "/auctions223", http.StatusFound)
			}
		}
	} else {
		tmpl.Execute(w, paramsData)
	}
}


func from_mail(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	params := url.Values{}
	var email_id string


	r.ParseForm()
	params = r.Form

	email_id = params.Get("e")

	var headers string
	for k, v := range r.Header {
		headers += fmt.Sprintf("%s: %s\n", k, v)
	}

	/*log.Println(r.RemoteAddr)
	log.Println("\n!!!\n")
	for k, v := range r.Trailer {
		log.Println(k, v, r.Trailer.Get(k))
	}
	log.Println("\n!!!\n")
	for k, v := range r.Header {
		log.Println(k, v, r.Header.Get(k))
	}*/

	_, err := db.Exec(`
		INSERT INTO email_send_statistics (send_id, click_link, headers, remote_ip, remote_user, email)
        VALUES (1, 'single_link', $1, $2, null, $3)
	`, headers, r.RemoteAddr, email_id)

	if err != nil {
		log.Fatal("Save from_mail INSERT:   ", err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func AppendShards(slice []string, elements []string) []string {
	n := len(slice)
	total := len(slice) + len(elements)
	if total > cap(slice) {
		// Reallocate. Grow to 1.5 times the new size, so we can still grow.
		newSize := total*3/2 + 1
		newSlice := make([]string, total, newSize)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[:total]
	copy(slice[n:], elements)
	return slice
}

func execTest() {
	cwd, _ := os.Getwd()
	ppid := os.Getppid()

	ps := exec.Command("ps")
	grep := exec.Command("grep", strconv.Itoa(ppid))

	r, w := io.Pipe()
	ps.Stdout = w
	grep.Stdin = r

	var output bytes.Buffer
	grep.Stdout = &output

	ps.Start()
	grep.Start()
	ps.Wait()
	w.Close()
	grep.Wait()

	pcmd := strings.Replace(strings.Trim(output.String(), "\t\n "), "\n", " | ", -1)

	log.Println("[debug] Arguments:", os.Args)
	log.Println("[debug] Cwd:", cwd)
	log.Println("[debug] Parent pid:", ppid)
	log.Println("[debug] Parent cmd:", pcmd)
}

var chttp = http.NewServeMux()
var mhttp = http.NewServeMux()

func main() {

	var flagFile string
	var flagTasksRoot string
	var flagStateStorage string
	var flagDBCons int
	var flagPort int
	var flagNoLog bool
	var flagUseSourcesDub bool
	var flagUseSourcesOrig bool

	flag.BoolVar(&flagNoLog, "nolog", false, "Turn off log2file")
	flag.BoolVar(&flagUseSourcesDub, "use-sources-dub", false, "Force use sphinx index and database dublicates")
	flag.BoolVar(&flagUseSourcesOrig, "use-sources-orig", false, "Force use sphinx index and database originals")
	flag.StringVar(&flagFile, "log2file", "/var/log/tender-one.log", "Record logs to file")
	flag.StringVar(&flagHookIP, "hook-ip", "127.0.0.5", "Allow hook requests from this ip")
	flag.StringVar(&flagTasksRoot, "tasks-root", "./", "Root directory with tasks scripts")
	flag.StringVar(&flagStateStorage, "state-storage", "/var/tender-one", "Directory to store files with states")
	flag.IntVar(&flagDBCons, "dbcons", 1, "Count of database connections")
	flag.IntVar(&flagPort, "port", 80, "Bind to custom port.")
	flag.Parse()

	if !flagNoLog && flagFile != "" {
			
		logFile, err := os.OpenFile(flagFile, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		defer logFile.Close()
		
		log.Println("Writing log to file")
		log.SetOutput(logFile)
	}
	
	execTest()
	cbr.InitScheduler(&USDRate)
	//mailer.Init("smtp.yandex.ru", 465, "bot@gkcrm.ru", "20172017")
	tasks.SetTasksRoot(flagTasksRoot)
	state.Load(flagStateStorage)
	sphinxSwitch.LoadState()
	psqlSwitch.LoadState()

	if flagUseSourcesDub != flagUseSourcesOrig {
		sphinxSwitch.Switch(flagUseSourcesDub)
		psqlSwitch.Switch(flagUseSourcesDub)
	}

	sphinxSwitch.Init("127.0.0.1", "xxx.xxx.xxx.xxx")

	dbSite, _ := psql.Init(&psql.Database{"X", "X", "localhost", "", "X", flagDBCons, flagDBCons})
	defer dbSite.Close()

	dbM, _ = psqlSwitch.Init(
              &psql.Database{"x", "x", "localhost", "", "x", flagDBCons, flagDBCons},
	      &psql.Database{"x", "x", "localhost", "", "x", flagDBCons, flagDBCons},
        )
	defer dbM.Close()
	
	db,_ := psqlSwitch.Init(
              &psql.Database{"x", "x", "localhost", "", "x", flagDBCons, flagDBCons},
              &psql.Database{"x", "x", "localhost", "", "x", flagDBCons, flagDBCons},
        )
        defer dbM.Close()
	defer db.Close()

	user.InitScheduler(dbSite)

	var store = sessions.NewCookieStore([]byte("something-very-secret"))

	static_dir, _ := os.Getwd()
	static_dir += "/static/"

	chttp.Handle("/", http.FileServer(http.Dir(static_dir)))

	mhttp.HandleFunc("/signin/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 0) {
			http.Redirect(w, r, "/auctions223", http.StatusFound)
			//notifications223Handler(w, r, db, dbSite)
		} else {
			//hiHandler(w, r, dbSite)
			//http.Redirect(w, r, "/signin", http.StatusFound)
			hiHandler(w, r, dbSite)
		}


	}) // Login
	mhttp.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		byeHandler(w, r, store)
	}) // Logout
	mhttp.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("ЕСТЬ ПЕРЕХОД!!!!!!!!!!!!!!!!!!!!!!!!!!!")

		requestHandler(w, r, dbSite)
	}) // Send Mail
	mhttp.HandleFunc("/restore/", func(w http.ResponseWriter, r *http.Request) {
		remindHandler(w, r, dbSite)
	}) // restore forgotten password
	mhttp.HandleFunc("/demo/", func(w http.ResponseWriter, r *http.Request) {
		regHandler(w, r, dbSite)
	}) // Logout
	mhttp.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		testHandler(w, r, dbSite)
	}) // Logout
	mhttp.HandleFunc("/user/password/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 0) {
			changePasswordHandler(w, r, dbSite)
		} else {
			//hiHandler(w, r, dbSite)
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
	}) // Change Password
	mhttp.HandleFunc("/old", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Redirecting to old databse...")
		//oldDBHandler(w, r, dbSite)
	})
	mhttp.HandleFunc("/auctions/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 0) {
			notificationsHandler(w, r, db, dbSite)
		} else {
			//hiHandler(w, r, dbSite)
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
		// ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		// log.Println("Customer IP: ", ip)
		// log.Println("X-Forwarded-For :" + r.Header.Get("X-FORWARDED-FOR"))
	})
	mhttp.HandleFunc("/auctions223/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 0) {
			notifications223Handler(w, r, db, dbSite)
		} else {
			//hiHandler(w, r, dbSite)
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
		// ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		// log.Println("Customer IP: ", ip)
		// log.Println("X-Forwarded-For :" + r.Header.Get("X-FORWARDED-FOR"))
	})
	mhttp.HandleFunc("/contracts/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 0) {
			contractsHandler(w, r, db, dbSite, "hope_products", false)
		} else {
			//hiHandler(w, r, dbSite)
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
		// ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		// log.Println("Customer IP: ", ip)
		// log.Println("X-Forwarded-For :" + r.Header.Get("X-FORWARDED-FOR"))
	})
	mhttp.HandleFunc("/org/", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("HEY")
		if canI(w, r, dbSite, 0) {
			customerSearchHandler(w, r, db)
		} else {
			hiHandler(w, r, dbSite)
		}
	})
	mhttp.HandleFunc("/supplier/", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("HEY")
		if canI(w, r, dbSite, 0) {
			supplierSearchHandler(w, r, db)
		} else {
			hiHandler(w, r, dbSite)
		}
	})
	mhttp.HandleFunc("/okpd/", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("HEY")
		if canI(w, r, dbSite, 0) {
			okpdSearchHandler(w, r, db)
		} else {
			hiHandler(w, r, dbSite)
		}
	})
	mhttp.HandleFunc("/api/okpd/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 0) {
			api.OKPDHandler(w, r, db)
		} else {
			hiHandler(w, r, dbSite)
		}
	})
	mhttp.HandleFunc("/api/okpd-children/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 0) {
			api.OKPDChildrenHandler(w, r, db)
		} else {
			hiHandler(w, r, dbSite)
		}
	})
	mhttp.HandleFunc("/templates/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 0) {
			templatesHandler(w, r, dbSite)
		} else {
			//hiHandler(w, r, dbSite)
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
	})

	mhttp.HandleFunc("/request/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 50) {
			// log.Println("Admin mode!")
			controlUsersHandler(w, r, dbSite)
		} else {
			//hiHandler(w, r, dbSite)
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
	})

	mhttp.HandleFunc("/control/", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 50) {
			// log.Println("Admin mode!")
			controlUsersHandler(w, r, dbSite)
		} else {
			//hiHandler(w, r, dbSite)
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
	})
	mhttp.HandleFunc("/control/a/users/list", func(w http.ResponseWriter, r *http.Request) {
		if canI(w, r, dbSite, 50) {
			// log.Println("API Admin mode!")
			listUsersHandler(w, r, dbSite)
		} else {
			//hiHandler(w, r, dbSite)
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
	})
	// EMAIL CLICK
	mhttp.HandleFunc("/from_mail", func(w http.ResponseWriter, r *http.Request) {
		from_mail(w, r, dbSite)
	})




	mhttp.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//log.Println(r.URL.Path)
		if r.URL.Path == "/" {
//			if canI(w, r, dbSite, 0) {
//				notifications223Handler(w, r, db, dbSite)
				//contractsHandler(w, r, db, dbSite, "hope_products", false) // indexHandler should be here
//			} else {
				//hiHandler(w, r, dbSite)
				//http.Redirect(w, r, "/signin", http.StatusFound)
				LandingHandler(w, r, dbSite) // Lending
//			}
		} else {
			chttp.ServeHTTP(w, r)
		}
	})

	/**
	 * Hooks
	 */

	mhttp.HandleFunc("/hook/", func(w http.ResponseWriter, r *http.Request) {

		if canHook(r) {
			hooksHandler(w, r)
		} else {
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
	})

	log.Println("Started at localhost:", flagPort)
	log.Println("CBRTEST:", fmt.Sprintf("%02d.%02d.%d", time.Now().Day(), time.Now().Month(), time.Now().Year()), ":", cbrate.GetCurrencyRate(fmt.Sprintf("%d.%d.%d", time.Now().Day(), time.Now().Month(), time.Now().Year()), USDRate))
	//http.ListenAndServe(":"+strconv.Itoa(flagPort), nil)

	l, err := net.Listen("tcp4", ":"+strconv.Itoa(flagPort))
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(l, context.ClearHandler(mhttp))

}
