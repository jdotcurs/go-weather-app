# Weather Forecast Application

This project is a weather forecast application that provides a 16-day hourly forecast for precipitation probability, temperature, and cloud cover. It's built with Go for the backend API and uses HTML/JavaScript with Chart.js for the frontend. The application demonstrates skills in full-stack development, API integration, and data visualization.

## Features

- Fetch 16-day hourly weather forecasts
- Display precipitation probability, temperature, and cloud cover
- Interactive charts for data visualization
- Automatic geolocation (with user permission)
- Timezone-aware data display
- Reverse geocoding to show location names

## Technologies Used

- Go (Golang) for backend
- HTML/CSS/JavaScript for frontend
- Chart.js for data visualization
- Open-Meteo API for weather data
- OpenStreetMap Nominatim API for reverse geocoding

## Project Structure

- `main.go`: Contains the main Go backend code
- `internal/weather/`: Contains the core logic for fetching weather data
- `templates/index.html`: Frontend interface

## Setup and Running

1. Ensure you have Go installed on your system.

2. Clone the repository:
   ```
   git clone https://github.com/yourusername/weather-forecast-app.git
   cd weather-forecast-app
   ```

3. Build and run the project:
   ```
   go build
   ./weather-forecast-app
   ```

4. Open a web browser and navigate to `http://localhost:8080`.

## API Endpoints

- GET `/`: Serves the main HTML page
- GET `/forecast`: Fetches weather forecast data (requires latitude, longitude, and timezone parameters)

## What This Project Demonstrates

### Go Backend Development:
- Creating a web server using the standard library
- Handling HTTP requests and responses
- JSON encoding/decoding
- Error handling and appropriate HTTP status codes
- External API integration (Open-Meteo)

### Code Organization:
- Structuring a Go application
- Separating concerns (weather service, HTTP handling, etc.)

### Frontend Development:
- HTML/CSS for user interface
- JavaScript for dynamic content and API interactions
- Data visualization with Chart.js
- Geolocation API usage
- Asynchronous operations with Fetch API

### API Design:
- RESTful principles
- Query parameter handling

### Testing:
- Unit tests for weather service
- Test coverage reporting

## Why This Project Was Created

This project serves as a portfolio piece, demonstrating:

- Ability to create a full-stack application
- Understanding of web development concepts
- External API integration
- Data visualization skills
- Frontend and backend communication
- Code organization and best practices in Go
- Test-Driven Development (TDD) approach
- Handling of time zones and geolocation data

It's designed to showcase a range of skills in both backend and frontend development, as well as the ability to work with real-world data and APIs.

## Testing

To run the tests for this project, use the following command:
```
go test ./...
```
For more verbose output, you can use the `-v` flag:

```
go test -v ./...
```

To check the coverage of the tests, you can use the following command:

```
go test -cover ./...
```

or to generate a coverage report in HTML format:

```
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Future Improvements

- Implement caching to reduce API calls
- Add more weather parameters (wind speed, humidity, etc.)
- Enhance error handling and input validation
- Containerize the application using Docker
- Implement a more sophisticated UI with responsive design
