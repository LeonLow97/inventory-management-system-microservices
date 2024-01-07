package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func handler(urlString string) gin.HandlerFunc {
	// define URLs for microservices
	parsedURL, _ := url.Parse(urlString)

	// create reverse proxies for each microservice
	proxy := reverseProxy(parsedURL)

	return func(c *gin.Context) {
		c.Request.Host = parsedURL.Host

		// request forwarding (rewriting)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func reverseProxy(address *url.URL) *httputil.ReverseProxy {
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
		} else if response.StatusCode == http.StatusOK {
			log.Println("Woohoo it worked!")
		} else {
			// log via gRPC to logger service
			// Read the response body
			body, err := io.ReadAll(response.Body)
			if err != nil {
				// Handle error reading the response body
				return err
			}

			// Log the received JSON response
			log.Println("Received JSON response:", string(body))

			// Re-assign the response body with the original content
			response.Body = io.NopCloser(bytes.NewReader(body))
		}
		return nil
	}
}
