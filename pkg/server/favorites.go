package server

import (
	"net/http"

	"github.com/robbymilo/rgallery/pkg/queries"
	"github.com/robbymilo/rgallery/pkg/render"
)

func ServeFavorites(w http.ResponseWriter, r *http.Request) {

	params := r.Context().Value(ParamsKey{}).(FilterParams)
	c := r.Context().Value(ConfigKey{}).(Conf)

	page := params.Page
	var pageSize = 100
	var rating = 5
	offset := (page - 1) * pageSize
	media, err := queries.GetFavorites(pageSize, offset, rating, params, c)
	if err != nil {
		c.Logger.Error("error getting favorites", "error", err)
	}
	total, err := queries.GetTotalFavorites(rating, c)
	if err != nil {
		c.Logger.Error("error getting total of favorites", "error", err)
	}

	response := ResponseMediaItems{
		Title:         "Favorites",
		Collection:    "favorites",
		MediaItems:    media,
		OrderBy:       "date",
		Page:          page,
		PageSize:      pageSize,
		Total:         total,
		Direction:     params.Direction,
		Section:       "favorites",
		HideNavFooter: false,
		Meta:          c.Meta,
	}

	err = render.Render(w, r, response, "images")
	if err != nil {
		c.Logger.Error("error rendering favorites response", "error", err)
	}
}
