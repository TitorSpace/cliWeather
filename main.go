package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		Epochtime int     `json:"last_updated_epoch"`
		Temp      float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		WindKph  float64 `json:"wind_kph"`
		Humidity int     `json:"humidity"`
	}

	Forecast struct {
		Forecastday []struct {
			Day struct {
				Maxtemp float64 `json:"maxtemp_c"`
				Mintemp float64 `json:"mintemp_c"`
			} `json:"day"`
			Astro struct {
				Sunrise string `json:"sunrise"`
				Sunset  string `json:"sunset"`
			} `json:"astro"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	fmt.Println("Hello World")

	debug := flag.Bool("debug", false, "Imprimir respuesta bruta y estructuras para depuración")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	city_query := "Vigo"

	// if args := flag.Args(); len(args) > 0 {
	// 	city_query = args[0]
	// 	fmt.Printf("This is the cityquery: %s", city_query)
	// }

	baseURL := "https://api.weatherapi.com/v1/forecast.json"

	//Get the API key
	apiKey := os.Getenv("WEATHER_API_KEY")
	fmt.Printf("This is the api key: %s", apiKey)
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Missing WEATHER_API_KEY (export WEATHER_API_KEY=tu_api_key)")
		os.Exit(1)
	}

	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("key", apiKey)
	q.Set("q", city_query)
	q.Set("lang", "es")
	q.Set("days", "1")
	q.Set("aqi", "no")
	q.Set("alerts", "no")
	u.RawQuery = q.Encode()

	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout a nivel de cliente
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Timeout a nivel de request
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating the request:", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "mruiz")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing the request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Weather API is not available. HTTP %d. Response: %s\n", resp.StatusCode, string(b))
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error leyendo respuesta:", err)
		os.Exit(1)
	}

	if *debug {
		fmt.Println("----- RAW JSON -----")
		fmt.Println(string(body))
		fmt.Println("--------------------")
	}

	var w Weather
	err = json.Unmarshal(body, &w)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		os.Exit(1)
	}

	if *debug {
		fmt.Printf("Parsed struct: %+v\n\n", w)
	}

	city, country := w.Location.Name, w.Location.Country

	windSpeed := w.Current.WindKph
	humidity := w.Current.Humidity
	description := w.Current.Condition.Text

	if len(w.Forecast.Forecastday) == 0 {
		fmt.Fprintln(os.Stderr, "No hay datos de pronóstico (forecastday está vacío).")
		os.Exit(1)
	}

	maxtemp := w.Forecast.Forecastday[0].Day.Maxtemp
	avgtemp := w.Current.Temp
	mintemp := w.Forecast.Forecastday[0].Day.Mintemp

	sunrise := w.Forecast.Forecastday[0].Astro.Sunrise
	sunset := w.Forecast.Forecastday[0].Astro.Sunset

	epoch := int64(w.Current.Epochtime)
	dateStr := time.Unix(epoch, 0).Format("Mon 02 Jan 2006 15:04:05 MST")

	fmt.Println(w)

	fmt.Printf("Good day %s, %s!\nToday's date is: %s\n", city, country, dateStr)

	fmt.Printf(`Today we have day with %s, where:
  max_temp: %.0f°C
  average_temp: %.0f°C
  min_temp: %.0f°C
  wind_speed: %.0f km/h
  humidity: %d%%
  SUNSET: %s
  SUNRISE: %s
`,
		description, maxtemp, avgtemp, mintemp, windSpeed, humidity, sunset, sunrise,
	)

}
