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

type versionInfoResponse struct {
	Version string
	Ruby    string
	API     string
}

func New(url string) *API {
	// TODO ensure url responds before continuing
	api := new(API)
	api.URL = url
	return api
}

func (api *API) AuthLogin(user, pass string) error {
	var loginResponse map[string]interface{}
	var login = []string{"auth.login", user, pass}
	err := api.request(&login, &loginResponse)
	if err != nil {
		return err
	}
	if loginResponse["error"] != nil {
		if loginResponse["error"].(bool) {
			return errors.New(fmt.Sprintf("%#v, %#v",
				loginResponse["error_class"].(string),
				loginResponse["error_message"].(string),
			))
		}
	}
	api.Token = fmt.Sprintf("%v", loginResponse["token"])
	return nil
}

func (api *API) AuthLogout() error {
	var response map[string]interface{}
	var request = []string{"auth.logout", api.Token, api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return err
	}
	if response["error"] != nil {
		if response["error"].(bool) {
			return errors.New(fmt.Sprintf("%#v, %#v",
				response["error_class"].(string),
				response["error_message"].(string),
			))
		}
	}
	api.Token = fmt.Sprintf("%v", response["token"])
	return nil
}

func (api *API) CoreVersion() (versionInfoResponse, error) {
	if err := api.ensureToken(); err != nil {
		return versionInfoResponse{}, err
	}
	var response map[string]string
	var request = []string{"core.version", api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return versionInfoResponse{}, err
	}
	versionInfo := versionInfoResponse{
		Version: response["version"],
		Ruby:    response["ruby"],
		API:     response["api"],
	}
	return versionInfo, nil
}

func (api *API) ModuleExecute(mType, name string, mapp map[string]interface{}) (int64, error) {
	if err := api.ensureToken(); err != nil {
		return 0, err
	}
	var response map[string]interface{}
	var request = []interface{}{"module.execute", api.Token, mType, name, mapp}
	err := api.request(&request, &response)
	if err != nil {
		return 0, err
	}
	if response["error"] != nil {
		if response["error"].(bool) {
			return 0, errors.New(fmt.Sprintf("%#v, %#v",
				response["error_class"].(string),
				response["error_message"].(string),
			))
		}
	}
	jobID := int64ify(response["job_id"])
	return jobID, nil
}

func (api *API) JobList() (map[string]string, error) {
	if err := api.ensureToken(); err != nil {
		return map[string]string{}, err
	}
	var response map[string]string
	var request = []string{"job.list", api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return map[string]string{}, err
	}
	return response, nil
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
		// return errors.New(fmt.Sprintf("problem with marshaling:\n\t%v\n", err))
		return err
	}

	responseReader := bytes.NewReader(packedBytes)
	resp, err := http.Post(api.URL, "binary/message-pack", responseReader)
	if err != nil {
		// return errors.New(fmt.Sprintf("problem with posting:\n\t%v\n", err))
		return err
	}
	defer resp.Body.Close()

	receiveBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// return errors.New(fmt.Sprintf("problem with ReadAll:\n\t%v\n", err))
		return err
	}

	err = msgpack.Unmarshal(receiveBytes, response)
	if err != nil {
		// return errors.New(fmt.Sprintf("problem unmarshaling:\n\t%v", err))
		return err
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
