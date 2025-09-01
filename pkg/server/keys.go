package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/robbymilo/rgallery/pkg/render"
	"github.com/robbymilo/rgallery/pkg/users"
)

// CreateKey handles a post request to create a new api key.
func CreateKey(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(ParamsKey{}).(FilterParams)
	c := r.Context().Value(ConfigKey{}).(Conf)

	// get signin as JSON
	creds := &ApiCredentials{}
	if params.Json {
		err := json.NewDecoder(r.Body).Decode(creds)
		if err != nil {
			c.Logger.Error("error decoding json", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("400\n"))
			if err != nil {
				c.Logger.Error("error writing 400 response", "error", err)
			}

			return
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			c.Logger.Error("error parsing form", "error", err)
		}
		creds.Name = r.Form["name"][0]
	}

	// create error URL
	errorUrl, _ := url.Parse("/admin")
	errorParams := url.Values{}
	errorParams.Add("error", "true")
	errorUrl.RawQuery = errorParams.Encode()

	var user UserKey
	if r.Context().Value(UserKey{}) != nil {
		user = r.Context().Value(UserKey{}).(UserKey)
	}

	if user.UserRole == "viewer" {
		c.Logger.Error("error adding key", "error", fmt.Errorf("api key cannot be added by viewers"))
		http.Redirect(w, r, errorUrl.String(), http.StatusFound)
		return
	}

	key, err := users.AddKey(creds, c)
	if err != nil {
		c.Logger.Error("error adding key", "error", err)
		http.Redirect(w, r, errorUrl.String(), http.StatusFound)
		return
	}

	// key added
	var credentials ApiCredentials
	credentials.Name = creds.Name
	credentials.Key = key

	c.Logger.Info("api key added: " + creds.Name)

	// show the created API key once
	keys, err := users.GetKeyNames(c)
	if err != nil {
		c.Logger.Error("error listing keys", "error", err)
		http.Redirect(w, r, errorUrl.String(), http.StatusFound)
		return
	}

	users, err := users.ListUsers(c)
	if err != nil {
		c.Logger.Error("error listing users", "error", err)
	}

	// needed to show the API key once
	response := ResponseAdmin{
		HideNavFooter: false,
		HideAuth:      c.DisableAuth,
		Section:       "admin",
		Key:           credentials,
		Keys:          keys,
		Users:         users,
		UserName:      user.UserName,
		UserRole:      user.UserRole,
		Meta:          c.Meta,
	}
	err = render.Render(w, r, response, "admin")
	if err != nil {
		c.Logger.Error("error rendering admin create key response", "error", err)
	}

}

// RemoveKey handles a post request to remove an API key.
func RemoveKey(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(ParamsKey{}).(FilterParams)
	c := r.Context().Value(ConfigKey{}).(Conf)

	// get signin as JSON
	creds := &ApiCredentials{}
	if params.Json {
		err := json.NewDecoder(r.Body).Decode(creds)
		if err != nil {
			c.Logger.Error("error decoding json", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("400\n"))
			if err != nil {
				c.Logger.Error("error writing 400 response", "error", err)
			}

			return
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			c.Logger.Error("error parsing form", "error", err)
		}
		creds.Name = r.Form["name"][0]
	}

	err := users.RemoveKey(creds, c)
	if err != nil {
		c.Logger.Error("error removing key", "error", err)
		return
	}

	c.Logger.Info("api key removed: " + creds.Name)

	http.Redirect(w, r, "/admin", http.StatusFound)

}
