package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
	"yottadb/dbdriver"
	"yottadb/dbdriver/document"
	"yottadb/dbdriver/keyvalue"
)

type Config struct {
	NodeTree *[]string
	Port     string
	HashKey  string
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func HttpHandlerFactory(config Config) (func(http.ResponseWriter, *http.Request), error) {

	dd, err := document.New(config.HashKey, config.NodeTree)
	kvd, err := keyvalue.New(config.HashKey, config.NodeTree)
	if err != nil {
		log.Println("Error instantiating driver: ", err)
		return nil, err
	}

	handler := func(w http.ResponseWriter, r *http.Request) {

		var req dbdriver.Request
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Request error: ", err)
			if _, err := w.Write([]byte("Malformed YottaDB request")); err != nil {
				log.Println("ERROR: ", err)
			}
			return
		}

		log.Println("Request: ", req)

		switch req.Driver {

		case "document":
			document.HttpHandler(w, req, dd)

		case "collection":
			document.HttpHandler(w, req, dd)

		case "keyvalue":
			keyvalue.HttpHandler(w, req, kvd)

		case "columnar":

		case "pubsub":

		default:
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("YottaDB driver not found")); err != nil {
				log.Println("ERROR: ", err)
			}
		}

	}

	return handler, nil

}
