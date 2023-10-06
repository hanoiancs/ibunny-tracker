package tracer

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

var tr = &http.Transport{
	MaxIdleConns:    10,
	IdleConnTimeout: 5 * time.Second,
}

var client = &http.Client{
	Transport: tr,
	Timeout:   5 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type URL struct {
	Location         string
	Status           int
	RedirectLocation string
}

func trace(link string) (*URL, error) {
	_, errParseURI := url.ParseRequestURI(link)

	if errParseURI != nil {
		return nil, errParseURI
	}

	u := URL{
		Location:         link,
		RedirectLocation: "",
		Status:           http.StatusOK,
	}
	req, _ := http.NewRequest("GET", link, nil)
	req.Header.Add(
		"User-Agent",
		`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`,
	)
	resp, errRequest := client.Do(req)

	if errRequest != nil {
		return nil, errRequest
	} else {
		location, errLocation := resp.Location()
		u.Status = resp.StatusCode
		if errLocation == nil {
			u.RedirectLocation = location.String()
		}
	}

	//if resp.StatusCode == http.StatusOK {
	//	defer func(Body io.ReadCloser) {
	//		errBodyClose := Body.Close()
	//		if errBodyClose != nil {
	//
	//		}
	//	}(resp.Body)
	//	body, errReadBody := io.ReadAll(resp.Body)
	//	if errReadBody == nil {
	//		fmt.Println(string(body))
	//	}
	//}

	return &u, nil
}

func Trace(link string) (*[]URL, error) {
	var urls []URL

	n := 1
	for link != "" && n <= 5 {
		tracer, err := trace(link)

		if tracer != nil {
			urls = append(urls, *tracer)
		}

		if err != nil {
			return &urls, err
		}

		if tracer.RedirectLocation == "" {
			break
		}
		link = tracer.RedirectLocation
		n += 1
	}

	return &urls, nil
}

func TraceHandler(c *gin.Context) {
	link := c.PostForm("url")
	urls, err := Trace(link)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  "ERROR",
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "OK",
		"trace":   urls,
	})
}
