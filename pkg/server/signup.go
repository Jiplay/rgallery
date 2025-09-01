package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/robbymilo/rgallery/pkg/users"
)

// SignUp handles a post request to create a new user.
func SignUp(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(ParamsKey{}).(FilterParams)
	c := r.Context().Value(ConfigKey{}).(Conf)

	// get signin as JSON
	creds := &UserCredentials{}
	if params.Json {
		err := json.NewDecoder(r.Body).Decode(creds)
		if err != nil {
			fmt.Println("error decoding json:", err)

			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("400\n"))
			if err != nil {
				fmt.Println("error writing 400 response:", err)
			}

			return
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("error parsing form:", err)
		}
		creds.Username = r.Form["username"][0]
		creds.Password = r.Form["password"][0]
		creds.Role = r.Form["role"][0]
	}

	// create error URL
	errorUrl, _ := url.Parse("/adduser")
	errorParams := url.Values{}
	errorParams.Add("error", "true")
	errorUrl.RawQuery = errorParams.Encode()

	err := users.AddUser(*creds, c)
	if err != nil {
		fmt.Println("error adding user:", err)
		http.Redirect(w, r, errorUrl.String(), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)

}
