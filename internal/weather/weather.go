package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WeatherService interface {
	GetHourlyForecast(latitude, longitude float64, timezone string) (HourlyForecast, error)
}

type OpenMeteoService struct {
	BaseURL string
	Client  *http.Client
}

type HourlyForecast struct {
	Time                     []time.Time `json:"time"`
	PrecipitationProbability []float64   `json:"precipitation_probability"`
	Temperature              []float64   `json:"temperature"`
	CloudCover               []float64   `json:"cloud_cover"`
}

type apiResponse struct {
	Hourly struct {
		Time                     []string  `json:"time"`
		PrecipitationProbability []float64 `json:"precipitation_probability"`
		Temperature              []float64 `json:"temperature_2m"`
		CloudCover               []float64 `json:"cloud_cover"`
	} `json:"hourly"`
}

func NewOpenMeteoService() *OpenMeteoService {
	return &OpenMeteoService{
		BaseURL: "https://api.open-meteo.com/v1/forecast",
		Client:  &http.Client{},
	}
}

func (s *OpenMeteoService) GetHourlyForecast(latitude, longitude float64, timezone string) (HourlyForecast, error) {
	url := fmt.Sprintf("%s?latitude=%.6f&longitude=%.6f&hourly=precipitation_probability,temperature_2m,cloud_cover&forecast_days=16&timezone=%s", s.BaseURL, latitude, longitude, timezone)

	resp, err := s.Client.Get(url)
	if err != nil {
		return HourlyForecast{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HourlyForecast{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return HourlyForecast{}, fmt.Errorf("error decoding response: %w", err)
	}

	forecast := HourlyForecast{
		Time:                     make([]time.Time, len(apiResp.Hourly.Time)),
		PrecipitationProbability: apiResp.Hourly.PrecipitationProbability,
		Temperature:              apiResp.Hourly.Temperature,
		CloudCover:               apiResp.Hourly.CloudCover,
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return HourlyForecast{}, fmt.Errorf("error loading timezone: %w", err)
	}

	for i, timeStr := range apiResp.Hourly.Time {
		t, err := time.ParseInLocation("2006-01-02T15:04", timeStr, loc)
		if err != nil {
			return HourlyForecast{}, fmt.Errorf("error parsing time: %w", err)
		}
		forecast.Time[i] = t
	}

	return forecast, nil
}
