package words

import (
	"encoding/json"
	"github.com/streadway/handy/report"
	"log"
	"net/http"
	"os"
)

// WordsHandler returns similar tags
func WordsHandler(ws WordsService) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get(":tag")
		if tag == "" {
			msg := "Missing tag parameter."
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		similar, _ := ws(tag, 5) // similar will be empty if !ok. this is fine.

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
