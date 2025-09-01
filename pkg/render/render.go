package render

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/robbymilo/rgallery/pkg/static"
	"github.com/robbymilo/rgallery/pkg/templates"
	"github.com/robbymilo/rgallery/pkg/types"
)

type FilterParams = types.FilterParams
type Conf = types.Conf
type ConfigKey = types.ConfigKey
type ParamsKey = types.ParamsKey
type UserKey = types.UserKey

var funcs = template.FuncMap{
	"formatInt": func(number int) string {
		output := strconv.Itoa(number)
		startOffset := 3
		if number < 0 {
			startOffset++
		}
		for outputIndex := len(output); outputIndex > startOffset; {
			outputIndex -= 3
			output = output[:outputIndex] + "," + output[outputIndex:]
		}
		return output
	},
	"formatFloat": func(f float64, format string, precision int) string {
		return strconv.FormatFloat(f, format[0], precision, 64)
	},
	"formatDate": func(t time.Time) string {
		return t.Format("Mon, 02 January 2006 15:04:05.000")
	},
	"formatLocalDate": func(t time.Time, offset float64) string {
		// Add the offset to the time (offset is in minutes)
		t = t.Add(time.Duration(offset) * time.Minute)

		// Format the time into the desired format: "Mon, 02 Jan 2006 15:04:05.000"
		return t.Format("Mon, 02 January 2006 15:04:05.000")
	},
	"readFile": func(path string) template.HTML {
		f, err := static.StaticDir.ReadFile(path)

		if err != nil {
			fmt.Println("error reading file", err)
		}

		return template.HTML(string(f))
	},
	"urlquery": url.QueryEscape,
}

// Render coordinates the sending of raw JSON, an embedded template, or a local template to the response.
func Render(w http.ResponseWriter, r *http.Request, response interface{}, layout string) error {
	c := r.Context().Value(ConfigKey{}).(Conf)
	params := r.Context().Value(ParamsKey{}).(FilterParams)
	var user UserKey
	if r.Context().Value(UserKey{}) != nil {
		user = r.Context().Value(UserKey{}).(UserKey)
	}

	e := setEtag(r.URL, response, user, params)
	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", e))

	if params.Json {
		err := RenderJson(w, r, response)
		if err != nil {
			return err
		}
	} else {
		if c.Dev {
			err := renderLocalTemplate(w, r, response, layout)
			if err != nil {
				return fmt.Errorf("error rendering local template: %v", err)
			}
		} else {
			err := renderEmbeddedTemplate(w, r, response, layout)
			if err != nil {
				return fmt.Errorf("error rendering embedded template: %v", err)
			}
		}
	}

	return nil

}

// RenderJson sends raw JSON to the request.
func RenderJson(w http.ResponseWriter, r *http.Request, response interface{}) error {
	params := r.Context().Value(ParamsKey{}).(FilterParams)
	var user UserKey
	if r.Context().Value(UserKey{}) != nil {
		user = r.Context().Value(UserKey{}).(UserKey)
	}

	e := setEtag(r.URL, response, user, params)
	w.Header().Set("Etag", fmt.Sprintf("\"%s\"", e))
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error rendering json response: %v", err)
	}
	_, err = w.Write(json)
	if err != nil {
		return fmt.Errorf("error writing json response: %v", err)

	}

	return nil
}

// renderEmbeddedTemplate loads embedded templates for production.
func renderEmbeddedTemplate(w http.ResponseWriter, r *http.Request, response interface{}, layout string) error {

	t, err := template.New("base").Funcs(funcs).Funcs(sprig.FuncMap()).ParseFS(templates.TemplatesDir, dirsEmbed(templates.TemplatesDir, "layouts/"+layout+".html", []string{"_default", "partials"})...)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.ExecuteTemplate(w, "base", response)
	if err != nil {
		return fmt.Errorf("error rendering embedded template: %v", err)
	}

	return nil
}

// renderLocalTemplate loads templates from the host to load template changes on browser reload during local development.
func renderLocalTemplate(w http.ResponseWriter, r *http.Request, response interface{}, layout string) error {
	t, err := template.New("base").Funcs(funcs).Funcs(sprig.FuncMap()).ParseFiles(dirsLocal("layouts/"+layout+".html", "_default", "partials")...)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.ExecuteTemplate(w, "base", response)
	if err != nil {
		return fmt.Errorf("error rendering local template: %v", err)
	}

	return nil
}

// dirsLocal returns templates from the host.
func dirsLocal(layout string, dirs ...string) []string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting working directory:", err)
	}

	var d []string

	for _, n := range dirs {

		path, err := filepath.Abs("pkg/templates/" + n)
		if err != nil {
			fmt.Println("error getting template dir:", err)
		}

		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("error reading template dir:", err)
		}

		for _, f := range files {
			d = append(d, fmt.Sprintf("%s/%s", path, f.Name()))
		}

	}

	d = append(d, fmt.Sprintf("%s/%s/%s", path, "pkg/templates", layout))

	return d
}

// dirsEmbed returns the embedded templates from the build binary.
func dirsEmbed(efs embed.FS, layout string, dirs []string) []string {
	var d []string

	f, err := getAllEmbedFilenames(&efs)
	if err != nil {
		fmt.Println("error getting all embedded filenames:", err)
	}
	for _, r := range f {

		for _, dir := range dirs {
			if strings.Contains(r, dir) {
				d = append(d, r)
			}
		}

	}

	d = append(d, layout)

	return d
}

func getAllEmbedFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
