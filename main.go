package main

import (
	"net/http"
	"log"
	"html/template"
	"encoding/json"
	"strconv"
	
)

var state = []string{"", "", "", "", "", "", "", "", ""}

var player = "X"

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index.html"] = template.Must(template.ParseFiles("static/index.html"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	stateJson, err := json.Marshal(state)

	if err != nil {
		log.Fatal(err)
	}

	tmpl := templates["index.html"]

	tmpl.ExecuteTemplate(w, "index.html", map[string]template.JS{"state": template.JS(stateJson)})
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	stateJson, err := json.Marshal(state)

	action := r.PostFormValue("cell-id")

	log.Default().Println("Action: ", action)

	actionId, err := strconv.Atoi(action)

	tmpl := templates["index.html"]

	if state[actionId] != "" {
		tmpl.ExecuteTemplate(w, "index.html", map[string]template.JS{"state": template.JS(stateJson)})
	} else {
		state[actionId] = player

		log.Default().Println("State: ", state)

		if player == "X" {
			player = "O"
		} else {
			player = "X"
		}

		tmpl.ExecuteTemplate(w, "index.html", map[string]template.JS{"state": template.JS(stateJson)})
	}

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
	//     fmt.Fprintf(w, "getting ti")
	// })

	// fs := http.FileServer(http.Dir("static/"))
	// http.Handle("/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/action", actionHandler)

	log.Fatal(http.ListenAndServe(":7009", nil))
}
