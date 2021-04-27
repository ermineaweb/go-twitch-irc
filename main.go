package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jrm780/gotirc"
)

var MAX_RANK, _ = strconv.Atoi(GetEnv("NB_VIEWER_DISPLAY"))
var IRC_ADDRESS = GetEnv("IRC_ADDRESS")
var IRC_PORT, _ = strconv.Atoi(GetEnv("IRC_PORT"))
var CLIENT_USERNAME = GetEnv("CLIENT_USERNAME")
var CLIENT_AUTH_SECRET = GetEnv("CLIENT_AUTH_SECRET")
var STREAMERS = strings.Split(GetEnv("STREAMERS"), ",")

func main() {
	streamers := createStreamers(STREAMERS)
	listenStreamerChat(streamers)
}

func listenStreamerChat(streamers []Streamer) {
	var channels []string
	for _, streamer := range streamers {
		channels = append(channels, "#"+streamer.Username)
	}

	options := gotirc.Options{Host: IRC_ADDRESS, Port: IRC_PORT, Channels: channels}

	client := gotirc.NewClient(options)

	client.OnChat(func(channel string, tags map[string]string, msg string) {
		streamer := getStreamerByChannel(channel, streamers)
		saveMessage(streamer, tags["display-name"], msg)
		displayResults(streamers)
	})

	// Connect and authenticate with the given Username and oauth token
	client.Connect(CLIENT_USERNAME, CLIENT_AUTH_SECRET)
}

func getStreamerByChannel(channel string, streamers []Streamer) Streamer {
	for _, s := range streamers {
		if "#"+s.Username == channel {
			return s
		}
	}
	return Streamer{}
}

func saveMessage(streamer Streamer, Username string, msg string) {
	// add the new message to the viewer's messages
	messages := append(streamer.Viewers[Username].Messages, msg)
	// upate the count and the messages list
	viewer := Viewer{Username: Username, Messages: messages, MessagesCount: len(messages)}
	// update the viewer
	streamer.Viewers[viewer.Username] = viewer
}

func getPodium(streamer Streamer) []Viewer {
	var ViewersResult []Viewer
	for _, viewer := range streamer.Viewers {
		ViewersResult = append(ViewersResult, viewer)
	}
	sort.Sort(ByMessagesCount(ViewersResult))
	return ViewersResult
}

var nbRequest = 0

func displayResults(streamers []Streamer) {
	nbRequest++
	// clear screen
	fmt.Print("\033[H\033[2J")
	fmt.Printf("Requests: %v\nTime: %v\n", nbRequest, time.Now().String())
	fmt.Println("---------------------------------------------")

	for _, streamer := range streamers {
		fmt.Printf("Viewers messages count for streamer [%v]\n---------------------------------------------\n", streamer.Username)

		ViewersResult := getPodium(streamer)

		for key, viewer := range ViewersResult {
			if key < MAX_RANK {
				fmt.Printf("%v - [%v]: %d\n", key+1, viewer.Username, viewer.MessagesCount)
			}
		}
		fmt.Println("---------------------------------------------")
	}
}

func createStreamers(Usernames []string) []Streamer {
	streamers := []Streamer{}
	for _, Username := range Usernames {
		streamers = append(streamers, Streamer{Username: Username, Viewers: make(map[string]Viewer)})
	}
	return streamers
}
