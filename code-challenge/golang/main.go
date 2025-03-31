package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	data, err := loadPollutionData()
	if err != nil {
		panic(err)
	}

	da, err := parsePollutionData(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("PM-10 averages:")
	fmt.Printf("\tMorning: %.2f μg/m³\n", da.PM10.Morning)
	fmt.Printf("\tAfternoon: %.2f μg/m³\n", da.PM10.Afternoon)
	fmt.Printf("\tNight: %.2f μg/m³\n", da.PM10.Night)

	fmt.Println("PM-2.5 averages:")
	fmt.Printf("\tMorning: %.2f μg/m³\n", da.PM25.Morning)
	fmt.Printf("\tAfternoon: %.2f μg/m³\n", da.PM25.Afternoon)
	fmt.Printf("\tNight: %.2f μg/m³\n", da.PM25.Night)
}

type DailyAverages struct {
	PM10 struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	PM25 struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
}

func parsePollutionData(raw []byte) (DailyAverages, error) {
	var data struct {
		Hourly struct {
			Time []string  `json:"time"`
			PM10 []float64 `json:"pm10"`
			PM25 []float64 `json:"pm2_5"`
		} `json:"hourly"`
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		return DailyAverages{}, err
	}

	var morningPM10, morningPM25 float64
	var afternoonPM10, afternoonPM25 float64
	var nightPM10, nightPM25 float64
	var morningCount, afternoonCount, nightCount int

	for idx, timeOfDay := range data.Hourly.Time {
		dt, err := time.Parse("2006-01-02T15:04", timeOfDay)
		if err != nil {
			return DailyAverages{}, err
		}

		// Changed to >= to include 6. Added counter to see how many entries we got between these times.
		// Created hour variable since it should be same for all the calls.
		hour := dt.Hour()
		if hour >= 6 && hour < 12 {
			morningPM10 += data.Hourly.PM10[idx]
			morningPM25 += data.Hourly.PM25[idx]
			morningCount++
			continue
		}

		// Changed to >= to include 12. Added counter to see how many entries we got between these times.
		if hour >= 12 && hour < 18 {
			afternoonPM10 += data.Hourly.PM10[idx]
			afternoonPM25 += data.Hourly.PM25[idx]
			afternoonCount++
			continue
		}

		// Removed conditions since it should be all other cases. And added a Counter
		nightPM10 += data.Hourly.PM10[idx]
		nightPM25 += data.Hourly.PM25[idx]
		nightCount++
	}

	// Added the counts to divide with the actual number of inputs instead, added a check to avoid division by zero.
	da := DailyAverages{}
	if morningCount > 0 {
		da.PM10.Morning = morningPM10 / float64(morningCount)
		da.PM25.Morning = morningPM25 / float64(morningCount)
	}
	if afternoonCount > 0 {
		da.PM10.Afternoon = afternoonPM10 / float64(afternoonCount)
		da.PM25.Afternoon = afternoonPM25 / float64(afternoonCount)
	}
	if nightCount > 0 {
		da.PM10.Night = nightPM10 / float64(nightCount)
		da.PM25.Night = nightPM25 / float64(nightCount)
	}

	return da, nil
}

func loadPollutionData() ([]byte, error) {
	// Create the API URL with yesterday's date
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	url := fmt.Sprintf("https://air-quality-api.open-meteo.com/v1/air-quality?latitude=52.5235&longitude=13.4115&hourly=pm10,pm2_5&start_date=%s&end_date=%s", yesterday, yesterday)

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}
