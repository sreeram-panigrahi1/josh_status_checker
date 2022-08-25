package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var URLs = make(map[string]string)

func main() {
	fmt.Println("Launching Server.....")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Printf("Error in initializing server %v", err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		param := r.URL.Query().Get("name")
		if param != "" {
			status, _ := URLs[param]
			fmt.Fprintln(w, fmt.Sprintf("Status of %s: %s", param, status))

		} else {
			for key, val := range URLs {
				fmt.Fprintln(w, key, " ", val)
			}
		}
	case "POST":
		var websites []string
		err := json.NewDecoder(r.Body).Decode(&websites)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Error in post method: %+v", err))
		}

		c := make(chan string)
		for _, url := range websites {
			go checkStatus(url, c)
		}
		for link := range c {
			go func(l string) {
				time.Sleep(5 * time.Second)
				checkStatus(l, c)
			}(link)
		}
	default:
		fmt.Fprint(w, "Error, invalid request")
	}
}
func checkStatus(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, " : Invalid, NOT working")
		URLs[link] = "DOWN"
		c <- link
	} else {
		fmt.Println(link, " : Status 200, Working")
		URLs[link] = "UP"
		c <- link
	}
}
