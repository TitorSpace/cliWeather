package cliweather

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cliWeather/internal/api/weatherapi"
	"cliWeather/internal/config"
	"cliWeather/internal/render"

	"github.com/spf13/cobra"
)

var (
	flagCity     string
	flagDays     int
	flagLang     string
	flagAPIKey   string
	flagDebug    bool
	flagDayIndex int
	flagJSON     bool
	flagUnits    string
	flagHours    string
)

var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Displays the weather forecast",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.FromEnv()

		if flagLang == "" {
			flagLang = cfg.Language
		}
		if flagDays == 0 {
			flagDays = cfg.Days
		}
		if flagAPIKey == "" {
			flagAPIKey = cfg.APIKey
		}
		if flagAPIKey == "" {
			return fmt.Errorf("missing WEATHER_API_KEY (export WEATHER_API_KEY=tu_api_key o usa --apikey)")
		}

		client := weatherapi.NewClient(flagAPIKey, flagLang, cfg.Timeout)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		w, err := client.Forecast(ctx, flagCity, flagDays, false, false)
		if err != nil {
			return err
		}

		if flagDebug {
			fmt.Printf("%+v\n\n", *w)
		}

		if len(w.Forecast.Forecastday) == 0 {
			fmt.Fprintln(os.Stdout, "No Forecast available.")
			return nil
		}

		// Raw JSON output if needed
		if flagJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(w)
		}

		useColor := !noColor && !envNoColor() && isTerminal(os.Stdout)
		useEmoji := !noEmoji // (podrías condicionar por OS o TTY si quisieras)
		opt := render.Options{
			Color: useColor,
			Emoji: useEmoji,
			Units: render.ParseUnits(flagUnits),
		}
		hr, err := render.ParseHourRange(flagHours)
		if err != nil {
			return fmt.Errorf("invalid --hours: %w", err)
		}
		opt.HourRange = hr

		// General header + rendered days
		render.RenderHeader(w, os.Stdout, opt)
		if flagDayIndex >= 0 {
			return render.RenderDay(w, flagDayIndex, len(w.Forecast.Forecastday), os.Stdout, opt)
		}
		return render.RenderAll(w, os.Stdout, opt)
	},
}

func init() {
	rootCmd.AddCommand(forecastCmd)

	forecastCmd.Flags().StringVarP(&flagCity, "city", "c", "Vigo", "City name or query")
	forecastCmd.Flags().IntVarP(&flagDays, "days", "d", 1, "Forecast days (1-3 on free tier)")
	forecastCmd.Flags().StringVarP(&flagLang, "lang", "l", "en", "Language (e.g., es, en, fr)")
	forecastCmd.Flags().StringVar(&flagAPIKey, "apikey", "", "WeatherAPI key (or set WEATHER_API_KEY)")
	forecastCmd.Flags().BoolVar(&flagDebug, "debug", false, "Print raw structs for debugging")
	forecastCmd.Flags().IntVar(&flagDayIndex, "day-index", -1, "Show only this forecast day index (0..days-1)")
	forecastCmd.Flags().BoolVar(&flagJSON, "json", false, "Print raw JSON response")
	forecastCmd.Flags().StringVarP(&flagUnits, "units", "u", "", "Units: metric|imperial")
	forecastCmd.Flags().StringVar(&flagHours, "hours", "", "Filter hours range, e.g. 08-18 (local time)")
}
