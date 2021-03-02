package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"html/template"

	"github.com/gorilla/mux"
)

var (
	USBC_SITE     = "http://webservices.bowl.com"
	ID            = "/USBC.Search.Services/api/v1/members/id"
	AVERAGES      = "/USBC.Search.Services/api/v1/compositeaverages"
	RERATE        = "/USBC.Search.Services/api/v1/reratedaverage"
	LEAGUES       = "/USBC.Search.Services/api/v1/leagueactivities"
	MEMBERSHIPS   = "/USBC.Search.Services/api/v1/memberships"
	ACHIEVEMENTS  = "/USBC.Search.Services/api/v1/achievements"
	TOURNAMENTS   = "/USBC.Search.Services/api/v1/tournamentaverages"

	//IGBO_SITE     = "http://old.igbo.org/"
	//IGBO_AVERAGES = "/tournaments/get-igbo-tournament-tad-average/"
	//IGBO_ID       = "/tournaments/igbots-id-lookup/"
)

type Results struct {
	Profiles []USBCProfile `json:"results"`
}

type USBCProfile struct {
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
	First string `json:"first"`
	Init string `json:"init"`
	Last string `json:"last"`
	Gender string `json:"gender"`
	PBA bool `json:"pba"`
	Association string `json:"assn"`
}

/*type IGBO struct {
	ID   string
	Name string
	City string
	USBC string
}
type Bowler struct {
	IGBO       string
	Name       string
	City       string
	USBC       string
	TotalPins  string
	TotalGames string
	Average    string
	Events     []Event
} */

/*type Event struct {
	Date       string
	Tournament string
	Type       string
	Series     int
	Games      int
	Average    int
} */

func UsageHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"/usbc/id":           "Get USBC ID: pass in prefix/suffix/size",
		"/usbc/averages":     "Get USBC Averages: pass in prefix/suffix/size",
		"/usbc/rerates":      "Get USBC Rerates: pass in prefix/suffix/size",
		"/usbc/leagues":      "Get USBC Leagues: pass in prefix/suffix/size",
		"/usbc/memberships":  "Get USBC Memberships: pass in prefix/suffix/size",
		"/usbc/achievements": "Get USBC Achievements: pass in prefix/suffix/size",
		"/usbc/tournaments":  "Get USBC Tournaments: pass in prefix/suffix/size",
		"/igbo/id":           "Get IGBO ID: pass in first/last/usbc",
		"/igbo/averages":     "Get IGBO averages: pass in id/yearrange",
	}

	d, _ := json.MarshalIndent(data, "", "    ")
	w.Write(d)
}

func queryForUSBC(url string, r *http.Request) []byte {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	s := r.URL.Query()
	prefix := s["prefix"][0]
	suffix := s["suffix"][0]
	size := s["size"][0]

	q := req.URL.Query()
	q.Add("prefix", prefix)
	q.Add("suffix", suffix)
	q.Add("size", size)
	req.URL.RawQuery = q.Encode()

	data, err := http.Get(req.URL.String())
	body, err := ioutil.ReadAll(data.Body)

	return body
}

func IDHandler(w http.ResponseWriter, r *http.Request) {
	body := queryForUSBC(USBC_SITE+ID, r)
	w.Header().Set("Content-Type", "application/json")

	var results Results
	err := json.Unmarshal(body, &results)
	if err != nil {
		fmt.Println(err)
	}

	profile := results.Profiles[0]

	f, err := os.Create("profiles/" + profile.Prefix + profile.Suffix)
	if err != nil {
		fmt.Println(err)
	}

	t, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(f, profile)
	http.ServeFile(w, r, "profiles/" + profile.Prefix + profile.Suffix)
}

func AverageHandler(w http.ResponseWriter, r *http.Request) {
	body := queryForUSBC(USBC_SITE+AVERAGES, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
func RerateHandler(w http.ResponseWriter, r *http.Request) {
	body := queryForUSBC(USBC_SITE+RERATE, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func LeagueHandler(w http.ResponseWriter, r *http.Request) {
	body := queryForUSBC(USBC_SITE+LEAGUES, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
func MembershipHandler(w http.ResponseWriter, r *http.Request) {
	body := queryForUSBC(USBC_SITE+MEMBERSHIPS, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func AchievementHandler(w http.ResponseWriter, r *http.Request) {
	body := queryForUSBC(USBC_SITE+ACHIEVEMENTS, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func TournamentHandler(w http.ResponseWriter, r *http.Request) {
	body := queryForUSBC(USBC_SITE+TOURNAMENTS, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", UsageHandler)
	r.HandleFunc("/usbc/id", IDHandler)
	r.HandleFunc("/usbc/averages", AverageHandler)
	r.HandleFunc("/usbc/rerates", RerateHandler)
	r.HandleFunc("/usbc/leagues", LeagueHandler)
	r.HandleFunc("/usbc/memberships", MembershipHandler)
	r.HandleFunc("/usbc/achievements", AchievementHandler)
	r.HandleFunc("/usbc/tournaments", TournamentHandler)
	//r.HandleFunc("/igbo/id", IGBOIDHandler)
	//r.HandleFunc("/igbo/averages", IGBOHandler)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
