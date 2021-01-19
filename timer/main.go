package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {

	http.HandleFunc("/shutdown", shutdown) // immediately exit
	http.HandleFunc("/", sayTime)          // output date-time
	fmt.Printf("Starting server at port 8002\n")
	if err := http.ListenAndServe(":8002", nil); err != nil {
		fmt.Println(err)
	}
}

func shutdown(w http.ResponseWriter, r *http.Request) { os.Exit(1) }

func sayTime(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, time.Now().String())
}
