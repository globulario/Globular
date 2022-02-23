package main

import (
	"fmt"
	"strings"

	"github.com/davecourtois/Utility"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/title/title_client"
	"github.com/globulario/services/golang/title/titlepb"
	"github.com/go-acme/lego/v4/log"
	colly "github.com/gocolly/colly/v2"
)

var (
	title_client_ *title_client.Title_Client
)

// return the actual title client.
func getTitleClient() (*title_client.Title_Client, error) {
	if title_client_ == nil {
		address, err := config.GetAddress()
		if err != nil {
			return nil, err
		}
		title_client_, err = title_client.NewTitleService_Client(address, "title.TitleService")
		if err != nil {
			return nil, err
		}
	}
	return title_client_, nil
}


// Upload a video from porn hub and index it in the seach engine on the server side.
func indexPornhubVideo(token, video_url, index_path, video_path string) error {

	currentVideo := new(titlepb.Video)
	currentVideo.Casting = make([]*titlepb.Person, 0)
	currentVideo.Genres = []string{"adult"}
	currentVideo.Tags = []string{} // keep empty...
	currentVideo.URL = video_url
	currentVideo.ID = strings.Split(video_url, "=")[1]

	movieCollector := colly.NewCollector(
		colly.AllowedDomains("pornhub.com", "www.pornhub.com"),
	)

	/////////////////////////////////////////////////////////////////////////
	// Movie info
	/////////////////////////////////////////////////////////////////////////

	// function call on visition url...
	movieCollector.OnHTML(".inlineFree", func(e *colly.HTMLElement) {
		currentVideo.Description = strings.TrimSpace(e.Text)
	})

	movieCollector.OnHTML(".pstar-list-btn", func(e *colly.HTMLElement) {
		p := new(titlepb.Person)
		p.ID = e.Attr("data-id")
		p.FullName = strings.TrimSpace(e.Text)
		p.URL = "https://www.pornhub.com" + e.Attr("href")
		if len(p.ID) > 0 {
			currentVideo.Casting = append(currentVideo.Casting, p)
		}
	})

	movieCollector.OnHTML("#hd-leftColVideoPage > div:nth-child(1) > div.video-actions-container > div.video-actions-tabs > div.video-action-tab.about-tab.active > div.video-detailed-info > div.video-info-row.userRow > div.userInfo > div > a", func(e *colly.HTMLElement) {
		currentVideo.PublisherId = new(titlepb.Publisher)
		currentVideo.PublisherId.ID = e.Text
		currentVideo.PublisherId.Name = e.Text
		currentVideo.PublisherId.URL = e.Attr("href")

	})

	movieCollector.OnHTML(".count", func(e *colly.HTMLElement) {
		count := e.Text
		if strings.Contains(count, "K") {
			count = strings.ReplaceAll(count, "K", "")
			currentVideo.Count = int64(Utility.ToInt(e.Text)) * 100000
		} else if strings.Contains(count, "M") {
			count = strings.ReplaceAll(count, "M", "")
			currentVideo.Count = int64(Utility.ToInt(e.Text)) * 1000000
		}

	})

	movieCollector.OnHTML(".percent", func(e *colly.HTMLElement) {
		percent := e.Text
		percent = strings.ReplaceAll(percent, "%", "")
		currentVideo.Rating = float32(Utility.ToNumeric(percent) / 10)
	})

	movieCollector.OnHTML(".videoElementPoster", func(e *colly.HTMLElement) {
		// The poster
		currentVideo.Poster = new(titlepb.Poster)
		currentVideo.Poster.ID = currentVideo.ID + "-thumnail"
		currentVideo.Poster.ContentUrl = e.Attr("src")
		currentVideo.Poster.URL = video_url
		currentVideo.Poster.TitleId = currentVideo.ID
	})

	movieCollector.OnHTML(".categoriesWrapper a", func(e *colly.HTMLElement) {
		tag := strings.TrimSpace(e.Text)
		if tag != "Suggest" {
			currentVideo.Tags = append(currentVideo.Tags, tag)
		}
	})

	movieCollector.OnHTML(".show-more-btn", func(e *colly.HTMLElement) {

		title_client_, err := getTitleClient()
		if err != nil {
			return
		}

		if currentVideo == nil {
			return
		}

		err = title_client_.CreateVideo(token, index_path, currentVideo)
		if err == nil {
			err := title_client_.AssociateFileWithTitle(index_path, currentVideo.ID, video_path)
			if err != nil {
				fmt.Println("fail to associate file with video information ", err)
			}
		}else{
			fmt.Println("fail to associate file with video information ", err)
		}

	})

	movieCollector.Visit(video_url)

	return nil
}
