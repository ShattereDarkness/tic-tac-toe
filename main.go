package main

import (
	"net/http"
	"log"
	"html/template"
	"encoding/json"
	"strconv"
	
)

var state = []string{"", "", "", "", "", "", "", "", ""}

var winStates = [][]int{{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6}}

var player = "X"

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index.html"] = template.Must(template.ParseFiles("static/index.html"))
	// templates["styles.css"] = template.Must(template.ParseFiles("static/styles.css"))
}

func checkWin() string {
	for _, winState := range winStates {
		if state[winState[0]] == state[winState[1]] && state[winState[1]] == state[winState[2]] {
			if state[winState[0]] != "" {
				return state[winState[0]]
			}
		}
	}

	// Check for draw
    isDraw := true
    for _, cell := range state {
        if cell == "" {
            isDraw = false
            break
        }
    }

    if isDraw {
        return "D"
    }

    return "N" // No winner yet
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	stateJson, err := json.Marshal(state)

	if err != nil {
		log.Fatal(err)
	}

	tmpl := templates["index.html"]
	data := map[string]interface{}{
		"state":   template.JS(stateJson),
		"winnerMessage": "\"\"",
	}
	tmpl.Execute(w, data)

	// tmpl.ExecuteTemplate(w, "index.html", map[string]template.JS{"state": template.JS(stateJson), "winner": template.JS("N")})
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	action := r.PostFormValue("cell-id")

	log.Default().Println("Action: ", action)

	actionId, err := strconv.Atoi(action)

	winner := checkWin()

	log.Println("Winner:", winner)

	if winner == "N" {
		if state[actionId] == "" {
			state[actionId] = player

			if player == "X" {
				player = "O"
			} else {
				player = "X"
			}
		}
	}

	winner = checkWin()

	stateJson, err := json.Marshal(state)

	log.Default().Println("State: ", state)

	if err != nil {
		log.Println("Error marshaling state:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	winnerMessage := ""

	if winner == "D" {
		winnerMessage = "Both of you suck"
	} else if winner != "N" {
		winnerMessage = "Player " + winner + " is the G!"
	}

	tmpl := templates["index.html"]
	data := map[string]interface{}{
		"state":   template.JS(stateJson),
		"winnerMessage": "\"" + winnerMessage + "\"",
	}
	tmpl.Execute(w, data)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	// Reset game state
	winner := checkWin()
	player = "X"

	if winner != "D" && winner != "N" {
		player = winner
	}

	state = []string{"", "", "", "", "", "", "", "", ""}

	indexHandler(w, r)
}


func main() {
	// http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
	//     fmt.Fprintf(w, "getting ti")
	// })

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/action", actionHandler)
	http.HandleFunc("/reset", resetHandler)

	log.Fatal(http.ListenAndServe(":7009", nil))
}
