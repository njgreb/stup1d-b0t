package commands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var tenorKey string

func getKey() string {
	if len(strings.TrimSpace(tenorKey)) == 0 {
		tenorKey = os.Getenv("tenorKey")
	}

	return tenorKey
}

func Gif(searchTerm string) string {

	tenorUrl := "https://g.tenor.com/v1/search?q=" + searchTerm + "&key=" + getKey() + "&limit=1&contentfilter=off&media_filter=minimal&locale=en_US"
	res, err := http.Get(tenorUrl)
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// derp
	}

	var gif_result results_set
	json.Unmarshal(body, &gif_result)

	return gif_result.Results[0].Media[0].MediumGif.GifUrl + "\nPowered by Tenor yo"
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
