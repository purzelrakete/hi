package words

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/handy/report"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Handler returns similar tags
func Handler(ws Service, max, minFq int, minSimilarity float32) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get(":tag")
		if tag == "" {
			errf(w, "Missing tag parameter.")
			return
		}

		if pMax := r.URL.Query().Get("max-results"); pMax != "" {
			i, err := strconv.ParseInt(pMax, 10, 32)
			if err != nil {
				errf(w, "Error parsing k: %s", err.Error())
				return
			}

			max = int(i)
		}

		if pMinSimilarity := r.URL.Query().Get("min-similarity"); pMinSimilarity != "" {
			f, err := strconv.ParseFloat(pMinSimilarity, 32)
			if err != nil {
				errf(w, "Error parsing theta: %s", err.Error())
				return
			}

			minSimilarity = float32(f)
		}

		if pMinFq := r.URL.Query().Get("min-frequency"); pMinFq != "" {
			i, err := strconv.ParseInt(pMinFq, 10, 32)
			if err != nil {
				errf(w, "Error parsing min frequency: %s", err.Error())
				return
			}

			minFq = int(i)
		}

		// similar will be empty if !ok. this is fine.
		similar, _, _ := ws.NN(tag, max, minFq, minSimilarity)

		type APIResult struct {
			Tag     string `json:"tag"`
			Similar []Hit  `json:"hits"`
		}

		json, err := json.Marshal(APIResult{
			Tag:     tag,
			Similar: similar,
		})

		if err != nil {
			errf(w, "Error creating JSON.")
			return
		}

		w.Write(json)
	}

	return middleware(handler)
}

func errf(w http.ResponseWriter, msg string, things ...interface{}) {
	log.Println(fmt.Sprintf(msg, things...))
	http.Error(w, msg, http.StatusBadRequest)
}

// middleware sets json and cache control headers, adds logging.
func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=60")

		report.JSON(os.Stdout, next).ServeHTTP(w, r)
	}
}
