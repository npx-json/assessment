package api

import (
	"database/sql"
	"encoding/json"
	"ipcheck/internal/geoip"
	"log"
	"net/http"
)

type CheckRequest struct {
	IP string `json:"ip"`
}

type CheckResponse struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Allowed bool   `json:"allowed"`
	Message string `json:"message"`
}

var sqlDb *sql.DB // Declare the global `db` variable

func SetDB(database *sql.DB) {
	sqlDb = database // Assign the global `db` variable
}

func MakeCheckHandler(s *geoip.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		// Fetch allowed countries from the database
		allowedCountries, err := GetAllowedCountries()
		if err != nil {
			http.Error(w, "Error fetching allowed countries: "+err.Error(), http.StatusExpectationFailed)
			return
		}

		var req CheckRequest
		allowed, country, err := s.CheckIP(req.IP, allowedCountries)

		if err != nil {
			http.Error(w, "Error checking IP: "+err.Error(), http.StatusExpectationFailed)
			return
		}

		res := CheckResponse{
			IP:      req.IP,
			Country: country,
			Allowed: allowed,
			Message: "Access Denied",
		}
		if allowed {
			res.Message = "Access Allowed"
		}
		json.NewEncoder(w).Encode(res)
	}
}

func GetAllowedCountries() ([]string, error) {
	rows, err := sqlDb.Query("SELECT country FROM whitelisted_tb")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []string
	for rows.Next() {
		var country string
		if err := rows.Scan(&country); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		countries = append(countries, country)
	}

	return countries, nil
}
