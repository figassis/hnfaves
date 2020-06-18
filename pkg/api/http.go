package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/figassis/hnfaves/pkg/utl/zaplog"

	echo "github.com/labstack/echo/v4"
)

func getFeed(c echo.Context) error {
	user := c.Param("token")
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

	switch c.Param("type") {
	case "atom":
		r, err := rss.ToAtom()
		if err = zaplog.ZLog(err); err != nil {
			return err
		}
		return c.String(http.StatusOK, r)
	}

	r, err := rss.ToRss()
	if err = zaplog.ZLog(err); err != nil {
		return err
	}
	return c.String(http.StatusOK, r)
}
