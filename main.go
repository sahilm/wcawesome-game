package main

import (
	"os"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
)

type team struct {
	Country string `json:"country"`
	Code    string `json:"code"`
	Goals   int    `json:"goals"`
}

type event struct {
	Id          int64  `json:"id"`
	Player      string `json:"player"`
	TypeOfEvent string `json:"type_of_event"`
	Time        string `json:"time"`
}

type game struct {
	Venue             string  `json:"venue"`
	Location          string  `json:"location"`
	Status            string  `json:"status"`
	Time              string  `json:"time"`
	FifaId            string  `json:"fifa_id"`
	Datetime          string  `json:"datetime"`
	LastEventUpdateAt string  `json:"last_event_update_at"`
	LastScoreUpdateAt string  `json:"last_score_update_at"`
	HomeTeam          team    `json:"home_team"`
	AwayTeam          team    `json:"away_team"`
	Country           string  `json:"teams"`
	Winner            string  `json:"winner"`
	WinnerCode        string  `json:"winner_code"`
	HomeTeamEvents    []event `json:"home_team_events"`
	AwayTeamEvents    []event `json:"away_team_events"`
}

func main() {
	fifaID := os.Getenv("FIFA_ID")
	country := os.Getenv("COUNTRY")
	url := "http://worldcup.sfg.io/matches"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var games []game
	err = json.Unmarshal(bodyBytes, &games)
	if err != nil {
		log.Fatal(err)
	}

	for g := range games {
		ok, homeOrAway := isHomeOrAway(country, games[g].HomeTeam.Country, games[g].AwayTeam.Country)
		if games[g].FifaId == fifaID && ok {
			fmt.Printf(games[g].Status + " " + homeOrAway)
			if homeOrAway == "home" {
				for i := range games[g].HomeTeamEvents {
					println(games[g].HomeTeamEvents[i].TypeOfEvent)
					time.Sleep(1 * time.Second)
				}

			} else {
				for i := range games[g].AwayTeamEvents {
					println(games[g].AwayTeamEvents[i].TypeOfEvent)
					time.Sleep(1 * time.Second)
				}
			}
		}
	}
}

func isHomeOrAway(country string, home string, away string) (bool, string) {
	switch country {
	case home:
		return true, "home"
	case away:
		return true, "away"
	default:
		return false, ""
	}
}
