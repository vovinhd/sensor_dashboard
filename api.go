package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sensor_dashboard/db"
	"time"
)

//go:embed public
var staticFiles embed.FS // Holds embedded static files

func Api(queries db.Queries, config Config) (http.Handler, error) {

	mux := http.NewServeMux()

	var publicDir fs.FS
	if config.ServeFromFS {
		publicDir = os.DirFS("./public")
	} else {
		publicDir, _ = fs.Sub(staticFiles, "public")
	}

	fileServer := http.FileServer(http.FS(publicDir))

	mux.Handle("/humidity", GetHumidityLogByDate(queries, config))
	mux.Handle("/power", GetPowerLogByDate(queries, config))
	mux.Handle("/state", GetSwitchStateLogByDate(queries, config))
	mux.Handle("/tags/devices", GetDevicesForTag(queries, config))
	mux.Handle("/tags", GetTags(queries, config))

	mux.Handle("/", fileServer)

	return mux, nil
}

func tagAndDateFromQuery(r *http.Request) (string, time.Time, error) {
	q := r.URL.Query()
	tag := q.Get("tag")
	startDateStr := q.Get("from")
	now := time.Now()

	if tag == "" {
		return "", now, errors.New("no tag specified")
	}

	var startDate time.Time
	if startDateStr == "" || startDateStr == "null" {
		startDate = now.AddDate(0, 0, -7)
	} else {
		startDate, _ = time.Parse("2006-01-02T15:04:05Z", startDateStr)
	}

	return tag, startDate, nil
}

func GetTags(queries db.Queries, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rows, err := queries.ListDevices(context.Background())

		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(rows); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error encoding JSON: %v", err)
		}

	})
}

func GetDevicesForTag(queries db.Queries, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()
		tag := q.Get("tag")
		if tag == "" {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		rows, err := queries.ListSensorsForDevice(context.Background(), tag)

		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(rows); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error encoding JSON: %v", err)
		}

	})
}

func GetPowerLogByDate(queries db.Queries, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag, startDate, err := tagAndDateFromQuery(r)

		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			log.Println(err)
			return
		}

		rows, err := queries.PowerLogForDeviceToDate(context.Background(), db.PowerLogForDeviceToDateParams{
			Tag:      tag,
			FromTime: startDate,
			ToTime:   time.Now(),
		})

		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(rows); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error encoding JSON: %v", err)
		}

	})
}

func GetSwitchStateLogByDate(queries db.Queries, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag, startDate, err := tagAndDateFromQuery(r)

		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			log.Println(err)
			return
		}

		rows, err := queries.SwitchStateLogForDeviceToDate(context.Background(), db.SwitchStateLogForDeviceToDateParams{
			Tag:      tag,
			FromTime: startDate,
			ToTime:   time.Now(),
		})

		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(rows); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error encoding JSON: %v", err)
		}

	})
}

func GetHumidityLogByDate(queries db.Queries, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag, startDate, err := tagAndDateFromQuery(r)

		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			log.Println(err)
			return
		}

		rows, err := queries.HumidityLogForDeviceToDate(context.Background(), db.HumidityLogForDeviceToDateParams{
			Tag:      tag,
			FromTime: startDate,
			ToTime:   time.Now(),
		})

		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(rows); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error encoding JSON: %v", err)
		}

	})
}
