package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (app *Config) handler(urlString string) gin.HandlerFunc {
	// define URLs for microservices
	parsedURL, _ := url.Parse(urlString)

	// create reverse proxies for each microservice
	proxy := app.reverseProxy(parsedURL)

	return func(c *gin.Context) {
		c.Request.Host = parsedURL.Host

		// request forwarding (rewriting)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (app *Config) reverseProxy(address *url.URL) *httputil.ReverseProxy {
	// Create a new reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(address)

	proxy.Director = func(request *http.Request) {
		request.Host = address.Host
		request.URL.Scheme = address.Scheme
		request.URL.Host = address.Host
		request.URL.Path = address.Path
	}

	proxy.ModifyResponse = modifyResponse()

	return proxy
}

// modifies the response and logs errors
func modifyResponse() func(response *http.Response) error {
	return func(response *http.Response) error {
		if response.StatusCode == http.StatusInternalServerError {
			log.Println("Something went wrong!! Internal Server Error!!")
		} else if response.StatusCode == 200 {
			log.Println("Woohoo it worked!")
		}
		return nil
	}
}
