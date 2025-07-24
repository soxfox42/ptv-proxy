package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// If rate limits need to be enforced down the track,
// having unique tokens will be useful.
// Adding them now means that current clients will work with future APIs.
func makeToken() string {
	bytes := make([]byte, 24)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

type SearchResult struct {
	Stops []any
}

type DeparturesResult struct {
	Departures []any
}

func main() {
	http.HandleFunc("/generateToken", func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]string{"token": makeToken()}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/searchStops", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		searchTerm := r.FormValue("searchTerm")

		if token == "" || searchTerm == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		api := fmt.Sprintf("/search/%s", url.PathEscape(searchTerm))
		ptvResponse, ok := ptvRequest[SearchResult](api)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := map[string]any{
			"results": ptvResponse.Stops,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/getDepartures", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		routeType := r.FormValue("routeType")
		stopID := r.FormValue("stopID")

		if token == "" || routeType == "" || stopID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		api := fmt.Sprintf("/departures/route_type/%s/stop/%s?expand=3&expand=4", routeType, stopID)
		ptvRoutesResponse, ok := ptvRequest[any](api)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ptvRoutesResponse)
	})

	fmt.Println("listening on port 5690")
	http.ListenAndServe("localhost:5690", nil)
}
