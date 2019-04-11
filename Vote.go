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
	Age           string
	Qualification string
	Phone         string
	Assets        string
	Livelihood    string
	Party         string
}

type Candidate struct {
	Id               string
	Name             string
	Age              string
	Qualification    string
	Phone            string
	Assets           string
	Livelihood       string
	Party            string
	Gender           string
	ITR              []string
	Cases            string `json:"cases_filed"`
	Convictions      string
	MovableAssets    string `json:"movable_assets"`
	ImmovableAssets  string `json:"immovable_assets"`
	Liabilities      string
	Political        string `json:"political_background"`
	PoliticalLink    string `json:"political_background_link"`
	Affidavit        string
	Const            string `json:"constituency"`
	ConstituencyName string `json:"constituency_name"`
}

type Constituency struct {
	Id         string
	Name       string
	Areas      []string
	Current    MP `json:"mp"`
	Candidates []CandidateMeta
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
							pintmpl.Execute(w, test.Data.Constituencies)
							return
						}
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
		tmpl.Execute(w, result)
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
		caninfo.Execute(w, result)
	}
}

var pincodes map[string]string
var constituencies map[string]Constituency
var tmpl *template.Template
var pintmpl *template.Template
var pinerror *template.Template
var caninfo *template.Template
var candidates map[string]Candidate

func main() {
	constituencyJSON, err := os.Open("constituencies.json")
	pincodeJSON, err := os.Open("pincode.json")
	candidateJSON, err := os.Open("candidates.json")

	if err != nil {
		//		log.Println(err)
	}

	//	log.Println("Successfully Opened constituencies.json")
	defer constituencyJSON.Close()

	byteValue, _ := ioutil.ReadAll(constituencyJSON)
	var temp Constituencies
	err = json.Unmarshal(byteValue, &temp)
	if err != nil {
		//		log.Println(err)
	}

	byteValue, _ = ioutil.ReadAll(candidateJSON)
	err = json.Unmarshal(byteValue, &candidates)

	constituencies = temp.Constituencies
	byteValue, _ = ioutil.ReadAll(pincodeJSON)
	err = json.Unmarshal(byteValue, &pincodes)
	if err != nil {
		//		log.Println(err)
	}

	tmpl = template.Must(template.ParseFiles("templates/layout.html", "templates/constituency.html"))
	pintmpl = template.Must(template.ParseFiles("templates/pintable.html"))
	pinerror = template.Must(template.ParseFiles("templates/pinalert.html"))
	caninfo = template.Must(template.ParseFiles("templates/layout.html", "templates/candidate.html"))
	r := mux.NewRouter()
	r.HandleFunc("/findconstituency", ConstituencyFinder)
	r.HandleFunc("/getconstituency", getConstituencyDetails)
	r.HandleFunc("/getcandidate", getCandidate)
	s := http.StripPrefix("/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/").Handler(s)
	http.Handle("/", r)
	//	http.ListenAndServe(":8080", r)
	appengine.Main()

}
