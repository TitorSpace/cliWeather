package render

import (
	"fmt"
	"io"
	"mruiz/cliWeather/internal/api/weatherapi"
	"strings"
	"time"
)

// Options controla cómo se muestra el texto
type Options struct {
	Color bool
	Emoji bool
}

// ======= Tema de colores ANSI =======

type theme struct {
	reset string
	bold  func(string) string
	dim   func(string) string

	header func(string) string
	label  func(string) string
	value  func(string) string
	hot    func(string) string
	cold   func(string) string
	ok     func(string) string
	warn   func(string) string
}

func makeTheme(enable bool) theme {
	if !enable {
		noop := func(s string) string { return s }
		return theme{
			reset: "",
			bold:  noop, dim: noop,
			header: noop, label: noop, value: noop,
			hot: noop, cold: noop, ok: noop, warn: noop,
		}
	}
	wrap := func(code string) func(string) string {
		return func(s string) string { return code + s + "\x1b[0m" }
	}
	return theme{
		reset: "\x1b[0m",
		bold:  wrap("\x1b[1m"),
		dim:   wrap("\x1b[90m"),

		header: wrap("\x1b[96m"), // bright cyan
		label:  wrap("\x1b[37m"), // gray
		value:  wrap("\x1b[97m"), // bright white

		hot:  wrap("\x1b[31m"), // red
		cold: wrap("\x1b[34m"), // blue
		ok:   wrap("\x1b[32m"), // green
		warn: wrap("\x1b[33m"), // yellow
	}
}

// ======= Emoji helpers =======

func em(enabled bool, s string) string {
	if !enabled || s == "" {
		return ""
	}
	return s + " "
}

func pickConditionEmoji(enabled bool, text string) string {
	if !enabled {
		return ""
	}
	lt := strings.ToLower(text)
	switch {
	case strings.Contains(lt, "soleado"), strings.Contains(lt, "despejado"), strings.Contains(lt, "sunny"), strings.Contains(lt, "clear"):
		return "☀️"
	case strings.Contains(lt, "parcial"), strings.Contains(lt, "partly"), strings.Contains(lt, "intervals"):
		return "⛅️"
	case strings.Contains(lt, "nublado"), strings.Contains(lt, "cloud"):
		return "☁️"
	case strings.Contains(lt, "lluvia"), strings.Contains(lt, "rain"), strings.Contains(lt, "chubasc"):
		return "🌧️"
	case strings.Contains(lt, "tormenta"), strings.Contains(lt, "thunder"):
		return "⛈️"
	case strings.Contains(lt, "nieve"), strings.Contains(lt, "snow"):
		return "❄️"
	case strings.Contains(lt, "niebla"), strings.Contains(lt, "fog"), strings.Contains(lt, "mist"):
		return "🌫️"
	case strings.Contains(lt, "viento"), strings.Contains(lt, "wind"):
		return "💨"
	default:
		return "🌤️"
	}
}

// ======= Helpers de formateo =======

// tempColor elige el color según temperatura (Celsius)
func tempColor(th theme, c float64) func(string) string {
	switch {
	case c >= 30:
		return th.hot
	case c <= 10:
		return th.cold
	default:
		return th.value
	}
}

func fmtTemp(th theme, c float64) string {
	color := tempColor(th, c)
	return color(fmt.Sprintf("%.0f°C", c))
}

// Si ChanceOfRain es int en tus tipos, cambia la firma a (th theme, p int) string
func fmtPercent(th theme, p float64) string {
	var apply func(string) string
	switch {
	case p >= 60:
		apply = th.warn
	case p >= 20:
		apply = th.value
	default:
		apply = th.dim
	}
	return apply(fmt.Sprintf("%.0f%%", p))
}

// ======= API =======

func RenderHeader(w *weatherapi.Weather, out io.Writer, opt Options) {
	th := makeTheme(opt.Color)
	loc := fmt.Sprintf("%s, %s", w.Location.Name, w.Location.Country)
	t := time.Unix(int64(w.Current.LastUpdatedEpoch), 0).Local().Format("Mon 02 Jan 2006 15:04:05 MST")

	_, _ = fmt.Fprintf(out, "%s%s %s\n",
		em(opt.Emoji, "📍")+th.header("¡Buen día! "),
		th.bold(loc),
		"",
	)
	_, _ = fmt.Fprintf(out, "%s%s %s\n",
		em(opt.Emoji, "📅"), th.label("Fecha:"), th.value(t),
	)
}

func RenderAll(w *weatherapi.Weather, out io.Writer, opt Options) error {
	total := len(w.Forecast.Forecastday)
	if total == 0 {
		_, _ = fmt.Fprintln(out, "No Forecast available.")
		return nil
	}
	for i := range w.Forecast.Forecastday {
		if err := RenderDay(w, i, total, out, opt); err != nil {
			return err
		}
	}
	return nil
}

func RenderDay(w *weatherapi.Weather, idx, total int, out io.Writer, opt Options) error {
	th := makeTheme(opt.Color)

	if idx < 0 || idx >= len(w.Forecast.Forecastday) {
		_, _ = fmt.Fprintf(out, "Invalid day index %d (available: 0..%d)\n", idx, len(w.Forecast.Forecastday)-1)
		return nil
	}
	fd := w.Forecast.Forecastday[idx]

	// Fecha del bloque
	headerTime := time.Now().Local()
	if len(fd.Hour) > 0 {
		headerTime = time.Unix(int64(fd.Hour[0].TimeEpoch), 0).Local()
	}
	dayTitle := fmt.Sprintf("%s (día %d/%d)", headerTime.Format("Mon 02 Jan 2006"), idx+1, total)
	_, _ = fmt.Fprintf(out, "\n%s %s\n", th.bold("==="), th.bold(th.header(dayTitle)))

	// Condición representativa
	conditionText := w.Current.Condition.Text
	if len(fd.Hour) > 0 {
		mid := len(fd.Hour) / 2
		conditionText = fd.Hour[mid].Condition.Text
	}

	// Media del día por horas
	avg := w.Current.TempC
	if len(fd.Hour) > 0 {
		var sum float64
		for _, h := range fd.Hour {
			sum += h.TempC
		}
		avg = sum / float64(len(fd.Hour))
	}

	// Iconos
	iconCond := pickConditionEmoji(opt.Emoji, conditionText)
	iconMax := em(opt.Emoji, "🔺")
	iconAvg := em(opt.Emoji, "📊")
	iconMin := em(opt.Emoji, "🔻")
	iconWind := em(opt.Emoji, "💨")
	iconHum := em(opt.Emoji, "💧")
	iconSunrise := em(opt.Emoji, "🌅")
	iconSunset := em(opt.Emoji, "🌇")

	// Resumen del día

	_, _ = fmt.Fprintf(out, "%s%s %s\n", iconCond+th.label("Hoy:"), "", th.value(conditionText))
	_, _ = fmt.Fprintf(out, "  %s%s  %s  %s%s  %s  %s%s  %s\n",
		iconMax, th.label("max:"), fmtTemp(th, fd.Day.MaxtempC),
		iconAvg, th.label("avg:"), fmtTemp(th, avg),
		iconMin, th.label("min:"), fmtTemp(th, fd.Day.MintempC),
	)
	_, _ = fmt.Fprintf(out, "  %s%s %s  %s%s %s\n",
		iconWind, th.label("viento:"), th.value(fmt.Sprintf("%.0f km/h", w.Current.WindKph)),
		iconHum, th.label("humedad:"), th.value(fmt.Sprintf("%d%%", w.Current.Humidity)),
	)
	_, _ = fmt.Fprintf(out, "  %s%s %s  %s%s %s\n\n",
		iconSunrise, th.label("amanecer:"), th.value(fd.Astro.Sunrise),
		iconSunset, th.label("atardecer:"), th.value(fd.Astro.Sunset),
	)

	// Horas
	umbrella := em(opt.Emoji, "☔️")
	for _, hour := range fd.Hour {
		tm := time.Unix(int64(hour.TimeEpoch), 0).Local()
		hhmm := tm.Format("15:04")
		condEm := pickConditionEmoji(opt.Emoji, hour.Condition.Text)

		// Construimos cada parte ya coloreada
		timePart := th.dim(hhmm)
		condPart := condEm + th.value(hour.Condition.Text)
		tempPart := fmtTemp(th, hour.TempC)
		rainPart := umbrella + fmtPercent(th, hour.ChanceOfRain)

		_, _ = fmt.Fprintf(out, "%s %s - %s, %s\n", timePart, condPart, tempPart, rainPart)
	}

	return nil
}
