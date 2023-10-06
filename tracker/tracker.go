package tracker

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Click struct {
	id          string
	userID      string
	code        string
	destination string
	host        string
	requestURI  string
	remoteIP    string
	referer     string
	userAgent   string
	clickedAt   time.Time
}

type ClickLogger struct {
}

func ShortCodeHandler(c *gin.Context) {
	code := c.Param("code")
	//userID := 123
	path := "/go/123?code=code&url=https%3A%2F%2Frefactoring.guru%2Fdesign-patterns%2Fsingleton%2Fgo%2Fexample"
	c.Request.URL.Path = path

	c.String(http.StatusOK, code)
}

func GoHandler(c *gin.Context) {
	userId := c.Param("id")
	code := c.Query("code")
	link := c.Query("url")

	if userId == "" || link == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found.",
		})
		return
	}

	_, errParseURI := url.ParseRequestURI(link)
	// Cannot parse URL
	if errParseURI != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    "URL_IS_INVALID",
			"message": "URL is invalid.",
			"url":     link,
		})
		return
	}

	click := Click{
		userID:      userId,
		code:        code,
		destination: link,
		host:        c.Request.Host,
		requestURI:  c.Request.RequestURI,
		remoteIP:    c.ClientIP(),
		referer:     c.Request.Referer(),
		userAgent:   c.Request.UserAgent(),
		clickedAt:   time.Now(),
	}

	logClick(&click)
	// OK, redirect to destination URL.
	c.Redirect(http.StatusTemporaryRedirect, link)
}

func logClick(click *Click) {
	f, _ := os.OpenFile(
		"./logs/tracking.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	logger := zerolog.New(f).With().Logger()
	logger.Info().
		Str("user_id", click.userID).
		Str("code", click.code).
		Str("destination", click.destination).
		Str("host", click.host).
		Str("request_uri", click.requestURI).
		Str("remote_ip", click.remoteIP).
		Str("referer", click.referer).
		Str("user_agent", click.userAgent).
		Str("clicked_at", click.clickedAt.String()).
		Msg("")
}
