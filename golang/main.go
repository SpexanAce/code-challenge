package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// Set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve the main page
	http.HandleFunc("/", handleHome)

	// API endpoint for pollution data
	http.HandleFunc("/api/pollution", handlePollutionData)

	// Debug endpoint to see raw API response
	http.HandleFunc("/api/debug", handleDebug)

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handlePollutionData(w http.ResponseWriter, r *http.Request) {
	// Get latitude, longitude, data type, and date from query parameters
	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")
	dataType := r.URL.Query().Get("type")
	date := r.URL.Query().Get("date")

	// If not provided, use default values
	if lat == "" {
		lat = "57.7047"
	}
	if lng == "" {
		lng = "11.9684"
	}
	if dataType == "" {
		dataType = "particulate"
	}
	if date == "" {
		date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	data, err := loadPollutionData(lat, lng, dataType, date)
	if err != nil {
		log.Printf("Error loading pollution data: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}

	da, err := parsePollutionData(data, dataType)
	if err != nil {
		log.Printf("Error parsing pollution data: %v", err)
		http.Error(w, fmt.Sprintf("Failed to parse data: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert the data to JSON format based on type
	var response interface{}
	if dataType == "particulate" {
		response = struct {
			PM10 struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"pm10"`
			PM25 struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"pm25"`
		}{
			PM10: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.PM10.Morning,
				Afternoon: da.PM10.Afternoon,
				Night:     da.PM10.Night,
			},
			PM25: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.PM25.Morning,
				Afternoon: da.PM25.Afternoon,
				Night:     da.PM25.Night,
			},
		}
	} else if dataType == "pollen" {
		response = struct {
			ALDER struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"alder"`
			BIRCH struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"birch"`
			GRASS struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"grass"`
			MUGWORT struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"mugwort"`
			OLIVE struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"olive"`
			RAGWEED struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"ragweed"`
		}{
			ALDER: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.ALDER.Morning,
				Afternoon: da.ALDER.Afternoon,
				Night:     da.ALDER.Night,
			},
			BIRCH: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.BIRCH.Morning,
				Afternoon: da.BIRCH.Afternoon,
				Night:     da.BIRCH.Night,
			},
			GRASS: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.GRASS.Morning,
				Afternoon: da.GRASS.Afternoon,
				Night:     da.GRASS.Night,
			},
			MUGWORT: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.MUGWORT.Morning,
				Afternoon: da.MUGWORT.Afternoon,
				Night:     da.MUGWORT.Night,
			},
			OLIVE: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.OLIVE.Morning,
				Afternoon: da.OLIVE.Afternoon,
				Night:     da.OLIVE.Night,
			},
			RAGWEED: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.RAGWEED.Morning,
				Afternoon: da.RAGWEED.Afternoon,
				Night:     da.RAGWEED.Night,
			},
		}
	} else {
		response = struct {
			CO struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"co"`
			NO2 struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"no2"`
			SO2 struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"so2"`
			O3 struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			} `json:"o3"`
		}{
			CO: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.CO.Morning,
				Afternoon: da.CO.Afternoon,
				Night:     da.CO.Night,
			},
			NO2: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.NO2.Morning,
				Afternoon: da.NO2.Afternoon,
				Night:     da.NO2.Night,
			},
			SO2: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.SO2.Morning,
				Afternoon: da.SO2.Afternoon,
				Night:     da.SO2.Night,
			},
			O3: struct {
				Morning   float64 `json:"morning"`
				Afternoon float64 `json:"afternoon"`
				Night     float64 `json:"night"`
			}{
				Morning:   da.O3.Morning,
				Afternoon: da.O3.Afternoon,
				Night:     da.O3.Night,
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func handleDebug(w http.ResponseWriter, r *http.Request) {
	data, err := loadPollutionData("57.7047", "11.9684", "particulate", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
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
	ALDER struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	BIRCH struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	GRASS struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	MUGWORT struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	OLIVE struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	RAGWEED struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	CO struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	NO2 struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	SO2 struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
	O3 struct {
		Morning   float64
		Afternoon float64
		Night     float64
	}
}

func parsePollutionData(raw []byte, dataType string) (DailyAverages, error) {
	log.Printf("Starting to parse pollution data")

	var data struct {
		Hourly struct {
			Time    []string  `json:"time"`
			PM10    []float64 `json:"pm10,omitempty"`
			PM25    []float64 `json:"pm2_5,omitempty"`
			ALDER   []float64 `json:"alder_pollen,omitempty"`
			BIRCH   []float64 `json:"birch_pollen,omitempty"`
			GRASS   []float64 `json:"grass_pollen,omitempty"`
			MUGWORT []float64 `json:"mugwort_pollen,omitempty"`
			OLIVE   []float64 `json:"olive_pollen,omitempty"`
			RAGWEED []float64 `json:"ragweed_pollen,omitempty"`
			CO      []float64 `json:"carbon_monoxide,omitempty"`
			NO2     []float64 `json:"nitrogen_dioxide,omitempty"`
			SO2     []float64 `json:"sulphur_dioxide,omitempty"`
			O3      []float64 `json:"ozone,omitempty"`
		} `json:"hourly"`
		HourlyUnits struct {
			Time    string `json:"time"`
			PM10    string `json:"pm10,omitempty"`
			PM25    string `json:"pm2_5,omitempty"`
			ALDER   string `json:"alder_pollen,omitempty"`
			BIRCH   string `json:"birch_pollen,omitempty"`
			GRASS   string `json:"grass_pollen,omitempty"`
			MUGWORT string `json:"mugwort_pollen,omitempty"`
			OLIVE   string `json:"olive_pollen,omitempty"`
			RAGWEED string `json:"ragweed_pollen,omitempty"`
			CO      string `json:"carbon_monoxide,omitempty"`
			NO2     string `json:"nitrogen_dioxide,omitempty"`
			SO2     string `json:"sulphur_dioxide,omitempty"`
			O3      string `json:"ozone,omitempty"`
		} `json:"hourly_units"`
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		return DailyAverages{}, err
	}

	log.Printf("Successfully parsed JSON. Found %d time entries", len(data.Hourly.Time))
	log.Printf("Units - CO: %s, NO2: %s, SO2: %s, O3: %s",
		data.HourlyUnits.CO, data.HourlyUnits.NO2, data.HourlyUnits.SO2, data.HourlyUnits.O3)

	var morningPM10, morningPM25 float64
	var afternoonPM10, afternoonPM25 float64
	var nightPM10, nightPM25 float64
	var morningALDER, morningBIRCH, morningGRASS, morningMUGWORT, morningOLIVE, morningRAGWEED float64
	var afternoonALDER, afternoonBIRCH, afternoonGRASS, afternoonMUGWORT, afternoonOLIVE, afternoonRAGWEED float64
	var nightALDER, nightBIRCH, nightGRASS, nightMUGWORT, nightOLIVE, nightRAGWEED float64
	var morningCO, morningNO2, morningSO2, morningO3 float64
	var afternoonCO, afternoonNO2, afternoonSO2, afternoonO3 float64
	var nightCO, nightNO2, nightSO2, nightO3 float64
	var morningCount, afternoonCount, nightCount int

	for idx, timeOfDay := range data.Hourly.Time {
		dt, err := time.Parse("2006-01-02T15:04", timeOfDay)
		if err != nil {
			log.Printf("Failed to parse time '%s': %v", timeOfDay, err)
			return DailyAverages{}, err
		}

		hour := dt.Hour()
		if hour >= 6 && hour < 12 {
			if dataType == "particulate" {
				morningPM10 += data.Hourly.PM10[idx]
				morningPM25 += data.Hourly.PM25[idx]
			} else if dataType == "pollen" {
				morningALDER += data.Hourly.ALDER[idx]
				morningBIRCH += data.Hourly.BIRCH[idx]
				morningGRASS += data.Hourly.GRASS[idx]
				morningMUGWORT += data.Hourly.MUGWORT[idx]
				morningOLIVE += data.Hourly.OLIVE[idx]
				morningRAGWEED += data.Hourly.RAGWEED[idx]
			} else {
				// Convert CO from μg/m³ to mg/m³ if needed
				if data.HourlyUnits.CO == "μg/m³" {
					morningCO += data.Hourly.CO[idx] / 1000
				} else {
					morningCO += data.Hourly.CO[idx]
				}
				morningNO2 += data.Hourly.NO2[idx]
				morningSO2 += data.Hourly.SO2[idx]
				morningO3 += data.Hourly.O3[idx]
			}
			morningCount++
			continue
		}

		if hour >= 12 && hour < 18 {
			if dataType == "particulate" {
				afternoonPM10 += data.Hourly.PM10[idx]
				afternoonPM25 += data.Hourly.PM25[idx]
			} else if dataType == "pollen" {
				afternoonALDER += data.Hourly.ALDER[idx]
				afternoonBIRCH += data.Hourly.BIRCH[idx]
				afternoonGRASS += data.Hourly.GRASS[idx]
				afternoonMUGWORT += data.Hourly.MUGWORT[idx]
				afternoonOLIVE += data.Hourly.OLIVE[idx]
				afternoonRAGWEED += data.Hourly.RAGWEED[idx]
			} else {
				// Convert CO from μg/m³ to mg/m³ if needed
				if data.HourlyUnits.CO == "μg/m³" {
					afternoonCO += data.Hourly.CO[idx] / 1000
				} else {
					afternoonCO += data.Hourly.CO[idx]
				}
				afternoonNO2 += data.Hourly.NO2[idx]
				afternoonSO2 += data.Hourly.SO2[idx]
				afternoonO3 += data.Hourly.O3[idx]
			}
			afternoonCount++
			continue
		}

		if dataType == "particulate" {
			nightPM10 += data.Hourly.PM10[idx]
			nightPM25 += data.Hourly.PM25[idx]
		} else if dataType == "pollen" {
			nightALDER += data.Hourly.ALDER[idx]
			nightBIRCH += data.Hourly.BIRCH[idx]
			nightGRASS += data.Hourly.GRASS[idx]
			nightMUGWORT += data.Hourly.MUGWORT[idx]
			nightOLIVE += data.Hourly.OLIVE[idx]
			nightRAGWEED += data.Hourly.RAGWEED[idx]
		} else {
			// Convert CO from μg/m³ to mg/m³ if needed
			if data.HourlyUnits.CO == "μg/m³" {
				nightCO += data.Hourly.CO[idx] / 1000
			} else {
				nightCO += data.Hourly.CO[idx]
			}
			nightNO2 += data.Hourly.NO2[idx]
			nightSO2 += data.Hourly.SO2[idx]
			nightO3 += data.Hourly.O3[idx]
		}
		nightCount++
	}

	log.Printf("Counted entries - Morning: %d, Afternoon: %d, Night: %d", morningCount, afternoonCount, nightCount)

	da := DailyAverages{}
	if morningCount > 0 {
		if dataType == "particulate" {
			da.PM10.Morning = morningPM10 / float64(morningCount)
			da.PM25.Morning = morningPM25 / float64(morningCount)
		} else if dataType == "pollen" {
			da.ALDER.Morning = morningALDER / float64(morningCount)
			da.BIRCH.Morning = morningBIRCH / float64(morningCount)
			da.GRASS.Morning = morningGRASS / float64(morningCount)
			da.MUGWORT.Morning = morningMUGWORT / float64(morningCount)
			da.OLIVE.Morning = morningOLIVE / float64(morningCount)
			da.RAGWEED.Morning = morningRAGWEED / float64(morningCount)
		} else {
			da.CO.Morning = morningCO / float64(morningCount)
			da.NO2.Morning = morningNO2 / float64(morningCount)
			da.SO2.Morning = morningSO2 / float64(morningCount)
			da.O3.Morning = morningO3 / float64(morningCount)
		}
	}
	if afternoonCount > 0 {
		if dataType == "particulate" {
			da.PM10.Afternoon = afternoonPM10 / float64(afternoonCount)
			da.PM25.Afternoon = afternoonPM25 / float64(afternoonCount)
		} else if dataType == "pollen" {
			da.ALDER.Afternoon = afternoonALDER / float64(afternoonCount)
			da.BIRCH.Afternoon = afternoonBIRCH / float64(afternoonCount)
			da.GRASS.Afternoon = afternoonGRASS / float64(afternoonCount)
			da.MUGWORT.Afternoon = afternoonMUGWORT / float64(afternoonCount)
			da.OLIVE.Afternoon = afternoonOLIVE / float64(afternoonCount)
			da.RAGWEED.Afternoon = afternoonRAGWEED / float64(afternoonCount)
		} else {
			da.CO.Afternoon = afternoonCO / float64(afternoonCount)
			da.NO2.Afternoon = afternoonNO2 / float64(afternoonCount)
			da.SO2.Afternoon = afternoonSO2 / float64(afternoonCount)
			da.O3.Afternoon = afternoonO3 / float64(afternoonCount)
		}
	}
	if nightCount > 0 {
		if dataType == "particulate" {
			da.PM10.Night = nightPM10 / float64(nightCount)
			da.PM25.Night = nightPM25 / float64(nightCount)
		} else if dataType == "pollen" {
			da.ALDER.Night = nightALDER / float64(nightCount)
			da.BIRCH.Night = nightBIRCH / float64(nightCount)
			da.GRASS.Night = nightGRASS / float64(nightCount)
			da.MUGWORT.Night = nightMUGWORT / float64(nightCount)
			da.OLIVE.Night = nightOLIVE / float64(nightCount)
			da.RAGWEED.Night = nightRAGWEED / float64(nightCount)
		} else {
			da.CO.Night = nightCO / float64(nightCount)
			da.NO2.Night = nightNO2 / float64(nightCount)
			da.SO2.Night = nightSO2 / float64(nightCount)
			da.O3.Night = nightO3 / float64(nightCount)
		}
	}

	if dataType == "particulate" {
		log.Printf("Calculated averages - PM10: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.PM10.Morning, da.PM10.Afternoon, da.PM10.Night)
		log.Printf("Calculated averages - PM2.5: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.PM25.Morning, da.PM25.Afternoon, da.PM25.Night)
	} else if dataType == "pollen" {
		log.Printf("Calculated averages - ALDER: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.ALDER.Morning, da.ALDER.Afternoon, da.ALDER.Night)
		log.Printf("Calculated averages - BIRCH: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.BIRCH.Morning, da.BIRCH.Afternoon, da.BIRCH.Night)
		log.Printf("Calculated averages - GRASS: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.GRASS.Morning, da.GRASS.Afternoon, da.GRASS.Night)
		log.Printf("Calculated averages - MUGWORT: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.MUGWORT.Morning, da.MUGWORT.Afternoon, da.MUGWORT.Night)
		log.Printf("Calculated averages - OLIVE: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.OLIVE.Morning, da.OLIVE.Afternoon, da.OLIVE.Night)
		log.Printf("Calculated averages - RAGWEED: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.RAGWEED.Morning, da.RAGWEED.Afternoon, da.RAGWEED.Night)
	} else {
		log.Printf("Calculated averages - CO: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.CO.Morning, da.CO.Afternoon, da.CO.Night)
		log.Printf("Calculated averages - NO2: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.NO2.Morning, da.NO2.Afternoon, da.NO2.Night)
		log.Printf("Calculated averages - SO2: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.SO2.Morning, da.SO2.Afternoon, da.SO2.Night)
		log.Printf("Calculated averages - O3: Morning=%.2f, Afternoon=%.2f, Night=%.2f",
			da.O3.Morning, da.O3.Afternoon, da.O3.Night)
	}

	return da, nil
}

func loadPollutionData(lat, lng, dataType, date string) ([]byte, error) {
	var hourlyParams string
	if dataType == "particulate" {
		hourlyParams = "pm10,pm2_5"
	} else if dataType == "pollen" {
		hourlyParams = "alder_pollen,birch_pollen,grass_pollen,mugwort_pollen,olive_pollen,ragweed_pollen"
	} else {
		hourlyParams = "carbon_monoxide,nitrogen_dioxide,sulphur_dioxide,ozone"
	}

	url := fmt.Sprintf("https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%s&longitude=%s&hourly=%s&start_date=%s&end_date=%s",
		lat, lng, hourlyParams, date, date)

	log.Printf("Fetching data from URL: %s", url)

	// Create a client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the HTTP request
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		// Try alternative API endpoint if the first one fails
		altUrl := fmt.Sprintf("https://api.open-meteo.com/v1/air-quality?latitude=%s&longitude=%s&hourly=%s&start_date=%s&end_date=%s",
			lat, lng, hourlyParams, date, date)
		log.Printf("Trying alternative URL: %s", altUrl)

		resp, err = client.Get(altUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data from both endpoints: %w", err)
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-200 status code: %s", resp.Status)
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Successfully received response with length: %d bytes", len(body))
	return body, nil
}
