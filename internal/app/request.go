package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var errNoConnectionToPhonesServer = errors.New("haven't connection to phones server address")

type client struct {
	phonesServerAddress string
	httpC               *http.Client
}

func newClient() *client {

	newClient := &client{}

	newClient.phonesServerAddress = "http://localhost:3333"

	newClient.httpC = &http.Client{}

	return newClient

}

func (c *client) getPhone(id string) (phone []byte, statusCode int, err error) {

	resp, err := c.httpC.Get(c.phonesServerAddress + "/?id=" + id)
	if err != nil {
		return nil, 0, fmt.Errorf("trying to get phone %s: %w", err, errNoConnectionToPhonesServer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}

	phone, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("reading resp body: %w", err)
	}

	return phone, resp.StatusCode, nil

}
