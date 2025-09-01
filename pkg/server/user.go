package server

import (
	"net/http"

	"github.com/robbymilo/rgallery/pkg/render"
)

// ServeAddUser serves the add user page.
func ServeAddUser(w http.ResponseWriter, r *http.Request) {
	response := ResponsAuth{
		HideNavFooter: false,
		Section:       "auth",
	}
	_ = render.Render(w, r, response, "adduser")
}
