package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/figassis/hnfaves/pkg/utl/zaplog"

	echo "github.com/labstack/echo/v4"
)

func getFeed(c echo.Context) error {
	user := c.Param("user")
	if user == "" {
		if err := zaplog.ZLog(errors.New("Invalid user")); err != nil {
			return err
		}
	}

	var page int64
	if tempPage, err := strconv.Atoi(c.QueryParam("p")); err == nil && tempPage > 0 {
		page = int64(tempPage)
	}

	rss, err := apiFeed(c, user, page)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	header := "application/xml; charset=UTF-8"
	var r string

	switch c.QueryParam("type") {
	case "atom":
		r, err = rss.ToAtom()
	case "json":
		header = "application/json; charset=UTF-8"
		r, err = rss.ToJSON()
	default:
		r, err = rss.ToRss()
	}

	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", header)
	return c.String(http.StatusOK, r)
}
