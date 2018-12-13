package main

import (
	"net/http"
	"io"
	"log"
	"encoding/json"
	"bitbucket.org/company-one/tender-one/sphinx-switch"
	"bitbucket.org/company-one/tender-one/postgres-switch"
	"bitbucket.org/company-one/tender-one/tasks"
)

type Response struct {
	Status string `json:"status"`
}

/**
 * Hooks requests handler.
 * 
 * List of hooks:
 * - Switch to use dublicates of sphinx indexes and database
 *   /hook/?action=SwitchSources&toDublicate=<bool>
 * - Run task (shell script)
 *   /hook/?action=RunTask&name=<string>
 */
func hooksHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	params := r.Form
	action := params.Get("action")
		
	jsonResponse, err := json.Marshal(Response{"none"})
	sendJsonResponse := true

	if toDublicate := params.Get("toDublicate"); action == "SwitchSources" && len(toDublicate) > 0 {
	
		sphinxSwitch.Switch(toDublicate == "true")
		psqlSwitch.Switch(toDublicate == "true")
		jsonResponse, err = json.Marshal(Response{"ok"})

		if toDublicate == "true" {
			log.Println("Database and indexes are switched to dublicates")
		} else {
			log.Println("Database and indexes are switched to originals")
		}
	}

	if taskName := params.Get("name"); action == "RunTask" && len(taskName) > 0 {

		w.Header().Set("Content-type", "text/plain")
		sendJsonResponse = false

		task := tasks.GetTask(taskName)
		task.Stdout = w
		
		if err = task.Start(); err != nil {
			log.Println("Error while task starting: ", err)
		}
		
		if err = task.Wait(); err != nil {
			log.Println("Error while task waiting: ", err)
		}
	}

	if sendJsonResponse {

		w.Header().Set("Content-type", "application/json")

		if err != nil {
			log.Println("Hook handler Error: ", err)
		}

		io.WriteString(w, string(jsonResponse))
	}
}
