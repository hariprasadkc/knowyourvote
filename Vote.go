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

type NetaMapping struct {
	Area         string `json:"assembly_constituency_name"`
	Constituency string `json:"parliament_constituency_name"`
}

type Pincode struct {
	Data struct {
		Constituencies []NetaMapping
	}
}

//func ConstituencyFinder(w http.ResponseWriter, r *http.Request) {
//	keys := r.URL.Query()
//	pincode := keys.Get("pincode")
//	if isValidPinCode(pincode) {
//		if result, ok := pincodes[pincode]; ok {
//			responseString := string(result)
//			fmt.Fprint(w, responseString)
//		}
//	}
//}

func ConstituencyFinder(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	pincode := keys.Get("pincode")
	if isValidPinCode(pincode) {
		resp, err := http.Get("https://api.neta-app.com/v2/constituencies/postal_code?pin=" + pincode)
		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				test := Pincode{}
				err = json.Unmarshal(body, &test)
				if err == nil {
					if len(test.Data.Constituencies) > 0 {
						fmt.Println(test.Data.Constituencies)
						pintmpl.Execute(w, test.Data.Constituencies)
					}
				}
			}
		}
	}
	pinerror.Execute(w, "")
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
var pintmpl *template.Template
var pinerror *template.Template

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
	pintmpl = template.Must(template.ParseFiles("templates/pintable.html"))
	pinerror = template.Must(template.ParseFiles("templates/pinalert.html"))
	r := mux.NewRouter()
	r.HandleFunc("/findconstituency", ConstituencyFinder)
	r.HandleFunc("/getconstituency", getConstituencyDetails)
	r.HandleFunc("/test", test)
	s := http.StripPrefix("/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/").Handler(s)
	http.Handle("/", r)
	//	http.ListenAndServe(":8080", r)
	appengine.Main()

}
