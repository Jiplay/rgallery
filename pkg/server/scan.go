package server

import (
	"net/http"

	cache "github.com/patrickmn/go-cache"
	"github.com/robbymilo/rgallery/pkg/scanner"
)

// IsScanInProgress returns true if a scan is currently in progress
func IsScanInProgress() bool {
	return scanner.IsScanInProgress()
}

func Scan(w http.ResponseWriter, r *http.Request, scanType string, cache *cache.Cache) {
	c := r.Context().Value(ConfigKey{}).(Conf)

	var user UserKey
	if r.Context().Value(UserKey{}) != nil {
		user = r.Context().Value(UserKey{}).(UserKey)
	}

	// disable scanning for viewers
	if c.DisableAuth || (!c.DisableAuth && user.UserRole == "admin") || (!c.DisableAuth && user.UserRole == "key") {
		status, err := scanner.Scan(scanType, c, cache)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("503\n"))
			c.Logger.Error("error starting scan", "error", err)
		}

		_, err = w.Write([]byte(status))
		if err != nil {
			c.Logger.Error("error writing scan status", "error", err)
		}
	} else {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

}

func ThumbScan(w http.ResponseWriter, r *http.Request) {
	c := r.Context().Value(ConfigKey{}).(Conf)

	var user UserKey
	if r.Context().Value(UserKey{}) != nil {
		user = r.Context().Value(UserKey{}).(UserKey)
	}

	if c.DisableAuth || (!c.DisableAuth && user.UserRole == "admin") || (!c.DisableAuth && user.UserRole == "key") {
		status, err := scanner.ThumbScan(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("503\n"))
			c.Logger.Error("error starting thumbscan", "error", err)
		}

		_, err = w.Write([]byte(status))
		if err != nil {
			c.Logger.Error("error writing thumbscan status", "error", err)
		}

	} else {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

}
