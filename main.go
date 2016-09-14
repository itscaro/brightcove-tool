package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var config Config

func main() {
	//log.SetFlags(log.Lshortfile)
	log.Println("Starting brightcove")

	data, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("Cannot find ./config.yml")
		return
	} else {
		err := yaml.Unmarshal([]byte(data), &config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Printf("Config:\n%+v\n\n", config)
	}

	var wg sync.WaitGroup

	shareTickerChan := time.NewTicker(time.Second * 5).C
	importTickerChan := time.NewTicker(time.Minute * 15).C
	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {

		case <-shareTickerChan:
			log.Println("Check for new videos to share")
			for _, shareConfig := range config.Share {
				wg.Add(1)
				go func() {
					defer wg.Done()
					defer recoverFunc()

					log.Printf("Config: %+v\n", shareConfig)

					videos := FindModifiedVideos(time.Now().Truncate(time.Duration(time.Hour) * 100))
					videoIds := FindVideosWithTags(videos, shareConfig.Tags)

					log.Printf("%+v\n", videoIds)

					//ShareVideo(shareConfig.ShareeAccountIds, videoIds, true, true)

				}()
			}
			log.Println("--- DONE SHARE ---")

		case <-importTickerChan:
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer recoverFunc()

				log.Println("Check for new videos to import")

				time.Sleep(time.Duration(time.Second) * 5)

				log.Println("--- DONE IMPORT ---")
			}()

		case <-sigChan:
			log.Println("Going to quit, after all jobs done")
			wg.Wait()
			log.Println("Quitting...")
			return

		}
	}
}

func recoverFunc() {
	if r := recover(); r != nil {
		log.Println("Recovered in f", r)
	}
}
