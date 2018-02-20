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
	ControllerProtocol string `json:"controller_protocol"`
	ControllerHost     string `json:"controller_host"`
	ControllerPort     int    `json:"controller_port"`
	User               string `json:"user"`
	Password           string `json:"password"`
	Account            string `json:"account"`
}

var confFile *string
var client *appdrest.Client

func connect(conf Config) {
	client = appdrest.NewClient(conf.ControllerProtocol, conf.ControllerHost, conf.ControllerPort, conf.User, conf.Password, conf.Account)
}

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
	application, err := client.Application.GetApplication(app)
	if err != nil {
		fmt.Println("ERROR - Obtaining application", err)
		return ""
	}
	return strconv.Itoa(application.ID)

}

func redirect(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Recebido:\t ", r.URL.String())

	c := getConfig()
	connect(c)

	queries := r.URL.Query()
	toRedirect := ""

	app := queries["application"]
	if len(app) > 0 {
		appid := getApplicationIDFromURL(app[0])
		newURL := strings.Replace(r.URL.String(), app[0], appid, -1)
		newURL = strings.Replace(newURL, "?", "", 1)

		toRedirect = fmt.Sprintf("%s://%s:%d/controller/#%s", c.ControllerProtocol, c.ControllerHost, c.ControllerPort, newURL)
		fmt.Println("Redirecionando:\t ", toRedirect)

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
