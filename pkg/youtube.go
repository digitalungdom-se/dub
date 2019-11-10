package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type (
	VideoInfo struct {
		Title        string
		Author       *Author
		Duration     string
		ThumbnailURL string
	}

	Author struct {
		Name      string
		URL       string
		AvatarURL string
	}

	oembedResponse struct {
		Author    string `json:"author_name"`
		AuthorURL string `json:"author_url"`
		Title     string `json:"title"`
	}
)

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func GetYoutubeInfo(url string) (*VideoInfo, error) {
	videoInfo := new(VideoInfo)
	videoInfo.Author = new(Author)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("img").Each(func(index int, item *goquery.Selection) {
		imgTag := item
		_, found := imgTag.Attr("data-ytimg")
		if found {
			src, _ := imgTag.Attr("data-thumb")
			isProfilePic, _ := imgTag.Attr("alt")

			if isProfilePic != "" {
				videoInfo.Author.AvatarURL = src
			} else {
				videoInfo.ThumbnailURL = src
			}
		}
	})

	oembed := new(oembedResponse)
	getJson(fmt.Sprintf("https://www.youtube.com/oembed?url=%v&format=json", url), oembed)

	videoInfo.Title = oembed.Title
	videoInfo.Author.Name = oembed.Author
	videoInfo.Author.URL = oembed.AuthorURL

	return videoInfo, nil
}
