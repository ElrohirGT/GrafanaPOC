package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
)

func PostObj(ctx context.Context) error {
	url := "https://api.restful-api.dev/objects"

	log.Println("Adding obj to service...")
	reqBody := `
{
   "name": "Apple MacBook Pro 16",
   "data": {
      "year": 2019,
      "price": 1849.99,
      "CPU model": "Intel Core i9",
      "Hard disk size": "1 TB"
   }
}`
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	log.Println("Service responded!")
	if err != nil {
		return err
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
