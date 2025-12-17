package services

import (
	"encoding/json"
	"fmt"
	"log"
	"police-assistant-backend/config"

	"github.com/go-resty/resty/v2"
)

const (
	orsDirectionsURL = "https://api.openrouteservice.org/v2/directions/driving-car"
	orsGeocodeURL    = "https://api.openrouteservice.org/geocode/search"
	orsReverseURL    = "https://api.openrouteservice.org/geocode/reverse"
)

type ORSService struct {
	client *resty.Client
}

func NewORSService() *ORSService {
	client := resty.New()
	client.SetHeader("Authorization", config.AppConfig.ORSAPIKey)
	client.SetHeader("Accept", "application/json, application/geo+json")
	client.SetHeader("Content-Type", "application/json")

	log.Println("✅ OpenRouteService initialized (Free Maps API)")

	return &ORSService{
		client: client,
	}
}

// GetTrafficInfo gets current route information around a location
func (s *ORSService) GetTrafficInfo(lat, lng float64) (map[string]interface{}, error) {
	// Create a small route to nearby point to estimate traffic
	destLat := lat + 0.01 // ~1km away
	destLng := lng + 0.01

	// ORS uses POST with JSON body for directions
	requestBody := map[string]interface{}{
		"coordinates": [][]float64{
			{lng, lat},         // Start point (longitude, latitude)
			{destLng, destLat}, // End point
		},
	}

	log.Printf("🗺️  Getting traffic info for: %.6f, %.6f", lat, lng)

	var result map[string]interface{}
	resp, err := s.client.R().
		SetHeader("Accept", "application/json, application/geo+json").
		SetBody(requestBody).
		SetResult(&result).
		Post(orsDirectionsURL)

	if err != nil {
		return nil, fmt.Errorf("failed to get directions: %w", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("❌ ORS API error: Status %d, Body: %s", resp.StatusCode(), resp.String())
		return nil, fmt.Errorf("ORS API returned status %d: %s", resp.StatusCode(), resp.String())
	}

	// Parse ORS response
	routes, ok := result["routes"].([]interface{})
	if !ok || len(routes) == 0 {
		return map[string]interface{}{
			"status":    "no_data",
			"message":   "Tidak ada data rute tersedia",
			"condition": "unknown",
		}, nil
	}

	route := routes[0].(map[string]interface{})
	summary := route["summary"].(map[string]interface{})

	distance := summary["distance"].(float64) / 1000 // Convert to km
	duration := summary["duration"].(float64) / 60   // Convert to minutes

	// Estimate traffic based on speed
	avgSpeed := (distance / duration) * 60 // km/h

	var condition string
	var conditionEmoji string

	if avgSpeed < 20 {
		condition = "heavy"
		conditionEmoji = "🔴"
	} else if avgSpeed < 40 {
		condition = "moderate"
		conditionEmoji = "🟡"
	} else {
		condition = "light"
		conditionEmoji = "🟢"
	}

	trafficInfo := map[string]interface{}{
		"status":          "success",
		"distance":        fmt.Sprintf("%.2f km", distance),
		"duration":        fmt.Sprintf("%.1f min", duration),
		"avg_speed":       fmt.Sprintf("%.1f km/h", avgSpeed),
		"condition":       condition,
		"condition_emoji": conditionEmoji,
	}

	log.Printf("✅ Traffic condition: %s %s (avg speed: %.1f km/h)", conditionEmoji, condition, avgSpeed)

	return trafficInfo, nil
}

// GetAlternativeRoutes gets multiple route options
func (s *ORSService) GetAlternativeRoutes(origin, destination string) ([]map[string]interface{}, error) {
	// First, geocode origin and destination
	originCoords, err := s.geocode(origin)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode origin '%s': %w", origin, err)
	}

	destCoords, err := s.geocode(destination)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode destination '%s': %w", destination, err)
	}

	log.Printf("🗺️  Finding routes from %s (%.4f,%.4f) to %s (%.4f,%.4f)",
		origin, originCoords["lng"], originCoords["lat"],
		destination, destCoords["lng"], destCoords["lat"])

	// Simple request without alternative routes first (for debugging)
	requestBody := map[string]interface{}{
		"coordinates": [][]float64{
			{originCoords["lng"].(float64), originCoords["lat"].(float64)},
			{destCoords["lng"].(float64), destCoords["lat"].(float64)},
		},
	}

	// Debug log
	reqBodyBytes, _ := json.Marshal(requestBody)
	log.Printf("📤 Request body: %s", string(reqBodyBytes))

	var result map[string]interface{}
	resp, err := s.client.R().
		SetHeader("Accept", "application/json, application/geo+json").
		SetBody(requestBody).
		SetResult(&result).
		Post(orsDirectionsURL)

	if err != nil {
		return nil, fmt.Errorf("failed to get routes: %w", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("❌ ORS API error: Status %d, Body: %s", resp.StatusCode(), resp.String())
		return nil, fmt.Errorf("ORS API returned status %d: %s", resp.StatusCode(), resp.String())
	}

	// Parse routes
	routesData, ok := result["routes"].([]interface{})
	if !ok || len(routesData) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	var routes []map[string]interface{}

	for i, routeData := range routesData {
		route := routeData.(map[string]interface{})
		summary := route["summary"].(map[string]interface{})

		distance := summary["distance"].(float64) / 1000 // km
		duration := summary["duration"].(float64) / 60   // minutes
		avgSpeed := (distance / duration) * 60           // km/h

		// Determine traffic condition
		var condition string
		var conditionEmoji string

		if avgSpeed < 20 {
			condition = "heavy"
			conditionEmoji = "🔴"
		} else if avgSpeed < 40 {
			condition = "moderate"
			conditionEmoji = "🟡"
		} else {
			condition = "light"
			conditionEmoji = "🟢"
		}

		// Extract steps if available
		var steps []map[string]interface{}
		if segments, ok := route["segments"].([]interface{}); ok && len(segments) > 0 {
			segment := segments[0].(map[string]interface{})
			if stepsData, ok := segment["steps"].([]interface{}); ok {
				for j, step := range stepsData {
					if j >= 5 { // Limit to first 5 steps
						break
					}
					s := step.(map[string]interface{})
					instruction := ""
					if inst, ok := s["instruction"].(string); ok {
						instruction = inst
					}
					steps = append(steps, map[string]interface{}{
						"instruction": instruction,
						"distance":    fmt.Sprintf("%.2f km", s["distance"].(float64)/1000),
						"duration":    fmt.Sprintf("%.1f min", s["duration"].(float64)/60),
					})
				}
			}
		}

		routeInfo := map[string]interface{}{
			"route_number":      i + 1,
			"summary":           fmt.Sprintf("Rute %d via OpenStreetMap", i+1),
			"distance":          fmt.Sprintf("%.2f km", distance),
			"duration":          fmt.Sprintf("%.0f min", duration),
			"avg_speed":         fmt.Sprintf("%.1f km/h", avgSpeed),
			"traffic_condition": condition,
			"condition_emoji":   conditionEmoji,
			"start_address":     origin,
			"end_address":       destination,
			"steps":             steps,
			"total_steps":       len(steps),
		}

		routes = append(routes, routeInfo)
	}

	log.Printf("✅ Found %d alternative route(s)", len(routes))

	return routes, nil
}

// ReverseGeocode converts coordinates to address using Nominatim (OSM)
func (s *ORSService) ReverseGeocode(lat, lng float64) (string, error) {
	params := map[string]string{
		"point.lon": fmt.Sprintf("%.6f", lng),
		"point.lat": fmt.Sprintf("%.6f", lat),
		"size":      "1",
	}

	var result map[string]interface{}
	resp, err := s.client.R().
		SetQueryParams(params).
		SetResult(&result).
		Get(orsReverseURL)

	if err != nil {
		return "", fmt.Errorf("failed to reverse geocode: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "Lokasi tidak diketahui", nil
	}

	features, ok := result["features"].([]interface{})
	if !ok || len(features) == 0 {
		return "Lokasi tidak diketahui", nil
	}

	feature := features[0].(map[string]interface{})
	properties := feature["properties"].(map[string]interface{})

	if label, ok := properties["label"].(string); ok {
		return label, nil
	}

	return "Lokasi tidak diketahui", nil
}

// geocode converts address to coordinates
func (s *ORSService) geocode(address string) (map[string]interface{}, error) {
	params := map[string]string{
		"text": address,
		"size": "1",
	}

	var result map[string]interface{}
	resp, err := s.client.R().
		SetQueryParams(params).
		SetResult(&result).
		Get(orsGeocodeURL)

	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		body := resp.String()
		log.Printf("❌ Geocoding error: %s", body)
		return nil, fmt.Errorf("geocoding returned status %d", resp.StatusCode())
	}

	features, ok := result["features"].([]interface{})
	if !ok || len(features) == 0 {
		return nil, fmt.Errorf("address not found: %s", address)
	}

	feature := features[0].(map[string]interface{})
	geometry := feature["geometry"].(map[string]interface{})
	coordinates := geometry["coordinates"].([]interface{})

	// Debug log
	coordBytes, _ := json.Marshal(coordinates)
	log.Printf("🗺️  Geocoded '%s' to: %s", address, string(coordBytes))

	return map[string]interface{}{
		"lng": coordinates[0].(float64),
		"lat": coordinates[1].(float64),
	}, nil
}
