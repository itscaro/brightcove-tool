package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pquerna/ffjson/ffjson"
)

const (
	URLPOST                 string = "http://qapi.brightcove.com/services/post"
	METHOD_POST_SHARE_VIDEO string = "share_video"
)

type PostRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type PostResponse struct {
	Id     interface{} `json:"id"`
	Error  interface{} `json:"error"`
	Result interface{} `json:"result"`
}

type ShareVideoParams struct {
	VideoId          int    `json:"video_id"`
	ShareeAccountIds []int  `json:"sharee_account_ids"`
	AutoAccept       bool   `json:"auto_accept "`
	ForceReshare     bool   `json:"force_reshare"`
	Token            string `json:"token"`
}

func ShareVideo(shareeAccountIds []int, videoIds []int, autoAccept bool, forceReshare bool) {
	for _, videoId := range videoIds {
		request := PostRequest{
			Method: METHOD_POST_SHARE_VIDEO,
			Params: ShareVideoParams{
				VideoId:          videoId,
				AutoAccept:       autoAccept,
				ForceReshare:     forceReshare,
				ShareeAccountIds: shareeAccountIds,
				Token:            config.Token,
			},
		}
		payload, err := ffjson.Marshal(request)

		log.Printf("payload %+v\n", string(payload))
		if err != nil {
			log.Printf("err %+v\n", err)
		} else {
			resp, err := http.PostForm(URLPOST, url.Values{"json": {string(payload)}})

			if err != nil {
				log.Printf("t l %+v\n", err)
			} else {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					log.Printf("err %+v\n", err)
				} else {
					var result PostResponse
					err = ffjson.Unmarshal(body, &result)

					if err != nil {
						log.Printf("err %+v\n", err)
					} else {
						log.Printf("body %+v\n", result)
					}
				}
			}
		}
	}
}
