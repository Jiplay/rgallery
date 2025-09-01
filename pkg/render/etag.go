package render

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/robbymilo/rgallery/pkg/middleware"
)

func setEtag(url *url.URL, response interface{}, user UserKey, params FilterParams) string {
	value := GenerateEtag(fmt.Sprint(response) + fmt.Sprint(user) + fmt.Sprint(params.Json))
	key := fmt.Sprint(url) + fmt.Sprint(user) + fmt.Sprint(params.Json)
	middleware.PersistEtag(key, value)
	return value

}

func GenerateEtag(body string) string {
	hash := md5.Sum([]byte(body))
	return hex.EncodeToString(hash[:])
}
