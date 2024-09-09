package weather

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetHourlyForecast(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedQuery := "latitude=52.520000&longitude=13.405000&hourly=precipitation_probability,temperature_2m,cloud_cover&forecast_days=16&timezone=Europe/Berlin"
		if r.URL.RawQuery != expectedQuery {
			t.Errorf("Unexpected query parameters: %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"hourly": {
				"time": ["2023-04-14T00:00", "2023-04-14T01:00"],
				"precipitation_probability": [30, 60],
				"temperature_2m": [15.5, 16.2],
				"cloud_cover": [25, 30]
			}
		}`))
	}))
	defer server.Close()

	service := &OpenMeteoService{
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	forecast, err := service.GetHourlyForecast(52.52, 13.405, "Europe/Berlin")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}

	expectedTime1 := time.Date(2023, 4, 14, 0, 0, 0, 0, loc)
	expectedTime2 := time.Date(2023, 4, 14, 1, 0, 0, 0, loc)

	expected := HourlyForecast{
		Time:                     []time.Time{expectedTime1, expectedTime2},
		PrecipitationProbability: []float64{30, 60},
		Temperature:              []float64{15.5, 16.2},
		CloudCover:               []float64{25, 30},
	}

	if len(forecast.Time) != len(expected.Time) {
		t.Fatalf("Expected %d time entries, got %d", len(expected.Time), len(forecast.Time))
	}

	for i := range forecast.Time {
		if !forecast.Time[i].Equal(expected.Time[i]) {
			t.Errorf("Expected time %v, got %v", expected.Time[i], forecast.Time[i])
		}
		if forecast.PrecipitationProbability[i] != expected.PrecipitationProbability[i] {
			t.Errorf("Expected precipitation probability %f, got %f", expected.PrecipitationProbability[i], forecast.PrecipitationProbability[i])
		}
		if forecast.Temperature[i] != expected.Temperature[i] {
			t.Errorf("Expected temperature %f, got %f", expected.Temperature[i], forecast.Temperature[i])
		}
		if forecast.CloudCover[i] != expected.CloudCover[i] {
			t.Errorf("Expected cloud cover %f, got %f", expected.CloudCover[i], forecast.CloudCover[i])
		}
	}
}

func TestNewOpenMeteoService(t *testing.T) {
	service := NewOpenMeteoService()
	if service.BaseURL != "https://api.open-meteo.com/v1/forecast" {
		t.Errorf("Expected BaseURL to be 'https://api.open-meteo.com/v1/forecast', got '%s'", service.BaseURL)
	}
	if service.Client == nil {
		t.Error("Expected Client to be initialized, got nil")
	}
}

func TestGetHourlyForecastErrors(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		latitude       float64
		longitude      float64
		timezone       string
		expectedError  string
	}{
		{
			name: "Request error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Simulate a network error by closing the connection
				hj, ok := w.(http.Hijacker)
				if !ok {
					t.Fatal("couldn't create hijacker")
				}
				conn, _, err := hj.Hijack()
				if err != nil {
					t.Fatal(err)
				}
				conn.Close()
			},
			latitude:      52.52,
			longitude:     13.405,
			timezone:      "Europe/Berlin",
			expectedError: "error making request",
		},
		{
			name: "Unexpected status code",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			latitude:      52.52,
			longitude:     13.405,
			timezone:      "Europe/Berlin",
			expectedError: "unexpected status code: 400",
		},
		{
			name: "Invalid JSON response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"invalid": "json"`))
			},
			latitude:      52.52,
			longitude:     13.405,
			timezone:      "Europe/Berlin",
			expectedError: "error decoding response",
		},
		{
			name: "Invalid timezone",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"hourly": {"time": ["2023-04-14T00:00"], "precipitation_probability": [30], "temperature_2m": [15.5], "cloud_cover": [25]}}`))
			},
			latitude:      52.52,
			longitude:     13.405,
			timezone:      "Invalid/Timezone",
			expectedError: "error loading timezone",
		},
		{
			name: "Invalid time format",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"hourly": {"time": ["invalid-time"], "precipitation_probability": [30], "temperature_2m": [15.5], "cloud_cover": [25]}}`))
			},
			latitude:      52.52,
			longitude:     13.405,
			timezone:      "Europe/Berlin",
			expectedError: "error parsing time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			service := &OpenMeteoService{
				BaseURL: server.URL,
				Client:  server.Client(),
			}

			_, err := service.GetHourlyForecast(tt.latitude, tt.longitude, tt.timezone)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
			}
		})
	}
}
