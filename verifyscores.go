package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	USBC_SITE     = "http://webservices.bowl.com"
	IGBO_SITE     = "http://old.igbo.org/"
	ID            = "/USBC.Search.Services/api/v1/members/id"
	AVERAGES      = "/USBC.Search.Services/api/v1/compositeaverages"
	RERATE        = "/USBC.Search.Services/api/v1/reratedaverage"
	LEAGUES       = "/USBC.Search.Services/api/v1/leagueactivities"
	MEMBERSHIPS   = "/USBC.Search.Services/api/v1/memberships"
	ACHIEVEMENTS  = "/USBC.Search.Services/api/v1/achievements"
	TOURNAMENTS   = "/USBC.Search.Services/api/v1/tournamentaverages"
	IGBO_AVERAGES = "/tournaments/get-igbo-tournament-tad-average/"
	IGBO_ID       = "/tournaments/igbots-id-lookup/"
)

type IGBO struct {
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
}

type Event struct {
	Date       string
	Tournament string
	Type       string
	Series     int
	Games      int
	Average    int
}

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

func queryForIGBO(url string, r *http.Request, id bool) *http.Response {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	if id {
		s := r.URL.Query()
		last := s["last"][0]

		q := req.URL.Query()
		q.Add("q", last)
		req.URL.RawQuery = q.Encode()
	} else {

		s := r.URL.Query()
		last := s["id"][0]
		years := s["yearrange"][0]

		q := req.URL.Query()
		q.Add("q", last)
		q.Add("yearrange", years)
		req.URL.RawQuery = q.Encode()
	}

	data, err := http.Get(req.URL.String())
	return data
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
	w.Write(body)
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

func IGBOIDHandler(w http.ResponseWriter, r *http.Request) {
	response := queryForIGBO(IGBO_SITE+IGBO_ID, r, true)
	tokens := html.NewTokenizer(response.Body)

	counter := 1
	igbo := IGBO{}
	newBowler := false
	igboBowlers := make([]IGBO, 0)

	s := r.URL.Query()
	first := s["first"][0]
	last := s["last"][0]
	usbc := s["usbc"][0]

	for {

		d := tokens.Next()

		if d == html.StartTagToken {

			data := tokens.Token().Data

			if data == "td" {

				tokens.Next()
				data = tokens.Token().Data

				if data == "a" {
					tokens.Next()
					data = tokens.Token().Data
				}

				if counter == 1 {
					igbo.ID = data
					newBowler = false
				}
				if counter == 2 {
					igbo.Name = strings.ToLower(data)
				}
				if counter == 3 {
					if data == "td" {
						igbo.City = ""
					} else {
						igbo.City = data
					}
				}
				if counter == 4 {
					if data == "td" {
						igbo.USBC = ""
					} else {
						igbo.USBC = data
					}

					igboBowlers = append(igboBowlers, igbo)
					counter = 1
					igbo = IGBO{}
					newBowler = true
				}

				if !newBowler {
					counter += 1
				}
			}
		}

		if d == html.ErrorToken {

			name := last + ", " + first
			for _, val := range igboBowlers {

				if name == val.Name {
					if val.USBC != "" {
						if val.USBC == usbc {
							w.Header().Set("Content-Type", "application/json")
							d, _ := json.Marshal(val.ID)
							w.Write(d)
							return
						}
					} else {
						w.Header().Set("Content-Type", "application/json")
						d, _ := json.Marshal(val.ID)
						w.Write(d)
						return
					}
				}

				if val.USBC == usbc {
					w.Header().Set("Content-Type", "application/json")
					d, _ := json.Marshal(val.ID)
					w.Write(d)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				d, _ := json.Marshal("bowler cannot be found")
				w.Write(d)
				return
			}
		}
	}
}

func IGBOHandler(w http.ResponseWriter, r *http.Request) {
	response := queryForIGBO(IGBO_SITE+IGBO_AVERAGES, r, false)
	tokens := html.NewTokenizer(response.Body)

	place := 0
	table := false
	table_count := 1
	table_2 := false
	count := 0
	bowler := Bowler{}

	for {

		d := tokens.Next()
		if d == html.StartTagToken {

			data := tokens.Token().Data

			if data == "table" {
				if table_count == 2 {
					table_2 = true
				}
				if table_count == 3 {
					bowler.Events = checkTable3(tokens)
					break
				}

				table = true
			} else if table {

				if data == "th" {
					tokens.Next()
				} else if data == "td" {
					if table_2 {
						tokens.Next()
						data = tokens.Token().Data
						if count < 3 {
							count += 1
						} else {
							if data == "strong" {
								tokens.Next()
								data = tokens.Token().Data
							}
							if place == 0 {
								bowler.TotalPins = data
							} else if place == 1 {
								bowler.TotalGames = data
							} else if place == 2 {
								bowler.Average = data
							}
							place += 1
							count += 1
						}
					} else {

						tokens.Next()
						data = tokens.Token().Data

						if place == 0 {
							bowler.IGBO = data
						} else if place == 1 {
							bowler.Name = data
						} else if place == 2 {
							bowler.City = data
						} else if place == 3 {
							bowler.USBC = data
						}

						place += 1
					}
				}
			}
		}

		if d == html.EndTagToken {

			if tokens.Token().Data == "table" {
				table = false
				table_count += 1
				place = 0
				table_2 = false

			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	d, _ := json.Marshal(bowler)
	w.Write(d)
}

func checkTable3(r *html.Tokenizer) []Event {

	counter := 0
	lis := make([]string, 0)
	place := 0
	events := []Event{}
	event := Event{}
	newEvent := false

	for {
		d := r.Next()
		if d == html.StartTagToken {

			token := r.Token()

			if token.Data == "th" {

				ty := r.Next()

				if ty == html.TextToken {
					data := r.Token().Data
					arr := strings.Split(data, "/")
					if len(arr) > 1 {
						counter += 1
						lis = append(lis, strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1]))

					} else {
						lis = append(lis, data)
					}
				}
			}

			if token.Data == "td" {

				ty := r.Next()

				if ty == html.TextToken {
					data := r.Token().Data
					if place == 0 {
						event.Date = data
						newEvent = false
					} else if place == 3 {
						in, _ := strconv.Atoi(data)
						event.Series = in
					} else if place == 4 {
						in, _ := strconv.Atoi(data)
						event.Games = in
					} else if place == 5 {
						in, _ := strconv.Atoi(data)
						event.Average = in
						events = append(events, event)
						event = Event{}
						place = 0
						newEvent = true
					}
					if !newEvent {
						place += 1
					}
				}

				if r.Token().Data == "strong" {
					strong := r.Next()

					if strong == html.TextToken {
						data := r.Token().Data
						if place == 1 {
							event.Tournament = data
						}
						place += 1

						//</strong>
						r.Next()
						//<Bullet point>
						bull := r.Next()

						//Continue on to what event it is
						if bull == html.TextToken {
							r.Next()
							r.Next()
							data := r.Token().Data
							if place == 2 {
								event.Type = data
							}
							place += 1
						}
					}
				}
			}
		}

		if d == html.ErrorToken {
			return events
		}
	}
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
	r.HandleFunc("/igbo/id", IGBOIDHandler)
	r.HandleFunc("/igbo/averages", IGBOHandler)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
