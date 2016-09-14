package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/ffjson/ffjson"
)

const (
	URLGET                        string = "http://api.brightcove.com/services/library"
	METHOD_GET_FindVideosByTags   string = "find_videos_by_tags"
	METHOD_GET_FindModifiedVideos string = "find_modified_videos"
)

type GetRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type ItemsCollection struct {
	Items      []Video `json:"items"`
	PageNumber int     `json:"page_number"`
	PageSize   int     `json:"page_size"`
	TotalCount int     `json:"total_count"`
}

type Video struct {
	Id                int         `json:"id"`
	Name              string      `json:"name"`
	AdKeys            string      `json:"adKeys"`
	ShortDescription  string      `json:"shortDescription"`
	LongDescription   string      `json:"longDescription"`
	CreationDate      string      `json:"creationDate"`
	PublishedDate     string      `json:"publishedDate"`
	LastModifiedDate  string      `json:"lastModifiedDate"`
	LinkURL           string      `json:"linkURL"`
	LinkText          string      `json:"linkText"`
	Tags              []string    `json:"tags"`
	VideoStillURL     string      `json:"videoStillURL"`
	ThumbnailURL      string      `json:"thumbnailURL"`
	ReferenceId       string      `json:"referenceId"`
	Length            int         `json:"length"`
	Economics         string      `json:"economics"`
	PlaysTotal        string      `json:"playsTotal"`
	PlaysTrailingWeek string      `json:"playsTrailingWeek"`
	FLVURL            string      `json:"FLVURL"`
	Renditions        []Rendition `json:"renditions"`
	FLVFullLength     Rendition   `json:"FLVFullLength"`
	VideoFullLength   Rendition   `json:"videoFullLength"`
}

type Rendition struct {
	AudioOnly             bool
	ControllerType        string
	displayName           string
	EncodingRate          int
	FrameHeight           int
	FrameWidth            int
	Id                    int
	ReferenceId           string
	RemoteStreamName      string
	RemoteUrl             string
	Size                  int
	UploadTimestampMillis int
	Url                   string
	VideoCodec            string
	VideoContainer        string
	VideoDuration         int
}

type FindVideosByTagsStruct struct {
	Token          string
	and_tags       interface{}
	or_tags        interface{}
	page_size      interface{}
	page_number    interface{}
	sort_by        interface{}
	sort_order     interface{}
	get_item_count interface{}
	video_fields   interface{}
	custom_fields  interface{}
	media_delivery interface{}
	output         interface{}
}

func FindVideosByTags(or_tags []string, and_tags []string) (videos []Video) {
	page := 0

	for {
		v := url.Values{}
		v.Add("token", config.Token)
		v.Add("command", METHOD_GET_FindVideosByTags)
		v.Add("page_number", strconv.Itoa(page))
		v.Add("page_size", "25")
		if or_tags != nil {
			v.Add("or_tags", strings.Join(or_tags, ","))
		}
		if and_tags != nil {
			v.Add("and_tags", strings.Join(or_tags, ","))
		}
		v.Add("video_fields", "id")
		getUrl, err := url.Parse(URLGET + "?" + v.Encode())

		if err != nil {
			log.Printf("%+v\n", err)
		}

		items, err := call(getUrl)
		if err == nil {
			count := len(items)
			log.Printf("Items in result %+v\n", count)

			for _, item := range items {
				videos = append(videos, item)
			}
			if count == 0 || count < 25 {
				break
			} else {
				page++
			}
		}
	}

	return
}

func FindModifiedVideos(fromdate time.Time) (videos []Video) {
	log.Printf("%+v\n", fromdate.String())

	page := 0

	for {
		v := url.Values{}
		v.Add("token", config.Token)
		v.Add("command", METHOD_GET_FindModifiedVideos)
		v.Add("page_number", strconv.Itoa(page))
		v.Add("page_size", "25")
		v.Add("from_date", strconv.Itoa(int(fromdate.Unix()/60)))
		v.Add("video_fields", "id,tags")
		getUrl, err := url.Parse(URLGET + "?" + v.Encode())

		if err != nil {
			log.Printf("%+v\n", err)
		}

		items, err := call(getUrl)
		if err == nil {
			count := len(items)
			log.Printf("Items in result %+v\n", count)

			for _, item := range items {
				videos = append(videos, item)
			}
			if count == 0 || count < 25 {
				break
			} else {
				page++
			}
		}
	}

	return
}

func call(url *url.URL) ([]Video, error) {
	log.Println(url.String())

	resp, err := http.Get(url.String())

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("err %+v\n", err)

		return []Video{}, err
	} else {
		var result ItemsCollection
		err := ffjson.Unmarshal(body, &result)

		if err != nil {
			log.Printf("err %+v\n", err)

			return []Video{}, err
		}

		return result.Items, err
	}
}

func FindVideosWithTags(videos []Video, tags []string) (videoIds []int) {
	for _, item := range videos {
		found := false
		for _, tag := range tags {
			if stringInSlice(tag, item.Tags) {
				found = true
				break
			}
		}

		if found {
			videoIds = append(videoIds, item.Id)
			log.Printf("body %+v\n", item.Id)

			log.Printf("%+v\n", item.Tags)
		}
	}

	return
}
