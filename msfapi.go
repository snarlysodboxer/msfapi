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

	err = msgpack.Unmarshal(receiveBytes, response)
	if err != nil {
		return errors.New(fmt.Sprintf("problem unmarshaling %v:\n\t%v", receiveBytes, err))
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
