package util

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/figassis/hnfaves/pkg/utl/zaplog"
	colly "github.com/gocolly/colly/v2"
)

const (
	maxThreads  = 4
	randomDelay = 2 // seconds
)

type (
	NewsItem struct {
		ID    string
		Title string
		Link  string
		Index int64
	}
)

func Crawl(user string, page int64) (result []NewsItem, err error) {
	url := fmt.Sprintf("https://news.ycombinator.com/favorites?id=%s", user)
	if page > 0 {
		url = fmt.Sprintf("https://news.ycombinator.com/favorites?id=%s&p=%d", user, page)
	}

	if err2 := zaplog.ZLog(GetCache(url, &result)); err2 == nil && len(result) > 0 {
		return
	}

	// Instantiate default collector and visit only approved domains
	c := colly.NewCollector(
		colly.AllowedDomains("news.ycombinator.com"),
		colly.Async(false),
		colly.URLFilters(
			//Only visit links mathcing filter
			regexp.MustCompile(`https://news\.ycombinator\.com/favorites\?id=.*`),
		),
	)

	if os.Getenv("CACHE") == "true" {
		c.CacheDir = "./cache"
	}

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	if err = zaplog.ZLog(c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: maxThreads, RandomDelay: randomDelay * time.Second})); err != nil {
		return
	}

	// Get data-module elements with country advisory content
	c.OnHTML(".athing", func(e *colly.HTMLElement) {

		item := NewsItem{
			ID: e.Attr("id"),
		}

		e.ForEach(".storylink", func(_ int, el *colly.HTMLElement) {
			item.Title = el.Text
			item.Link = el.Attr("href")
			if strings.HasPrefix(item.Link, "item?id") {
				item.Link = fmt.Sprintf("https://news.ycombinator.com/%s", item.Link)
			}
		})

		e.ForEach(".rank", func(_ int, el *colly.HTMLElement) {
			index, err := strconv.Atoi(strings.TrimSuffix(el.Text, "."))
			if err == nil {
				item.Index = int64(index)
			}
		})

		if item.ID != "" && item.Link != "" && item.Title != "" {
			result = append(result, item)
		}
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		zaplog.ZLog(fmt.Sprintf("Request URL: %s failed with response: %v", r.Request.URL, err))
	})

	err = zaplog.ZLog(c.Visit(url))
	go CacheTTL(url, result, time.Minute*10)
	return
}
