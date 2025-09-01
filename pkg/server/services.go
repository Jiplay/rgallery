package server

import "net/http"

// CheckResizeServiceHealth confirms the resize service is available to resize images.
func CheckResizeServiceHealth(c Conf) error {

	res, err := http.Get(c.ResizeService + "/healthz")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// CheckLocationServiceHealth confirms the resize service is available to resize images.
func CheckLocationServiceHealth(c Conf) error {

	res, err := http.Get(c.LocationService + "/healthz")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
