// Package main compares performance among three identical apis
// written in Go, Sinatra, and Rails
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	url2 "net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	contentType       = "Content-Type"
	applicationJSON   = "application/json"
	xUserID           = "X-User-Id"
	httpGET           = "GET"
	httpStr           = "http"
	localhost         = "localhost"
	golang            = "go"
	postgrest         = "postgrest"
	goHTTPPort        = "3001"
	railsHTTPPort     = "3000"
	sinatraHTTPPort   = "4567"
	postgrestHTTPPort = "3002"
	rails             = "rails"
	sinatra           = "sinatra"
	itemsPath         = "/api/v1/items"
	userIDsPath       = "/api/v1/user_ids"
)

var (
	client    *http.Client
	oneClient sync.Once
	workingOn = "working on %v\n"
)

// work is called workerCt times, each in its own goroutine
// It ranges over requestCh, grabs a request, performs that request, then
// sends each response to responseCh
func work(requestCh <-chan *http.Request, responseCh chan<- *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range requestCh {
		result, err := testClient().Do(req)
		if err != nil {
			panic(err)
		}
		responseCh <- result
	}
}

type app struct {
	name string
	port string
	path string
}

type result struct {
	name        string
	elapsedTime float64
}

type report []result

// String() implements the Stringer interface, so we get a nice
// clean reporting output
func (r report) String() string {
	output := []string{}
	sort.Slice(
		r, func(i, j int) bool {
			return r[i].elapsedTime < r[j].elapsedTime
		},
	)
	for _, res := range r {
		formatted := fmt.Sprintf("%10v: %v", res.name, res.elapsedTime)
		output = append(output, formatted)
	}
	return strings.Join(output, "\n")
}

// main iterates over three apps (Go, Rails, Sinatra) and performs
// 1000 requests against each app's identical api/v1/items JSON api,
// using 10 concurrent workers, then prints out the results comparing
// the performance of all three.
func main() {
	apps := []app{
		{
			name: golang,
			port: goHTTPPort,
			path: "/api/v1/items",
		},
		{
			name: rails,
			port: railsHTTPPort,
			path: "/api/v1/items",
		},
		{
			name: sinatra,
			port: sinatraHTTPPort,
			path: "/api/v1/items",
		},
		{
			name: postgrest,
			port: postgrestHTTPPort,
			path: "/get_user_items",
		},
	}

	userIds := userIDs()

	report := report{}

	for _, a := range apps {
		fmt.Printf(workingOn, a.name)
		startTime := time.Now()
		workerCt := 10
		reqCt := 1000

		// Create channels for requestCh (URLs) and responseCh
		requestCh := make(chan *http.Request, reqCt)
		responseCh := make(chan *http.Response)

		// Create a WaitGroup to wait for all workers to finish
		var wg sync.WaitGroup

		// Create worker pool
		for i := 1; i <= workerCt; i++ {
			wg.Add(1)
			go work(requestCh, responseCh, &wg)
		}

		requests := make([]*http.Request, reqCt)
		for i := 0; i < reqCt; i++ {
			randOffset := rand.Intn(10) * 10
			rawQuery := map[string]string{
				"limit":  "10",
				"offset": fmt.Sprintf("%d", randOffset),
			}
			requests[i] = createRequest(localhost, a.port, a.path, rawQuery, randomUserID(userIds))
		}

		go func() {
			for _, req := range requests {
				requestCh <- req
			}
			close(requestCh)
		}()

		go func() {
			wg.Wait()
			close(responseCh)
		}()

		for res := range responseCh {
			bytes, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(bytes))
		}
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		report = append(report, result{name: a.name, elapsedTime: elapsedTime.Seconds()})
	}
	fmt.Println(report)
}

// testClient uses sync.Once to ensure only one http.Client is ever made
func testClient() *http.Client {
	oneClient.Do(
		func() {
			client = http.DefaultClient
		},
	)
	return client
}

// createRequest creates a request using a random userID
func createRequest(host, port, path string, rawQuery map[string]string, userID int) *http.Request {
	rawQueryString, rawQuerySl := "", []string{}
	if rawQuery != nil {
		for k, v := range rawQuery {
			rawQuerySl = append(rawQuerySl, fmt.Sprintf("%v=%v", k, v))
		}
		rawQueryString = strings.Join(rawQuerySl, "&")
	}

	url := url2.URL{
		Scheme:   httpStr,
		Host:     fmt.Sprintf("%v:%v", host, port),
		Path:     path,
		RawQuery: rawQueryString,
	}

	header := map[string][]string{
		contentType: {applicationJSON},
		xUserID:     {fmt.Sprintf("%v", userID)},
	}

	r := http.Request{
		Method: httpGET,
		URL:    &url,
		Header: http.Header(header),
	}
	return &r
}

type user struct {
	Id int `json:"id"`
}

// randomUserID grabs a random user_id for each request
func randomUserID(ids []int) int {
	randInt := rand.Intn(len(ids))
	return ids[randInt]
}

// userIDs grabs all userIDs from sinatra-items-api, so we can grab a random
// user_id for each request
func userIDs() []int {
	users := []user{}
	// This request does not require auth, userID == 0
	req := createRequest(localhost, sinatraHTTPPort, userIDsPath, nil, 0)
	res, err := testClient().Do(req)
	if err != nil {
		panic(err)
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		panic(err)
	}
	ids := make([]int, len(users))
	for i := range users {
		ids[i] = users[i].Id
	}
	return ids
}
