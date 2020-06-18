package ratelimiter

import (
	tollbooth "github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/errors"
	limiter "github.com/didip/tollbooth/v6/limiter"
	"github.com/labstack/echo/v4"
)

func LimitMiddleware(lmt *limiter.Limiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			httpError := tollbooth.LimitByRequest(lmt, c.Response(), c.Request())
			if httpError != nil {
				return c.String(httpError.StatusCode, httpError.Message)
			}
			return next(c)
		})
	}
}

func LimitMiddlewareByTier(free, paid *limiter.Limiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			headers := paid.GetHeaders()
			var isPaid bool
			for header := range headers {
				if c.Request().Header.Get(header) != "" {
					isPaid = true
					break
				}
			}

			var httpError *errors.HTTPError
			if isPaid {
				httpError = tollbooth.LimitByRequest(paid, c.Response(), c.Request())
			} else {
				httpError = tollbooth.LimitByRequest(free, c.Response(), c.Request())
			}

			if httpError != nil {
				return c.String(httpError.StatusCode, httpError.Message)
			}
			return next(c)
		})
	}
}

func LimitHandler(lmt *limiter.Limiter) echo.MiddlewareFunc {
	return LimitMiddleware(lmt)
}

func LimitHandlerByTier(free, paid *limiter.Limiter) echo.MiddlewareFunc {
	return LimitMiddlewareByTier(free, paid)
}
