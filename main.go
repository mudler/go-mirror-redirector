package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"

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
	config := map[string][]string{}
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("yamlFile.Get err   #%v ", err))
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal err   #%v ", err))
	}

	m := macaron.Classic()

	m.Get("/", func(ctx *macaron.Context, req *http.Request, w http.ResponseWriter) {
		remoteIP := ctx.RemoteAddr()
		geo, err := geoip.New()
		if err != nil {
			ctx.Resp.WriteHeader(500)
			return
		}
		country := geo.Lookup(net.ParseIP(remoteIP))

		if mirrors, ok := config[country.Short]; ok {
			randomIndex := rand.Intn(len(mirrors))
			pick := mirrors[randomIndex]
			fmt.Println("Redirecting", remoteIP, country, "to", pick)
			ctx.Redirect(pick, 301)
			return
		}

		mirrors, ok := config["default"]
		if len(mirrors) == 0 || !ok {
			ctx.HTML(200, "Warning", "No mirrors configured in the 'default' key")
			return
		}

		randomIndex := rand.Intn(len(mirrors))
		pick := mirrors[randomIndex]
		fmt.Println("Redirecting", remoteIP, country, "to", pick)
		ctx.Redirect(pick, 301)

	})
	m.Run()
}
