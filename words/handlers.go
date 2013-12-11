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

// WordsHandler returns similar tags
func WordsHandler(ws WordsService, k int, θ float32) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")
		if tag == "" {
			msg := "Missing tag parameter."
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		if pK := r.URL.Query().Get("k"); pK != "" {
			i, err := strconv.ParseInt(pK, 10, 32)
			if err != nil {
				msg := fmt.Sprintf("Error parsing k: %s", err.Error())
				log.Println(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			k = int(i)
		}

		if pθ := r.URL.Query().Get(":theta"); pθ != "" {
			f, err := strconv.ParseFloat(pθ, 32)
			if err != nil {
				msg := fmt.Sprintf("Error parsing theta: %s", err.Error())
				log.Println(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			θ = float32(f)
		}

		similar, _ := ws(tag, k, θ) // similar will be empty if !ok. this is fine.

		type APIResult struct {
			Tag     string   `json:"tag"`
			Similar []string `json:"similar-tags"`
		}

		json, err := json.Marshal(APIResult{
			Tag:     tag,
			Similar: similar,
		})

		if err != nil {
			log.Println("Error creating JSON.")
			http.Error(w, "Error creating JSON.", http.StatusInternalServerError)
			return
		}

		w.Write(json)
	}

	return middleware(handler)
}

// middleware sets json and cache control headers, adds logging.
func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")

		report.JSON(os.Stdout, next).ServeHTTP(w, r)
	}
}
