package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mruiz/cliWeather/internal/api/weatherapi"
	"strings"
	"testing"
	"time"
)

// containsANSI detecta secuencias ANSI básicas
func containsANSI(s string) bool {
	return strings.Contains(s, "\x1b[")
}

// containsAnyEmoji detecta si hay al menos uno de estos emojis comunes que usamos
func containsAnyEmoji(s string) bool {
	emojis := []string{"☀️", "⛅️", "☁️", "🌧️", "⛈️", "❄️", "🌫️", "💨", "🌅", "🌇", "☔️"}
	for _, e := range emojis {
		if strings.Contains(s, e) {
			return true
		}
	}
	return false
}

func TestRender_NoColorNoEmoji(t *testing.T) {
	now := time.Now().Unix()

	var w weatherapi.Weather
	// Construimos el objeto con JSON para evitar incompatibilidades de structs anónimos.
	payload := fmt.Sprintf(`{
      "location": {"name":"Vigo","country":"Spain"},
      "current": {
        "last_updated_epoch": %d,
        "temp_c": 18,
        "condition": {"text":"Parcialmente nublado"},
        "wind_kph": 10,
        "humidity": 55
      },
      "forecast": {
        "forecastday": [{
          "day": {
            "maxtemp_c": 20,
            "mintemp_c": 15,
            "daily_chance_of_rain": 30,
            "daily_will_it_rain": 0
          },
          "astro": {"sunrise":"08:00 AM","sunset":"08:00 PM"},
          "hour": [
            {
              "time_epoch": %d,
              "temp_c": 17.5,
              "chance_of_rain": 10,
              "condition": {"text":"Despejado"}
            },
            {
              "time_epoch": %d,
              "temp_c": 19.0,
              "chance_of_rain": 20,
              "condition": {"text":"Soleado"}
            }
          ]
        }]
      }
    }`, now, now, now+3600)

	if err := json.Unmarshal([]byte(payload), &w); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	var buf bytes.Buffer
	opt := Options{Color: false, Emoji: false}

	RenderHeader(&w, &buf, opt)
	if err := RenderDay(&w, 0, 1, &buf, opt); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	// No debe haber secuencias ANSI ni emojis
	if containsANSI(out) {
		t.Fatalf("expected no ANSI sequences, got:\n%s", out)
	}
	if containsAnyEmoji(out) {
		t.Fatalf("expected no emoji, got:\n%s", out)
	}

	// No debe aparecer el error de formato de fmt
	if strings.Contains(out, "%!s(") {
		t.Fatalf("format error found (%%!s()):\n%s", out)
	}

	// Debe contener partes clave
	for _, want := range []string{
		"¡Buen día", // encabezado
		"Fecha:",
		"Hoy:",
		"max:",
		"min:",
		"viento:",
		"humedad:",
		"amanecer:",
		"atardecer:",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRender_ColorEmoji(t *testing.T) {
	now := time.Now().Unix()

	var w weatherapi.Weather
	// Usamos "Soleado" para forzar ☀️ en pickConditionEmoji
	payload := fmt.Sprintf(`{
      "location": {"name":"Madrid","country":"Spain"},
      "current": {
        "last_updated_epoch": %d,
        "temp_c": 25,
        "condition": {"text":"Soleado"},
        "wind_kph": 4,
        "humidity": 34
      },
      "forecast": {
        "forecastday": [{
          "day": {
            "maxtemp_c": 26,
            "mintemp_c": 12,
            "daily_chance_of_rain": 0,
            "daily_will_it_rain": 0
          },
          "astro": {"sunrise":"08:06 AM","sunset":"08:05 PM"},
          "hour": [
            {
              "time_epoch": %d,
              "temp_c": 16,
              "chance_of_rain": 0,
              "condition": {"text":"Despejado"}
            },
            {
              "time_epoch": %d,
              "temp_c": 19,
              "chance_of_rain": 5,
              "condition": {"text":"Soleado"}
            }
          ]
        }]
      }
    }`, now, now, now+3600)

	if err := json.Unmarshal([]byte(payload), &w); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	var buf bytes.Buffer
	opt := Options{Color: true, Emoji: true}

	RenderHeader(&w, &buf, opt)
	if err := RenderDay(&w, 0, 1, &buf, opt); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	// Debe haber ANSI y al menos un emoji (☀️ por "Soleado")
	if !containsANSI(out) {
		t.Fatalf("expected ANSI sequences (color) in output, got:\n%s", out)
	}
	if !containsAnyEmoji(out) {
		t.Fatalf("expected emoji in output, got:\n%s", out)
	}

	// No debe aparecer el error de formato de fmt
	if strings.Contains(out, "%!s(") {
		t.Fatalf("format error found (%%!s()):\n%s", out)
	}

	// Algunas comprobaciones de contenido
	for _, want := range []string{
		"Hoy:",
		"max:",
		"min:",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
