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
	// templates["styles.css"] = template.Must(template.ParseFiles("static/styles.css"))
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
	action := r.PostFormValue("cell-id")

	log.Default().Println("Action: ", action)

	actionId, err := strconv.Atoi(action)

	if state[actionId] == "" {
		state[actionId] = player

		if player == "X" {
			player = "O"
		} else {
			player = "X"
		}
	}
	
	stateJson, err := json.Marshal(state)


	log.Default().Println("State: ", state)

	if err != nil {
		log.Println("Error marshaling state:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("static/index.html"))
	tmpl.Execute(w, map[string]template.JS{"state": template.JS(stateJson)})
}

func main() {
	// http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
	//     fmt.Fprintf(w, "getting ti")
	// })

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/action", actionHandler)

	log.Fatal(http.ListenAndServe(":7009", nil))
}
