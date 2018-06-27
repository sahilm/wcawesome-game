package main

import (
	"os"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"strconv"
	"bytes"
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

type RefNotification struct {
	Country string `json:"country"`
	Event   event  `json:"event"`
}

func main() {
	fifaID := os.Getenv("FIFA_ID")
	country := os.Getenv("COUNTRY")
	refURL := os.Getenv("REF_URL")
	resp, err := http.Get("http://worldcup.sfg.io/matches")
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
	var refNotification RefNotification

	for i := range games {
		e := events(country, games[i], fifaID)
		if e != nil {
			for j := range e {
				time.Sleep(time.Duration(interval(e, j, j-1)) * 200 * time.Millisecond)
				refNotification.Country = country
				refNotification.Event = e[j]
				body, _ := json.Marshal(refNotification)
				req, err := http.NewRequest("POST", refURL, bytes.NewBuffer(body))
				if err != nil {
					fmt.Printf("%s", err)
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				resp.Body.Close()
				fmt.Println(refNotification)
			}
			refNotification.Country = "GAME_OVER"
			body, _ := json.Marshal(refNotification)
			req, err := http.NewRequest("POST", refURL, bytes.NewBuffer(body))
			if err != nil {
				fmt.Printf("%s", err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
			return
		}
	}
}

func interval(e []event, startIndex int, previousIndex int) int {
	if previousIndex < 0 {
		return timeparser(e[startIndex].Time)
	}
	return timeparser(e[startIndex].Time) - timeparser(e[previousIndex].Time)

}

func events(country string, g game, fifaID string) ([]event) {
	if fifaID != g.FifaId {
		return nil
	}
	switch country {
	case g.HomeTeam.Country:
		return g.HomeTeamEvents
	case g.AwayTeam.Country:
		return g.AwayTeamEvents
	default:
		return nil
	}
}

func timeparser(t string) int {
	t = strings.Replace(t, "'", "", -1)
	cleandtime := strings.Split(t, "+")
	x := 0
	for i := range cleandtime {
		intVal, err := strconv.Atoi(cleandtime[i])
		if err != nil {
			log.Fatal(err)
		}
		x += intVal
	}
	return x

}
