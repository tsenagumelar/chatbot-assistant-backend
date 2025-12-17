package services

import (
	"encoding/json"
	"fmt"
	"log"
	"police-assistant-backend/config"
	"strconv"
	"strings"

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
	// Parse or geocode origin
	originCoords, err := s.parseOrGeocode(origin)
	if err != nil {
		return nil, fmt.Errorf("failed to process origin '%s': %w", origin, err)
	}

	// Parse or geocode destination
	destCoords, err := s.parseOrGeocode(destination)
	if err != nil {
		return nil, fmt.Errorf("failed to process destination '%s': %w", destination, err)
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

		// Extract ALL steps (turn-by-turn directions)
		var steps []map[string]interface{}
		if segments, ok := route["segments"].([]interface{}); ok && len(segments) > 0 {
			for _, seg := range segments {
				segment := seg.(map[string]interface{})
				if stepsData, ok := segment["steps"].([]interface{}); ok {
					for _, step := range stepsData {
						s := step.(map[string]interface{})

						// Get instruction
						instruction := ""
						if inst, ok := s["instruction"].(string); ok {
							instruction = inst
						}

						// Get road name if available
						roadName := ""
						if name, ok := s["name"].(string); ok && name != "" && name != "-" {
							roadName = name
						}

						// Get step type (turn left, turn right, straight, etc)
						stepType := ""
						if sType, ok := s["type"].(float64); ok {
							stepType = getStepTypeName(int(sType))
						}

						// Get distance and duration
						stepDistance := 0.0
						if dist, ok := s["distance"].(float64); ok {
							stepDistance = dist
						}

						stepDuration := 0.0
						if dur, ok := s["duration"].(float64); ok {
							stepDuration = dur
						}

						// Build detailed step info
						stepInfo := map[string]interface{}{
							"step_number": len(steps) + 1,
							"instruction": instruction,
							"road_name":   roadName,
							"type":        stepType,
							"distance":    fmt.Sprintf("%.2f km", stepDistance/1000),
							"distance_m":  stepDistance,
							"duration":    fmt.Sprintf("%.1f min", stepDuration/60),
							"duration_s":  stepDuration,
						}

						steps = append(steps, stepInfo)
					}
				}
			}
		}

		// Log step count for debugging
		log.Printf("   Route %d: %d steps extracted", i+1, len(steps))

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

// getStepTypeName converts ORS step type code to human-readable name
func getStepTypeName(stepType int) string {
	switch stepType {
	case 0:
		return "Berangkat"
	case 1:
		return "Lurus"
	case 2:
		return "Belok kanan sedikit"
	case 3:
		return "Belok kanan"
	case 4:
		return "Belok kanan tajam"
	case 5:
		return "Putar balik"
	case 6:
		return "Belok kiri tajam"
	case 7:
		return "Belok kiri"
	case 8:
		return "Belok kiri sedikit"
	case 9:
		return "Terus lurus"
	case 10:
		return "Masuk bundaran"
	case 11:
		return "Keluar bundaran"
	case 12:
		return "Tetap di bundaran"
	case 13:
		return "Terus"
	case 14:
		return "Masuk jalan raya"
	case 15:
		return "Sampai tujuan"
	default:
		return "Lanjutkan"
	}
}

// parseOrGeocode tries to parse coordinates from string, or geocode if it's an address
func (s *ORSService) parseOrGeocode(location string) (map[string]interface{}, error) {
	// Try to parse as coordinates first (format: "lat,lng" or "lat, lng")
	location = strings.TrimSpace(location)
	parts := strings.Split(location, ",")

	if len(parts) == 2 {
		latStr := strings.TrimSpace(parts[0])
		lngStr := strings.TrimSpace(parts[1])

		lat, errLat := strconv.ParseFloat(latStr, 64)
		lng, errLng := strconv.ParseFloat(lngStr, 64)

		// If both parse successfully, it's coordinates
		if errLat == nil && errLng == nil {
			// Validate coordinate ranges
			if lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180 {
				log.Printf("🗺️  Parsed coordinates: %.6f, %.6f", lat, lng)
				return map[string]interface{}{
					"lat": lat,
					"lng": lng,
				}, nil
			}
		}
	}

	// If not valid coordinates, treat as address and geocode
	log.Printf("🔍 Geocoding address: %s", location)
	return s.geocode(location)
}

// geocode converts address to coordinates
func (s *ORSService) geocode(address string) (map[string]interface{}, error) {
	// First attempt with full address
	coords, err := s.geocodeAttempt(address)
	if err == nil {
		return coords, nil
	}

	log.Printf("⚠️  Full address geocoding failed, trying simplified query...")

	// Second attempt: Extract city/locality from address
	// Common patterns: "Something, City, Province" or "Something City"
	simplifiedAddress := s.extractMainLocation(address)
	if simplifiedAddress != address {
		log.Printf("🔍 Trying with simplified address: %s", simplifiedAddress)
		coords, err = s.geocodeAttempt(simplifiedAddress)
		if err == nil {
			return coords, nil
		}
	}

	return nil, fmt.Errorf("could not geocode address: %s", address)
}

// extractMainLocation tries to extract the main city/locality from full address
func (s *ORSService) extractMainLocation(address string) string {
	// Split by comma and look for city names
	parts := strings.Split(address, ",")

	// Try to find the most relevant part (usually city name)
	// Common Indonesian cities/areas to prioritize
	keywords := []string{
		"Sukabumi", "Jakarta", "Bandung", "Bogor", "Depok",
		"Tangerang", "Bekasi", "Cianjur", "Purwakarta",
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(part), strings.ToLower(keyword)) {
				return keyword
			}
		}
	}

	// If no known city found, return first non-street part
	if len(parts) >= 2 {
		// Skip first part (usually street/building), return city
		return strings.TrimSpace(parts[1])
	}

	return address
}

// geocodeAttempt performs a single geocoding attempt
func (s *ORSService) geocodeAttempt(address string) (map[string]interface{}, error) {
	params := map[string]string{
		"text": address,
		"size": "5",
		// Focus on Indonesia region for better results
		"boundary.country": "IDN",
		"focus.point.lat":  "-6.2",
		"focus.point.lon":  "106.8",
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

	// Log all found locations for debugging
	log.Printf("🔍 Found %d location(s) for '%s':", len(features), address)
	for i, f := range features {
		feat := f.(map[string]interface{})
		props := feat["properties"].(map[string]interface{})
		if label, ok := props["label"].(string); ok {
			log.Printf("   %d. %s", i+1, label)
		}
	}

	// Find the best match (not just "Indonesia")
	var selectedFeature map[string]interface{}
	var selectedLabel string

	for _, f := range features {
		feat := f.(map[string]interface{})
		props := feat["properties"].(map[string]interface{})

		if label, ok := props["label"].(string); ok {
			// Skip results that are too generic (just country name)
			if label == "Indonesia" || label == "Java" {
				continue
			}

			// Check if this is a valid location (has locality, region, or county)
			if locality, hasLoc := props["locality"].(string); hasLoc && locality != "" {
				selectedFeature = feat
				selectedLabel = label
				break
			}

			if region, hasReg := props["region"].(string); hasReg && region != "" {
				selectedFeature = feat
				selectedLabel = label
				break
			}

			if county, hasCounty := props["county"].(string); hasCounty && county != "" {
				selectedFeature = feat
				selectedLabel = label
				break
			}

			// If no specific fields but label is not too generic, use it
			if len(strings.Split(label, ",")) > 1 {
				selectedFeature = feat
				selectedLabel = label
				break
			}
		}
	}

	// If no good match found, return error
	if selectedFeature == nil {
		return nil, fmt.Errorf("no specific location found for: %s", address)
	}

	geometry := selectedFeature["geometry"].(map[string]interface{})
	coordinates := geometry["coordinates"].([]interface{})

	// Log selected location
	log.Printf("✅ Selected location: %s", selectedLabel)

	// Debug log coordinates
	coordBytes, _ := json.Marshal(coordinates)
	log.Printf("🗺️  Coordinates: %s", string(coordBytes))

	return map[string]interface{}{
		"lng": coordinates[0].(float64),
		"lat": coordinates[1].(float64),
	}, nil
}
