package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dlopes7/go-appdynamics-rest-api/appdrest"
)

// Config contains detail for the controller
type Config struct {
	ControllerURL string `json:"controller_url"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Account       string `json:"account"`
}

var confFile *string
var client = appdrest.NewClient("http", "demo2.appdynamics.com", 80, "resty", "Hwkh718b", "customer1")

func getConfig() Config {
	raw, err := ioutil.ReadFile(*confFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Config
	json.Unmarshal(raw, &c)
	return c
}

func getApplicationIDFromURL(app string) string {
	fmt.Println(app)
	application, err := client.Application.GetApplication(app)
	if err != nil {
		fmt.Println("Error obtaining application")
	}
	return strconv.Itoa(application.ID)

}

func redirect(w http.ResponseWriter, r *http.Request) {

	c := getConfig()

	queries := r.URL.Query()
	toRedirect := ""

	app := queries["application"]
	if len(app) > 0 {
		appid := getApplicationIDFromURL(app[0])
		newURL := strings.Replace(r.URL.String(), app[0], appid, -1)

		toRedirect := fmt.Sprintf("%s/controller/#%s\n", c.ControllerURL, newURL)
		fmt.Println("Redirecionando para: ", toRedirect)

	}
	http.Redirect(w, r, toRedirect, 301)
}

func main() {

	confFile = flag.String("c", "conf.json", "The config file to use")
	flag.Parse()

	http.HandleFunc("/", redirect)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
