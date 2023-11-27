package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/go-resty/resty/v2"
	"github.com/gobuffalo/packr/v2"
)

// open -gj -a /Applications/habiticaNotification
var path string

func main() {
	var (
		anyReminder bool
		wg          sync.WaitGroup
	)

	box := packr.New("Box", "./")
	path = box.Path

	// Create a Resty client
	client := resty.New()

	// Set the Habitica API base URL
	habiticaURL := "https://habitica.com/api/v3/"

	// Set your Habitica user ID and API token
	userID := "1cc6f0db-02eb-4d0c-b1ee-a68dba756599"
	apiToken := "b846a3ac-aaf1-4804-bbdd-65e15aba43bf"

	// Set the API endpoint you want to access
	endpoint := "tasks/user"

	// Make the GET request
	resp, err := client.R().
		SetHeader("x-api-user", userID).
		SetHeader("x-api-key", apiToken).
		SetQueryParam("type", "dailys").
		Get(habiticaURL + endpoint)
	if err != nil {
		panic(err)
	}

	var habiticaData HabiticaData

	// Unmarshal the response body into a HabiticaData struct
	err = json.Unmarshal(resp.Body(), &habiticaData)
	if err != nil {
		panic(err)
	}

	for _, habit := range habiticaData.Data {
		if !habit.IsDue {
			continue
		}

		log.Printf("habit text: %+v isDue: %v\n", habit.Text, habit.IsDue)
		log.Println("----------------------------")

		if len(habit.Reminders) == 0 {
			log.Printf("no reminders for habit: %+v\n", habit.Text)
		}

		for _, reminder := range habit.Reminders {
			// Get today's date
			today := time.Now().UTC()
			reminderTime := reminder.Time

			// Combine with today's date
			combinedTime := time.Date(
				today.Year(),
				today.Month(),
				today.Day(),
				reminderTime.Hour(),
				reminderTime.Minute(),
				reminderTime.Second(),
				0,
				time.UTC)

			tickerDuration := time.Until(combinedTime)
			log.Printf("reminder time: %v duration: %+v\n", reminder.Time.String(), tickerDuration)
			log.Printf("combined time: %v\n", combinedTime.String())
			log.Println("********************************************")

			if tickerDuration < 0 {
				continue
			}

			anyReminder = true

			go scheduleHabitNotification(tickerDuration, habit.Text, &wg)

			wg.Add(1)
		}
	}

	wg.Wait()

	if !anyReminder {
		os.Exit(1)
	}

	log.Println("all reminders are set. ")
}

func scheduleHabitNotification(tickerDuration time.Duration, text string, wg *sync.WaitGroup) {
	defer wg.Done()

	if tickerDuration < 0 {
		log.Println("returning")

		return
	}

	ticker := time.NewTicker(tickerDuration)
	<-ticker.C

	defer ticker.Stop()

	notification := gosxnotifier.NewNotification("Habitica Reminder!")
	notification.Title = text
	notification.Subtitle = text
	notification.Sound = gosxnotifier.Hero
	notification.ContentImage = path + "/habitica.png"
	notification.AppIcon = path + "/habitica.png"

	if err := notification.Push(); err != nil {
		panic(err)
	}
}
