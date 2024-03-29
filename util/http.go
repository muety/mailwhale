package util

import (
	"encoding/json"
	"github.com/emvi/logbuch"
	"net/http"
)

func RespondEmpty(w http.ResponseWriter, r *http.Request, status int) {
	if status <= 0 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
}

func RespondJson(w http.ResponseWriter, status int, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(object); err != nil {
		logbuch.Error("error while writing json response: %v", err)
	}
}

func RespondHtml(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	w.Write([]byte(data))
}

func RespondError(w http.ResponseWriter, r *http.Request, status int, err error) {
	logbuch.Error("request '%s %s' failed: %v", r.Method, r.URL.Path, err)
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
}

func RespondErrorMessage(w http.ResponseWriter, r *http.Request, status int, err error) {
	logbuch.Error("request '%s %s' failed: %v", r.Method, r.URL.Path, err)
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
