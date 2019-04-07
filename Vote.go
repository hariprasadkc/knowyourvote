// Vote
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"log"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

type MP struct {
	Name  string
	Party string
	Link  string
	Wiki  string
}

type Candidate struct {
	Id            string
	Name          string
	Age           int
	Qualification string
	Phone         string
	Assets        string
	Livelihood    string
	Party         string
}

type Constituency struct {
	Id         string
	Name       string
	Areas      []string
	Current    MP `json:"mp"`
	Candidates []Candidate
}

type Constituencies struct {
	Constituencies map[string]Constituency
}

func ConstituencyFinder(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	pincode := keys.Get("pincode")
	if isValidPinCode(pincode) {
		if result, ok := pincodes[pincode]; ok {
			log.Println("allclear")
			responseString := string(result)
			fmt.Fprint(w, responseString)
		}
	}
}

func getConstituencyDetails(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	constituency := keys.Get("constituency")
	if result, ok := constituencies[constituency]; ok {
		log.Println("allclear")
		response, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Fprint(w, string(response))
	}
}

func isValidPinCode(s string) bool {
	if len(s) == 6 {
		_, err := strconv.ParseFloat(s, 64)
		return err == nil
	}
	return false
}

func test(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, constituencies["con-1"])
}

var pincodes map[string]string
var constituencies map[string]Constituency
var tmpl *template.Template

func main() {
	log.Println("Hello World!")
	constituencyJSON, err := os.Open("constituencies.json")
	pincodeJSON, err := os.Open("pincode.json")

	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully Opened constituencies.json")
	defer constituencyJSON.Close()

	byteValue, _ := ioutil.ReadAll(constituencyJSON)
	var temp Constituencies
	err = json.Unmarshal(byteValue, &temp)
	if err != nil {
		log.Println("parse error")
		log.Println(err)
	}

	constituencies = temp.Constituencies
	byteValue, _ = ioutil.ReadAll(pincodeJSON)
	err = json.Unmarshal(byteValue, &pincodes)
	if err != nil {
		log.Println("parse error")
		log.Println(err)

	}

	tmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/constituency.html"))

	r := mux.NewRouter()
	r.HandleFunc("/findconstituency", ConstituencyFinder)
	r.HandleFunc("/getconstituency", getConstituencyDetails)
	r.HandleFunc("/test", test)
	s := http.StripPrefix("/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/").Handler(s)
	http.Handle("/", r)
	//http.ListenAndServe(":8080", r)
	appengine.Main()

}
