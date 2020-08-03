package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	SITE = "http://webservices.bowl.com"
	ID = "/USBC.Search.Services/api/v1/members/id"
	AVERAGES = "/USBC.Search.Services/api/v1/compositeaverages"
	RERATE = "/USBC.Search.Services/api/v1/reratedaverage"
	LEAGUES = "/USBC.Search.Services/api/v1/leagueactivities"
	MEMBERSHIPS = "/USBC.Search.Services/api/v1/memberships"
	ACHIEVEMENTS = "/USBC.Search.Services/api/v1/achievements"
	TOURNAMENTS = "/USBC.Search.Services/api/v1/tournamentaverages"
)

func query(url string, r *http.Request) []byte {
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
	body := query(SITE + ID, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func AverageHandler(w http.ResponseWriter, r *http.Request) {
	body := query(SITE + AVERAGES, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
func RerateHandler(w http.ResponseWriter, r *http.Request) {
	body := query(SITE + RERATE, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func LeagueHandler(w http.ResponseWriter, r *http.Request) {
	body := query(SITE + LEAGUES, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
func MembershipHandler(w http.ResponseWriter, r *http.Request) {
	body := query(SITE + MEMBERSHIPS, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func AchievementHandler(w http.ResponseWriter, r *http.Request) {
	body := query(SITE + ACHIEVEMENTS, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func TournamentHandler(w http.ResponseWriter, r *http.Request) {
	body := query(SITE + TOURNAMENTS, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func UsageHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := mux.NewRouter()
	r.Handle("/", UsageHandler)
	r.HandleFunc("/id", IDHandler)
	r.HandleFunc("/averages", AverageHandler)
	r.HandleFunc("/rerates", RerateHandler)
	r.HandleFunc("/leagues", LeagueHandler)
	r.HandleFunc("/memberships", MembershipHandler)
	r.HandleFunc("/achievements", AchievementHandler)
	r.HandleFunc("/tournaments", TournamentHandler)
	log.Fatal(http.ListenAndServe(":" + port, r))
}