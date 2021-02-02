package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

//nolint
var LOG *os.File

func main() {
	var err error
	if LOG, err = os.OpenFile("/d/echo.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777); err != nil {
		fmt.Println(err)
		LOG = nil
	} else {
		fmt.Fprint(LOG, "Echo-app started")
		defer LOG.Close()
	}
	http.HandleFunc("/ping", pingHndl)     // respond just "pong"
	http.HandleFunc("/var", varHndl)       // respond All ENV variables or only one named variable
	http.HandleFunc("/shutdown", shutdown) // immediately exit
	http.HandleFunc("/hello", hello)       // get WORLD_NAME env var and greeting.
	http.HandleFunc("/", hello)            // same as /hello
	fmt.Printf("Starting server at port 8001\n")
	if err = http.ListenAndServe(":8001", nil); err != nil {
		// log.Fatal(err)
		fmt.Fprintf(LOG, "Echo-app fatal err: %s", err)
		fmt.Println(err)
	}
	fmt.Fprint(LOG, "Echo-app stopped")

}

func shutdown(w http.ResponseWriter, r *http.Request) { os.Exit(1) }

func hello(w http.ResponseWriter, r *http.Request) {
	WW := os.Getenv("WORLD_NAME")
	if WW == `` {
		WW = `___NONAME___`
	}
	fmt.Fprintf(w, "Hello, %s!", WW)
}

func varHndl(w http.ResponseWriter, r *http.Request) {
	varToGet := r.URL.Query()

	if len(varToGet) > 0 {
		for i := range varToGet {
			fmt.Fprintf(w, "%s = %s\n\n", i, os.Getenv(strings.TrimSpace(i)))
			fmt.Fprintf(LOG, "Echo-app data: %s = %s\n\n", i, os.Getenv(strings.TrimSpace(i)))

		}
		return
	}
	for _, e := range os.Environ() {
		fmt.Fprintln(w, e)
	}
}

func pingHndl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong!")
}
