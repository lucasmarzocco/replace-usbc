package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"

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

func IDHandler(w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest("GET", SITE + ID, nil)
	if err != nil {
		log.Print(err)
		log.Fatal(err)
	}

	s := r.URL.Query()

	prefix := s["prefix"][0]
	suffix := s["suffix"][0]

	q := req.URL.Query()
	q.Add("prefix", prefix)
	q.Add("suffix", suffix)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	data, err := http.Get(req.URL.String())

	body, err := ioutil.ReadAll(data.Body)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func AverageHandler(w http.ResponseWriter, r *http.Request) {


}
func RerateHandler(w http.ResponseWriter, r *http.Request) {


}

func LeagueHandler(w http.ResponseWriter, r *http.Request) {


}
func MembershipHandler(w http.ResponseWriter, r *http.Request) {


}

func AchievementHandler(w http.ResponseWriter, r *http.Request) {


}

func TournamentHandler(w http.ResponseWriter, r *http.Request) {


}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/id", IDHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}