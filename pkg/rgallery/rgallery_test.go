package rgallery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/robbymilo/rgallery/pkg/database"
	"github.com/robbymilo/rgallery/pkg/scanner"
	"github.com/robbymilo/rgallery/pkg/types"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

type ResponseMedia = types.ResponseMedia
type ResponseTags = types.ResponseTags
type ResponseFolders = types.ResponseFolders
type ResponseFilter = types.ResponseFilter
type ResponseGear = types.ResponseGear

var c = Conf{
	Dev:              false,
	DisableAuth:      true,
	Media:            "../../testdata/media",
	Cache:            "../../testdata/cache",
	Data:             "../../testdata/data",
	PreGenerateThumb: true,
	ResizeService:    "",
	LocationDataset:  "Provinces10",
	Logger:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
	OnThisDay:        false,
}

var ca = cache.New(-1, -1)

func TestEndpoints(t *testing.T) {
	aliases := make(map[string]string)
	aliases["Nikon 105mm f/2.5 Ai-s"] = "Nikon Ai-s 105mm f/2.5"
	aliases["Nikon 105mm f/2.5 AI-s"] = "Nikon Ai-s 105mm f/2.5"
	aliases["Nikon AI-s 105mm f/2.5"] = "Nikon Ai-s 105mm f/2.5"
	aliases["Nikon AI-s 105mm f/2.5   "] = "Nikon Ai-s 105mm f/2.5"
	aliases["Nikon 105mm f/2.5 Ai-s "] = "Nikon Ai-s 105mm f/2.5"
	aliases["AF-S Nikkor 50mm f/1.8G"] = "123"
	c.Aliases.Lenses = aliases

	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		fmt.Println("error getting timezone", err)
	}
	time.Local = loc

	testEndpoints := []string{
		"../../testdata/media/2019/20190330-sawtooths/20190330-copper-mtn-robbymilo-1112.jpg",
		"../../testdata/media/2017/20170624-idaho/20170624-sawtooth-mountain-biking-robbymilo-0030.jpg",
		"../../testdata/media/2017/20170624-idaho/20170624-sawtooth-mountain-biking-robbymilo-0030 copy.jpg",
		"../../testdata/media/2015/boise/20150111rmilo-0775.jpg",
		"../../testdata/media/2015/boise/20150111rmilo-0776.jpg",
		"../../testdata/media/2016/boise/20160424-boise-robbymilo-0123.jpg",
		"../../testdata/media/2018/20180304-bogus/20180304-mores-mountain-ski-robbymilo-0230.jpg",
		"../../testdata/media/2018/20180304-bogus/20180304-mores-mountain-ski-robbymilo-0230 2.jpg",
		"../../testdata/media/2019/20190518-bogus-basin/20190518-shafer-butte-robbymilo-0007.jpg",
		"../../testdata/media/2024/105-5.jpg",
		"../../testdata/media/misc/105-1.jpg",
		"../../testdata/media/misc/105-2.jpg",
		"../../testdata/media/misc/105-3.jpg",
		"../../testdata/media/misc/105-4.jpg",
		"../../testdata/media/misc/51750950528.jpg",
	}

	for _, i := range testEndpoints {
		setModTime(i)
	}

	database.CreateDB(c)

	_, _ = scanner.Scan("default", c, ca)

	testResponse(t, "/?format=json", "../../testdata/ResponseFilter.json")

	testResponse(t, "/?format=json", "../../testdata/ResponseFilter.json")
	testResponse(t, "/?format=json&orderby=modified", "../../testdata/ResponseFilter-modified.json")
	testResponse(t, "/?format=json&orderby=modified&direction=asc", "../../testdata/ResponseFilter-modified-asc.json")
	testResponse(t, "/?format=json&term=copp", "../../testdata/ResponseFilter-term.json")
	testResponse(t, "/?format=json&term=su+špa", "../../testdata/ResponseFilter-term-1.json")
	testResponse(t, "/?format=json&camera=NIKON%20D800", "../../testdata/ResponseFilter-camera.json")
	testResponse(t, "/?format=json&lens=AF-S%20Nikkor%2050mm%20f%2f1.8G", "../../testdata/ResponseFilter-lens.json")
	testResponse(t, "/?format=json&lens=123", "../../testdata/ResponseFilter-lens-1.json")
	testResponse(t, "/?format=json&lens=Nikon Ai-s 105mm f%2f2.5", "../../testdata/ResponseFilter-lens-2.json")
	testResponse(t, "/?format=json&folder=2017/20170624-idaho", "../../testdata/ResponseFilter-folder.json")
	testResponse(t, "/?format=json&tag=idaho", "../../testdata/ResponseFilter-tag.json")
	testResponse(t, "/onthisday?format=json", "../../testdata/ResponseFilter-onthisday.json")

	testResponse(t, "/media/651935749?format=json", "../../testdata/ResponseImage-0.json")
	testResponse(t, "/media/3455659031?format=json", "../../testdata/ResponseImage-1.json")
	testResponse(t, "/media/4119775194?format=json", "../../testdata/ResponseImage-2.json")
	testResponse(t, "/media/651935749/in/folder/2017/20170624-idaho?format=json", "../../testdata/ResponseImage-folder.json")
	testResponse(t, "/media/651935749/in/tag/idaho?format=json", "../../testdata/ResponseImage-tag.json")
	testResponse(t, "/media/4119775194/in/tag/%40acconfb?format=json", "../../testdata/ResponseImage-tag-acc.json")
	testResponse(t, "/media/4119775194/in/tag/%23californiawildfires?format=json", "../../testdata/ResponseImage-tag-cal.json")
	testResponse(t, "/media/651935749/in/favorites?format=json", "../../testdata/ResponseImage-favorites.json")

	// prev/next responses
	testResponse(t, "/media/3455659031?camera=NIKON D800&format=json", "../../testdata/ResponseImage-camera.json")
	testResponse(t, "/media/651935749?lens=AF-S Nikkor 50mm f%2f1.8G&format=json", "../../testdata/ResponseImage-lens.json")
	testResponse(t, "/media/651935749?lens=123&format=json", "../../testdata/ResponseImage-lens.json")
	testResponse(t, "/media/525791494?lens=Nikon Ai-s 105mm f%2f2.5&format=json", "../../testdata/ResponseImage-lens-1.json")
	testResponse(t, "/media/264898052?focallength35=50&format=json", "../../testdata/ResponseImage-focallength35.json")
	testResponse(t, "/media/3216513272?software=darktable 4.4.2&format=json", "../../testdata/ResponseImage-software.json")
	testResponse(t, "/media/1129346697?term=bogus&format=json", "../../testdata/ResponseImage-term.json")

	testResponse(t, "/favorites?format=json", "../../testdata/ResponseImages-favorites.json")

	testResponse(t, "/folders?format=json", "../../testdata/ResponseFolders.json")
	testResponse(t, "/folder/2019/20190330-sawtooths?format=json", "../../testdata/ResponseFolder.json")

	testResponse(t, "/tags?format=json", "../../testdata/ResponseTags.json")
	testResponse(t, "/tag/idaho?format=json", "../../testdata/ResponseTag.json")
	testResponse(t, "/tag/%40acconfb?format=json", "../../testdata/ResponseTag-acc.json")
	testResponse(t, "/tag/%23californiawildfires?format=json", "../../testdata/ResponseTag-cal.json")

	testResponse(t, "/gear?format=json", "../../testdata/ResponseGear.json")
	testResponse(t, "/map?format=json", "../../testdata/ResponseMap.json")

	testStatusCode(t, "/media/123", 404)
	testStatusCode(t, "/folder/123", 404)
	testStatusCode(t, "/tag/123", 404)

	// test thumbnail generation
	testThumbnail(t, "3455659031", "3000")
	testThumbnail(t, "651935749", "4000")
	testThumbnail(t, "651935749", "2400")
	testThumbnail(t, "3455659031", "2400")

}

func testResponse(t *testing.T, path string, json_relative_path string) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	router(w, r, ca)
	response_json, err := io.ReadAll(w.Body)
	if err != nil {
		fmt.Println(err)
	}

	json_file, err := os.Open(json_relative_path)
	if err != nil {
		fmt.Println(err)
	}

	defer json_file.Close()

	json_byte, err := io.ReadAll(json_file)
	if err != nil {
		fmt.Println(err)
	}

	dst := &bytes.Buffer{}
	if err := json.Compact(dst, []byte(json_byte)); err != nil {
		panic(err)
	}

	assert.EqualValues(t, dst.String(), string(response_json), "they should be equal")
}

func testStatusCode(t *testing.T, path string, statusCode int) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	router(w, r, ca)
	res := w.Result()
	defer res.Body.Close()

	assert.EqualValues(t, statusCode, res.StatusCode, "they should be equal")

}

func testThumbnail(t *testing.T, id, size string) {
	var exists bool
	path := "../../testdata/cache/" + size + "/" + id + ".jpg"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		exists = false
	} else {
		exists = true
	}

	fmt.Println("exists", exists, path)
	assert.EqualValues(t, exists, true, "they should be equal")

}

func setModTime(test_path string) {
	m, err := time.Parse("2006-01-02T15:04:05.000Z", "2023-11-21T20:44:53.923Z")
	if err != nil {
		fmt.Println(err)
	}

	err = os.Chtimes(test_path, m, m)
	if err != nil {
		fmt.Println(err)
	}

}

func router(w http.ResponseWriter, r *http.Request, cache *cache.Cache) {

	SetupRouter(c, ca, "", "").ServeHTTP(w, r)

}
