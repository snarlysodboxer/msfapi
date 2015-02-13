package msfapi

import (
	"errors"
	"fmt"
)

type consoleInstanceResponse struct {
	ID     string
	Prompt string
	Busy   bool
}

func (api *API) ConsoleCreate() (consoleInstanceResponse, error) {
	var response map[string]interface{}
	var request = []string{"console.create", api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return consoleInstanceResponse{}, err
	}
	var consoleInstanceResponse = consoleInstanceResponse{
		ID:     response["id"].(string),
		Prompt: response["prompt"].(string),
		Busy:   response["busy"].(bool),
	}
	return consoleInstanceResponse, nil
}

func (api *API) ConsoleDestroy(consoleID string) error {
	var response map[string]string
	var request = []string{"console.destroy", api.Token, consoleID}
	err := api.request(&request, &response)
	if err != nil {
		return err
	}
	if response["result"] == "failure" {
		return errors.New(fmt.Sprintf("Invalid console ID %v", consoleID))
	}
	return nil
}

func (api *API) ConsoleList() ([]consoleInstanceResponse, error) {
	consoles := []consoleInstanceResponse{}
	var response map[string][]map[string]interface{}
	var request = []string{"console.list", api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return consoles, err
	}
	for _, console := range response["consoles"] {
		var console = consoleInstanceResponse{
			ID:     console["id"].(string),
			Prompt: console["prompt"].(string),
			Busy:   console["busy"].(bool),
		}
		consoles = append(consoles, console)
	}
	return consoles, nil
}

// NOTE: It's necessary to wait maybe 1 second after consoleWrite() before ConsoleRead()
func (api *API) ConsoleWrite(consoleID, command string) error {
	var response map[string]int
	var request = []string{"console.write", api.Token, consoleID, command}
	err := api.request(&request, &response)
	if err != nil {
		return err
	}
	if response["wrote"] != len(command) {
		return errors.New("Wrote length != command length")
	}
	return nil
}

type consoleReadResponse struct {
	Data   string
	Prompt string
	Busy   bool
}

func (api *API) ConsoleRead(consoleID string) (consoleReadResponse, error) {
	read := consoleReadResponse{}
	var response map[string]interface{}
	var request = []string{"console.read", api.Token, consoleID}
	err := api.request(&request, &response)
	if err != nil {
		return read, err
	}
	read.Data = response["data"].(string)
	read.Prompt = response["prompt"].(string)
	read.Busy = response["busy"].(bool)
	return read, nil
}
