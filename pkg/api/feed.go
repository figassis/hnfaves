package api

import (
	"fmt"
	"sort"
	"time"

	"github.com/figassis/hnfaves/pkg/utl/util"
	"github.com/figassis/hnfaves/pkg/utl/zaplog"
	"github.com/gorilla/feeds"
	echo "github.com/labstack/echo/v4"
)

func apiFeed(c echo.Context, user string, page int64) (*feeds.Feed, error) {
	return feed(user, page)
}

func feed(user string, page int64) (feed *feeds.Feed, err error) {
	url := fmt.Sprintf("https://hnfaves.com/%s", user)
	if page > 0 {
		url = fmt.Sprintf("https://hnfaves.com/%s&p=%d", user, page)
	}

	items, err := util.Crawl(user, page)
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Index < items[j].Index })

	now := time.Now()
	feed = &feeds.Feed{
		Title:       fmt.Sprintf("%s's Hacker News favorites", user),
		Link:        &feeds.Link{Href: url},
		Description: fmt.Sprintf("%s's Hacker News favorites", user),
		Author:      &feeds.Author{Name: user, Email: fmt.Sprintf("%s@hnfaves.com", user)},
		Created:     now,
		Updated:     now,
	}

	if len(items) > 0 {
		feed.Created = items[len(items)-1].Created
	}

	for _, item := range items {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: item.Link},
			Description: item.Title,
			Author:      feed.Author,
			Created:     item.Created,
			Source:      &feeds.Link{Href: "https://news.ycombinator.com"},
			Id:          fmt.Sprintf("https://news.ycombinator.com/%s", item.ID),
		})
	}

	return
}
