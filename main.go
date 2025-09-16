package main

import (
	"context"
	"flag"
	"fmt"
	"mruiz/cliWeather/internal/api/weatherapi"
	"mruiz/cliWeather/internal/config"
	"os"
	"time"
)

func main() {

	cfg := config.FromEnv()

	city := flag.String("city", "Vigo", "City name or query")
	days := flag.Int("days", cfg.Days, "Foracast days (1-3 on free tier)")
	lang := flag.String("lang", cfg.Language, "Choose your pref language, f.e. es, en, fr...")
	apikey := flag.String("apikey", cfg.APIKey, "The WeatherAPI key (you can choose to set WEATHER_API_KEY)")
	debug := flag.Bool("debug", false, "Print raw structs for debugging")
	flag.Parse()

	if *apikey == "" {
		fmt.Fprintln(os.Stderr, "Missing WEATHER_API_KEY (export WEATHER_API_KEY=tu_api_key)")
		os.Exit(1)
	}

	client := weatherapi.NewClient(*apikey, *lang, cfg.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)

	defer cancel()

	w, err := client.Forecast(ctx, *city, *days, false, false)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if *debug {
		fmt.Printf("%+v\n\n", *w)
	}

	if len(w.Forecast.Forecastday) == 0 {
		fmt.Printf("No Forecast available.")
		return
	}

	fd := w.Forecast.Forecastday[0]
	t := time.Unix(int64(w.Current.LastUpdatedEpoch), 0).Local().Format("Mon 02 Jan 2006 15:04:05 MST")

	fmt.Printf("¡Buen día %s, %s!\nFecha: %s\n", w.Location.Name, w.Location.Country, t)
	fmt.Printf(
		"Hoy: %s\n  max: %.0f°C  avg: %.0f°C  min: %.0f°C\n  viento: %.0f km/h  humedad: %d%%\n  amanecer: %s  atardecer: %s\n \n",
		w.Current.Condition.Text,
		fd.Day.MaxtempC, w.Current.TempC, fd.Day.MintempC,
		w.Current.WindKph, w.Current.Humidity,
		fd.Astro.Sunrise, fd.Astro.Sunset,
	)

	for _, hour := range fd.Hour {
		date := time.Unix(int64(hour.TimeEpoch), 0)
		fmt.Printf("%s - %.0f°C, %.0f°C, %s\n",
			date.Format("15:03"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
	}

}
