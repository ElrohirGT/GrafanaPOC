package main

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
)

var ErrSimulatedError = errors.New("emulating: something went wrong during parsing of body")

func GetUser(ctx context.Context) error {
	url := "https://randomuser.me/api/"

	log.Println("Getting user from api...")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	log.Println(ctx, "Service responded!")
	if err != nil {
		return err
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if rand.Float32() < 0.1 {
		return ErrSimulatedError
	}

	return nil
}
