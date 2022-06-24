package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// GetMockStream is a blackbox function which returns a mock stream for
// demonstration purposes
func GetMockStream() Stream {
	return Stream{0, mockdata}
}

// Stream is a mock stream for demonstration purposes, not threadsafe
type Stream struct {
	pos    int
	tweets []Tweet
}

// ErrEOF returns on End of File error
var ErrEOF = errors.New("End of File")

// Next returns the next Tweet in the stream, returns EOF error if
// there are no more tweets
func (s *Stream) Next() (*Tweet, error) {

	// simulate delay
	time.Sleep(320 * time.Millisecond)
	if s.pos >= len(s.tweets) {
		return &Tweet{}, ErrEOF
	}

	tweet := s.tweets[s.pos]
	s.pos++

	return &tweet, nil
}

// Tweet defines the simlified representation of a tweet
type Tweet struct {
	Username string
	Text     string
}

// IsTalkingAboutGo is a mock process which pretend to be a sophisticated procedure to analyse whether tweet is talking about go or not
func (t *Tweet) IsTalkingAboutGo() bool {
	// simulate delay
	time.Sleep(330 * time.Millisecond)

	hasGolang := strings.Contains(strings.ToLower(t.Text), "golang")
	hasGopher := strings.Contains(strings.ToLower(t.Text), "gopher")

	return hasGolang || hasGopher
}

var mockdata = []Tweet{
	{
		"davecheney",
		"#golang top tip: if your unit tests import any other package you wrote, including themselves, they're not unit tests.",
	}, {
		"beertocode",
		"Backend developer, doing frontend featuring the eternal struggle of centering something. #coding",
	}, {
		"ironzeb",
		"Re: Popularity of Golang in China: My thinking nowadays is that it had a lot to do with this book and author https://github.com/astaxie/build-web-application-with-golang",
	}, {
		"beertocode",
		"Looking forward to the #gopher meetup in Hsinchu tonight with @ironzeb!",
	}, {
		"vampirewalk666",
		"I just wrote a golang slack bot! It reports the state of github repository. #Slack #golang",
	},
}

func producer(stream Stream, tweets chan *Tweet, quit chan bool) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			quit <- true
			close(tweets)
			return
		}
		tweets <- tweet
	}
}

func consumer(tweets chan *Tweet) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// use two channels so we can gracefully stop listening for the tweets channel
	tweets := make(chan *Tweet)
	quit := make(chan bool, 1)

	// Producer
	go producer(stream, tweets, quit)

	var done bool
	for !done {
		select {
		case <-quit:
			done = true
		default:
			// Consumer
			consumer(tweets)
		}
	}

	fmt.Printf("Process took %s\n", time.Since(start))
}
