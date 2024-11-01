/*
* Copyright © 2024 Minand Manomohanan <minand.nell.mohan@gmail.com>
 */
package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandlerGet(c *fiber.Ctx, origin *string) error {
	urlPath := c.Path()

	// Strip stray / if it is there in origin
	trimmedOrigin := strings.TrimSuffix(*origin, "/")

	url := trimmedOrigin + urlPath

	// Return cached response if available
	cachedResponse, ok := urlCache.CheckAndReturnSavedResponse(url)
	if ok {
		fmt.Println("Returning value from cache...")
		c.Status(cachedResponse.Code)
		return c.JSON(cachedResponse.Body)
	}

	// Initialize request
	fmt.Println("redirecting to " + url)
	headers := c.GetReqHeaders()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Handle headers with multiple values
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make GET call
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to get Response! due to " + err.Error())
		return err
	}

	// Set status code
	c.Status(resp.StatusCode)

	// Handle if response is encoded
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("Failed to create gzip reader! due to " + err.Error())
			return err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	// Read body
	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Failed to read response body! due to " + err.Error())
		return err
	}
	defer resp.Body.Close()

	// Convert the body bytes to a JSON object
	var jsonBody interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		fmt.Println("Failed to unmarshal response body")
		return c.SendString(string(body))
	}

	// Put new entry into cache
	go urlCache.PutNewEntryInCache(url, resp.StatusCode, jsonBody)

	fmt.Println("Returning value from server..")
	return c.JSON(jsonBody)

}
