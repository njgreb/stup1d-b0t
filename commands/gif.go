package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var tenorKey string

func getKey() string {
	if len(strings.TrimSpace(tenorKey)) == 0 {
		tenorKey = os.Getenv("tenorKey")
	}

	return tenorKey
}

func Gif(searchTerm string, imageOnly bool, contentfilter string) string {

	var returnString string

	if len(searchTerm) == 0 {
		// Just get a trending gif
		tenorUrl := "https://g.tenor.com/v1/trending?key=" + getKey() + "&media_filter=minimal&contentfilter=" + contentfilter + "locale=en_US"
		spew.Dump(tenorUrl)
		res, err := http.Get(tenorUrl)
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			// derp
		}

		var gif_result results_set
		json.Unmarshal(body, &gif_result)

		rand.Seed(time.Now().UnixNano())

		spew.Dump(gif_result)

		randomIndex := rand.Intn(len(gif_result.Results))
		fmt.Printf("picking image %d out of %d", randomIndex, len(gif_result.Results))
		returnString = gif_result.Results[randomIndex].Media[0].MediumGif.GifUrl

	} else {
		// Find the perfect gif for the term
		tenorUrl := "https://g.tenor.com/v1/search?q=" + url.QueryEscape(searchTerm) + "&key=" + getKey() + "&limit=20&contentfilter=off&media_filter=minimal&locale=en_US"
		res, err := http.Get(tenorUrl)
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			// derp
		}

		var gif_result results_set
		json.Unmarshal(body, &gif_result)

		randomIndex := rand.Intn(len(gif_result.Results))
		fmt.Printf("picking image %d out of %d", randomIndex, len(gif_result.Results))
		returnString = gif_result.Results[randomIndex].Media[0].MediumGif.GifUrl
	}

	if imageOnly == false {
		returnString += "\nPowered by Tenor yo"
	}
	return returnString

}

// Tenor results struct
type tenor_result struct {
	Id         string        `json:"id,omitempty"`
	Title      string        `json:"title,omitempty"`
	H1Title    string        `json:"h1_title,omitempty"`
	Media      []tenor_media `json:"media,omitempty"`
	BgColor    string        `json:"bg_color,omitempty"`
	Created    string        `json:"created,omitempty"`
	Itemurl    string        `json:"itemurl,omitempty"`
	TenorUrl   string        `json:"url,omitempty"`
	Shares     int           `json:"shares,omitempty"`
	Hasaudio   bool          `json:"hasaudio,omitempty"`
	Hascaption bool          `json:"hascaption,omitempty"`
}

type results_set struct {
	Results []tenor_result `json:"results,omitempty"`
}

type tenor_media struct {
	MediumGif medium_gif `json:"gif,omitempty"`
}

type medium_gif struct {
	Preview string `json:"preview,omitempty"`
	GifUrl  string `json:"url,omitempty"`
}
