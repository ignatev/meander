package  main

import (
	"net/http"
	"encoding/json"
//	"runtime"
	"github.com/ignatev/meander"
	"os"
	"strings"
	"strconv"
	"log"
)

func main(){
	meander.APIKey = os.Getenv("GOOGLE_PLACES_API_KEY")
	http.HandleFunc("/journeys", cors(func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	}))
	http.HandleFunc("/recommendations", cors(func(w http.ResponseWriter, r *http.Request) {
	q := &meander.Query {
		Journey: strings.Split(r.URL.Query().Get("journey"), "|"),
	}
	var err error
	q.Lat, err = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	log.Println(q.Lat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q.Lng, err = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	log.Println(q.Lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q.Radius, err = strconv.Atoi(r.URL.Query().Get("radius"))
	log.Println(q.Radius)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q.CostRangeStr = r.URL.Query().Get("cost")
	places := q.Run()
	respond(w, r, places)
}))
	log.Println("serving meander API on :8080")
	http.ListenAndServe(":8080", http.DefaultServeMux)
}
func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))
	for i, d := range data {
		publicData[i] = meander.Public(d)
	}

	return json.NewEncoder(w).Encode(publicData)
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}