package main

import "fmt"
import "net/http"
import "io/ioutil"
import "strings"
import "sync"
import "time"

const worker_count = 10
const task_queue_count = 128

var cache *Cache

func LogActivity(activity func() error, errPrefix string, success string) {
	if err := activity(); err != nil {
		fmt.Println(errPrefix, err)
		return
	}
	fmt.Println(success)
}

var getLock sync.Mutex

func HttpGet(url string) string {

	if val, ok := cache.Get(url); ok {
		return val
	}

	// only allow 1 request at a time, and they each block for
	// at least 1 second to keep us at <= 1 request/second
	getLock.Lock()
	defer getLock.Unlock()
	time.Sleep(10 * time.Second)
	fmt.Printf("GET %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}

	bin, err := ioutil.ReadAll(resp.Body)
	data := string(bin)
	if err != nil {
		return ""
	}

	cache.Put(url, data)
	cache.Save() // This is a bad idea...
	return data
}

func main() {
	cache = NewCache("cache.json")
	LogActivity(cache.Load, "Error loading cache: ", "Cache loaded.")
	defer LogActivity(cache.Save, "Error writing cache: ", "Cache saved.")

	runner := NewRunner()
	runner.tasks <- scrapeMain
	runner.Run()
}

func scrapeMain(tasks chan<- Task) {
	// Get main page
	data := HttpGet("http://alibaba.com/vn")

	// Break it into [bad, goodprefix, goodprefix, ... ]
	parts := strings.Split(data, "\n                            <a href=\"")
	for _, line := range parts[1:] {
		url := strings.Split(line, "\"")[0]
		tasks <- scrapeCategory(url)
	}
}
func scrapeCategory(url string) Task {
	return func(tasks chan<- Task) {
		// Get the html
		data := HttpGet(url)

		// There's possibly a next page, enqueue that
		next := strings.Index(data, "class=\"next\"")
		if next != -1 {
			next += 21 // skip to url part
			fmt.Println("Next: ", strings.Split(data[next:], "\"")[0])
		}

		// split all the company_profile links into
		// trash + last-quote + url
		parts := strings.Split(data, "company_profile.html")

		parts = parts[:len(parts)-1] //all but last
		for _, line := range parts {
			idx := strings.LastIndex(line, "\"") + 1
			//fmt.Println(line[idx:])
			_ = line[idx:]
		}
	}
}
