package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/davecourtois/Utility"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/title/title_client"
	"github.com/globulario/services/golang/title/titlepb"
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

// get the thumbnail fil with help of youtube dl...
func downloadThumbnail(video_id, video_url, video_path string) (string, error){

	if !strings.Contains(video_path, ".mp4"){
		return "", errors.New("wrong file extension must be .mp4 video")
	}

	// The hidden folder path...
	path_ := video_path[0:strings.LastIndex(video_path, "/")]
	name_ := video_path[strings.LastIndex(video_path, "/")+1 : strings.LastIndex(video_path, ".")]
	thumbnail_path := path_ + "/.hidden/" + name_

	Utility.CreateDirIfNotExist(thumbnail_path)

	fmt.Println("try to run youtube-dl ",  video_url, thumbnail_path)
	cmd := exec.Command("youtube-dl", video_url, "-o", video_id,  "--write-thumbnail", "--skip-download")
	cmd.Dir = thumbnail_path
	
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	
	ext := "jpg"
	f, err := os.Open(thumbnail_path + "/" + video_id + "." + ext)
	if err != nil {
		ext = "webp"
		f, err = os.Open(thumbnail_path + "/" + video_id + "." + ext)
		if err != nil {
			return "", err
		}
	}

	defer f.Close()
    
    // Read entire JPG into byte slice.
    reader := bufio.NewReader(f)
    content, _ := ioutil.ReadAll(reader)
    
    // Encode as base64.
    encoded := base64.StdEncoding.EncodeToString(content)

	// cointain the data url...
	return "data:image/" + ext + ";base64," + encoded, nil
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
		} else {
			fmt.Println("fail to associate file with video information ", err)
		}

	})

	movieCollector.Visit(video_url)

	return nil
}

// Upload a video from porn hub and index it in the seach engine on the server side.
func indexXvideosVideo(token, video_url, index_path, video_path, file_path string) error {

	currentVideo := new(titlepb.Video)
	currentVideo.Casting = make([]*titlepb.Person, 0)
	currentVideo.Genres = []string{"adult"}
	currentVideo.Tags = []string{} // keep empty...
	currentVideo.URL = video_url
	currentVideo.ID = strings.Split(video_url, "/")[3]

	currentVideo.Poster = new(titlepb.Poster)
	currentVideo.Poster.ID = currentVideo.ID + "-thumnail"
	currentVideo.Poster.ContentUrl, _ = downloadThumbnail(currentVideo.ID, video_url, file_path) //e.Attr("src")
	currentVideo.Poster.URL = video_url
	currentVideo.Poster.TitleId = currentVideo.ID

	movieCollector := colly.NewCollector(
		colly.AllowedDomains("www.xvideos.com", "xvideos.com"),
	)

	// function call on visition url...
	movieCollector.OnHTML(".page-title", func(e *colly.HTMLElement) {

		currentVideo.Description = strings.TrimSpace(e.Text)
		e.ForEach(".duration", func(index int, child *colly.HTMLElement) {
			// The poster
			currentVideo.Duration = child.Text
			currentVideo.Description = currentVideo.Description[0:strings.Index(currentVideo.Description, currentVideo.Duration)]
			currentVideo.Description = strings.TrimSpace(currentVideo.Description)
		})

		e.ForEach(".video-hd-mark", func(index int, child *colly.HTMLElement) {
			// The poster
			tag := strings.TrimSpace(child.Text)
			if len(tag) > 3 {
				currentVideo.Tags = append(currentVideo.Tags, tag)
			}
		})
	})

	movieCollector.OnHTML(".label.profile", func(e *colly.HTMLElement) {
		p := new(titlepb.Person)
		p.URL = "https://www.xvideos.com" + e.Attr("href")
		e.ForEach(".name", func(index int, child *colly.HTMLElement) {
			// The poster
			p.ID = child.Text
			p.FullName = child.Text
		})
		if len(p.ID) > 0 {
			currentVideo.Casting = append(currentVideo.Casting, p)
		}
	})

	movieCollector.OnHTML(".uploader-tag ", func(e *colly.HTMLElement) {
		currentVideo.PublisherId = new(titlepb.Publisher)
		currentVideo.PublisherId.URL = e.Attr("href")

		e.ForEach(".name", func(index int, child *colly.HTMLElement) {
			// The poster
			currentVideo.PublisherId.ID = child.Text
			currentVideo.PublisherId.Name = child.Text

		})

	})

	// The number of view
	movieCollector.OnHTML("#v-actions-left > div.vote-actions > div.rate-infos > span", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		str := strings.Split(e.Text, " ")[0]
		str = strings.ReplaceAll(str, ",", "")
		currentVideo.Count = int64(Utility.ToInt(str))
	})

	// The raiting
	movieCollector.OnHTML("#v-actions-left > div.vote-actions > div.rate-infos > div.perc > div.good > span.rating-good-perc", func(e *colly.HTMLElement) {
		percent := e.Text
		percent = strings.ReplaceAll(percent, "%", "")
		currentVideo.Rating = float32(Utility.ToNumeric(percent) / 10)
	})
	
	// The tags
	movieCollector.OnHTML("#main > div.video-metadata.video-tags-list.ordered-label-list.cropped > ul > li > a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/tags/") {
			currentVideo.Tags = append(currentVideo.Tags, strings.TrimSpace(e.Text))
		}
	})

	movieCollector.OnHTML("#footer", func(e *colly.HTMLElement) {

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
				return
			}
		} else {
			fmt.Println("fail to associate file with video information ", err)
			return
		}
	})

	movieCollector.Visit(video_url)
	return nil
}


////////////////////////////////// Youtube ////////////////////////////////////////////////

func indexYoutubeVideo(token, video_url, index_path, video_path, file_path string) error {

	currentVideo := new(titlepb.Video)
	currentVideo.Casting = make([]*titlepb.Person, 0)
	currentVideo.Genres = []string{"youtube"}
	currentVideo.Tags = []string{} // keep empty...
	currentVideo.URL = video_url
	currentVideo.ID = video_url[strings.LastIndex(video_url, "=") + 1 :]

	currentVideo.Poster = new(titlepb.Poster)
	currentVideo.Poster.ID = currentVideo.ID + "-thumnail"
	currentVideo.Poster.ContentUrl, _ = downloadThumbnail(currentVideo.ID, video_url, file_path) //e.Attr("src")
	currentVideo.Poster.URL = video_url
	currentVideo.Poster.TitleId = currentVideo.ID

	// For that one I will made use of web-api from https://noembed.com/embed
	url := `https://noembed.com/embed?url=` + video_url
	var myClient = &http.Client{Timeout: 5 * time.Second}
	r, err := myClient.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

	target := make(map[string]interface{})
    json.NewDecoder(r.Body).Decode(&target)

	// console.log()
	currentVideo.PublisherId = new(titlepb.Publisher)
	currentVideo.PublisherId.URL = target["author_url"].(string)
	currentVideo.PublisherId.ID = target["author_name"].(string) 
	currentVideo.PublisherId.Name = target["author_name"].(string)
	currentVideo.Description =  target["title"].(string) 

	title_client_, err := getTitleClient()
	if err != nil {
		return err
	}

	if currentVideo == nil {
		return err
	}

	err = title_client_.CreateVideo(token, index_path, currentVideo)
	if err == nil {
		err := title_client_.AssociateFileWithTitle(index_path, currentVideo.ID, video_path)
		if err != nil {
			fmt.Println("fail to associate file with video information ", err)
			return err
		}
	} else {
		fmt.Println("fail to associate file with video information ", err)
		return err
	}


	return nil
}
