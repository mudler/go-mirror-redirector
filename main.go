package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/alecthomas/geoip"
	"gopkg.in/macaron.v1"
	"gopkg.in/yaml.v2"
)

// Take a yaml of the following format
// countrycode:
// - mirror1
// - mirror2
// default:
// - fallback1
// - fallback2
// And use it to redirect people when hitting our port

func main() {

	configFile := os.Getenv("CONFIG")
	if len(configFile) == 0 {
		configFile = "config.yaml"
	}
	var config = map[string][]string{}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("yamlFile.Get err   #%v ", err))
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal err   #%v ", err))
	}

	countryToMirror := func(urlpath string, mirrors []string) string {
		// select a random mirror
		randomIndex := rand.Intn(len(mirrors))
		pick := mirrors[randomIndex]

		// Join the original path requested with the mirror, so we redirect to the specific page
		u, err := url.Parse(pick)
		if err != nil {
			panic(err)
		}
		u.Path = path.Join(u.Path, urlpath)
		pick = u.String()
		return pick
	}

	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Get("/*", func(ctx *macaron.Context, req *http.Request, w http.ResponseWriter) {
		urlpath := req.URL.Path

		remoteIP := ctx.RemoteAddr()
		geo, err := geoip.New()
		if err != nil {
			ctx.Resp.WriteHeader(500)
			return
		}
		country := geo.Lookup(net.ParseIP(remoteIP))
		if country != nil {
			if mirrors, ok := config[country.Short]; ok {
				pick := countryToMirror(urlpath, mirrors)

				fmt.Println("Redirecting", remoteIP, country, "to", pick)
				ctx.Redirect(pick, 301)
				return
			}
		}

		mirrors, ok := config["default"]
		if len(mirrors) == 0 || !ok {
			ctx.PlainText(200, []byte("No mirrors configured in the 'default' key"))
			return
		}

		pick := countryToMirror(urlpath, mirrors)
		fmt.Println("Redirecting", remoteIP, country, "to", pick)

		ctx.Redirect(pick, 301)
	})

	m.Run()
}
