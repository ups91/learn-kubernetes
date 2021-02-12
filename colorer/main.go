package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//nolint
var LOG *os.File

func main() {
	var err error
	if LOG, err = os.OpenFile("/d/color.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777); err != nil {
		fmt.Println(err)
		LOG = nil
	} else {
		fmt.Fprint(LOG, "Colorer-app started")
		defer LOG.Close()
	}
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
	fmt.Fprintf(LOG, "Request procceded: %s", GetString(url))
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
