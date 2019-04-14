// Vote
package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type MP struct {
	Name  string
	Party string
	Link  string
	Wiki  string
}

type CandidateMeta struct {
	Id            string
	Name          string
	Age           int
	Gender        string
	Qualification string
	Phone         string
	Assets        string
	Livelihood    string
	Party         string
}

type Candidate struct {
	Id                 string
	Name               string
	Age                string
	Qualification      string `json:"education_level"`
	Education          string
	Phone              string
	Email              string
	Assets             string
	Livelihood         string `json:"profession"`
	Income             string
	Party              string
	PartySymbol        string `json:"symbol"`
	Gender             string
	ITR16              string
	ITR17              string
	ITR18              string
	PendingCases       string `json:"pending_cases"`
	PendingDescription string `json:"pending_description"`
	PendingIPC         string `json:"pending_ipc"`
	PendingNIPC        string `json:"pending_nipc"`
	ConvictedCases     string `json:"convicted_cases"`
	ConvictedIPC       string `json:"convicted_ipc"`
	ConvictedNIPC      string `json:"convicted_nipc"`
	Liabilities        string
	Political          string `json:"political_background"`
	PoliticalLink      string `json:"political_background_link"`
	Affidavit          string
	Const              string `json:"constituency"`
	ConstituencyName   string `json:"constituency_name"`
}

type Constituency struct {
	Id          string
	Name        string
	Areas       string
	Wiki        string
	Current     MP `json:"mp"`
	Candidates  []CandidateMeta
	Male        string
	Female      string
	Transgender string
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

func ConstituencyFinder(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	keys := r.URL.Query()
	pincode := keys.Get("pincode")
	if isValidPinCode(pincode) {
		resp, err := client.Get("https://api.neta-app.com/v2/constituencies/postal_code?pin=" + pincode)
		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				test := Pincode{}
				err = json.Unmarshal(body, &test)
				if err == nil {
					if test.Data.Constituencies != nil {
						if len(test.Data.Constituencies) > 0 {
							templates["pintable"].Execute(w, test.Data.Constituencies)
							return
						}
					}
				}
			}
		}
	}
	templates["pinerror"].Execute(w, "")
}

func getConstituencyDetails(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	constituency := keys.Get("constituency")
	if result, ok := constituencies[constituency]; ok {
		templates["constituency"].Execute(w, result)
	}
}

func isValidPinCode(s string) bool {
	if len(s) == 6 {
		_, err := strconv.ParseFloat(s, 64)
		return err == nil
	}
	return false
}

func getCandidate(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	candidate := keys.Get("candidate")
	if result, ok := candidates[candidate]; ok {
		templates["candidate"].Execute(w, result)
	}
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["pintable"] = template.Must(template.ParseFiles("templates/pintable.html"))
	templates["pinerror"] = template.Must(template.ParseFiles("templates/pinalert.html"))
	templates["constituency"] = template.Must(template.ParseFiles("templates/layout.html", "templates/constituency.html"))
	templates["candidate"] = template.Must(template.ParseFiles("templates/layout.html", "templates/candidate.html"))
}

func loadJSON() error {
	constituencyJSON, err := os.Open("constituencies.json")
	if err == nil {
		byteValue, _ := ioutil.ReadAll(constituencyJSON)
		defer constituencyJSON.Close()
		err = json.Unmarshal(byteValue, &constituencies)
		if err == nil {
			candidateJSON, err := os.Open("candidates.json")
			if err == nil {
				byteValue, _ := ioutil.ReadAll(candidateJSON)
				err = json.Unmarshal(byteValue, &candidates)
				defer candidateJSON.Close()
			}
		}
	}
	return err
}

var constituencies map[string]Constituency
var candidates map[string]Candidate
var templates map[string]*template.Template

func main() {

	loadTemplates()
	if err := loadJSON(); err != nil {
		panic("Unable to load JSON : " + err.Error())
	}
	r := mux.NewRouter()
	r.HandleFunc("/findconstituency", ConstituencyFinder)
	r.HandleFunc("/getconstituency", getConstituencyDetails)
	r.HandleFunc("/getcandidate", getCandidate)
	s := http.StripPrefix("/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/").Handler(s)
	http.Handle("/", r)
	appengine.Main()
	defer func() {
		if rec := recover(); rec != nil {

		}
	}()
}
