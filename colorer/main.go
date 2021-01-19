package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/job", doJob) // respond
	http.HandleFunc("/", doJob)    // same as /job
	fmt.Printf("Starting server at port 8003\n")
	if err := http.ListenAndServe(":8003", nil); err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
}

func doJob(w http.ResponseWriter, r *http.Request) {
	url := "http://service-timer:8002/"

	fmt.Fprintf(w, "<font color='#ff0000'><b>%s</b></font>", GetString(url))
}

func GetString(url string) string {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return ``
	}

	sl, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(sl)
}
