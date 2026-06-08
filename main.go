package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	_ "embed"
)

//go:embed index.html
var htmlContent string

// Struktura do odczytania JSONa z API
type WeatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		WindSpeed   float64 `json:"windspeed"`
	} `json:"current_weather"`
}

func main() {
	// Flaga do obsługi HEALTHCHECK
	checkHealth := flag.Bool("health", false, "Uruchamia test zdrowia")
	flag.Parse()

	if *checkHealth {
		resp, err := http.Get("http://127.0.0.1:8080/health")
		if err != nil || resp.StatusCode != 200 {
			os.Exit(1)
		}
		os.Exit(0)
	}

	log.Printf("Data uruchomienia: %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Printf("Autor programu: Kacper Kot")
	log.Printf("Serwer nasłuchuje na porcie TCP: 8080")

	// Endpoint główny
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			lat := r.FormValue("lat")
			lon := r.FormValue("lon")
			city := r.FormValue("city")

			url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true", lat, lon)
			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, "Błąd pobierania pogody (brak internetu?)", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			var weather WeatherResponse
			if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
				http.Error(w, "Błąd dekodowania danych", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "<div style='font-family: Arial; padding: 20px;'>")
			fmt.Fprintf(w, "<h2>Aktualna pogoda dla: %s</h2>", city)
			fmt.Fprintf(w, "<p><b>Temperatura:</b> %.1f °C</p>", weather.CurrentWeather.Temperature)
			fmt.Fprintf(w, "<p><b>Prędkość wiatru:</b> %.1f km/h</p>", weather.CurrentWeather.WindSpeed)
			fmt.Fprintf(w, "<br><a href='/'>Sprawdź inne miasto</a>")
			fmt.Fprintf(w, "</div>")
			return
		}
		
		tmpl := template.Must(template.New("index").Parse(htmlContent))
		tmpl.Execute(w, nil)
	})

	// Endpoint do healthcheck
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Uruchomienie serwera
	log.Fatal(http.ListenAndServe(":8080", nil))
}