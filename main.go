package main

import (
	"net/http"
)

func main() {
	// http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
	//     fmt.Fprintf(w, "getting ti")
	// })

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":7009", nil)
}
