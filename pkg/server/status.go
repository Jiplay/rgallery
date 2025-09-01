package server

import (
	"fmt"
	"net/http"

	"github.com/robbymilo/rgallery/pkg/scanner"
)

// ServeStatus handles requests for status on initial page loads.
func ServeStatus(w http.ResponseWriter, r *http.Request, c Conf) {
	scanInProgress := scanner.IsScanInProgress()
	_, err := w.Write([]byte(fmt.Sprintf("%t", scanInProgress)))
	if err != nil {
		c.Logger.Error("error writing scan status", "error", err)
	}
}
