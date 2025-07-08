package main

type GeoIPResult struct {
	IP           string  `json:"ip"`
	CountryCode  string  `json:"country_code"`
	Country      string  `json:"country"`
	Region       string  `json:"region"`
	City         string  `json:"city"`
	PostalCode   string  `json:"postal_code"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Organization string  `json:"organization"`
	Timezone     string  `json:"timezone"`
}
