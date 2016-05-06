package msfapi

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io/ioutil"
	"net/http"
)

type API struct {
	Token string
	URL   string
}

func New(url string) *API {
	// TODO ensure url responds before continuing
	api := new(API)
	api.URL = url
	return api
}

func (api *API) ensureToken() error {
	if api.Token == "" {
		return errors.New("Token is empty for some reason")
	}
	return nil
}

func (api *API) request(request, response interface{}) error {
	packedBytes, err := msgpack.Marshal(request)
	if err != nil {
		return errors.New(fmt.Sprintf("problem with marshaling:\n\t%v\n", err))
	}

	responseReader := bytes.NewReader(packedBytes)
	resp, err := http.Post(api.URL, "binary/message-pack", responseReader)
	if err != nil {
		return errors.New(fmt.Sprintf("problem with posting:\n\t%v\n", err))
	}
	defer resp.Body.Close()

	receiveBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("problem with ReadAll:\n\t%v\n", err))
	}

	switch resp.StatusCode {
	case 200:
		// "The request was successfully processed"
		var stringInterface map[string]interface{}
		err = msgpack.Unmarshal(receiveBytes, &stringInterface)
		if err == nil {
			if stringInterface["error"] != nil {
				if stringInterface["error"].(bool) {
					return errors.New(fmt.Sprintf("%s, %s", stringInterface["error_class"].(string), string(stringInterface["error_message"].([]uint8))))
				}
			}
		}
		err = msgpack.Unmarshal(receiveBytes, response)
		if err != nil {
			return errors.New(fmt.Sprintf("Error unmarshaling response: %s\n\t%s", err, string(receiveBytes)))
		}
		return nil
	case 500:
		return errors.New(fmt.Sprintf("The request resulted in an error: \n%s", string(receiveBytes)))
	case 401:
		return errors.New(fmt.Sprintf("The authentication credentials supplied were not valid"))
	case 403:
		return errors.New(fmt.Sprintf("The authentication credentials supplied were not granted access to the resource"))
	case 404:
		return errors.New(fmt.Sprintf("The request was sent to an invalid URI"))
	default:
		return errors.New(resp.Status)
	}
	return nil
}

func int64ify(n interface{}) int64 {
	switch n := n.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case uint:
		return int64(n)
	case uintptr:
		return int64(n)
	case uint8:
		return int64(n)
	case uint16:
		return int64(n)
	case uint32:
		return int64(n)
	case uint64:
		return int64(n)
	}
	return int64(0)
}
