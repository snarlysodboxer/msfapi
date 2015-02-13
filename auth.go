package msfapi

import (
	"errors"
	"fmt"
)

func (api *API) AuthLogin(user, pass string) (string, error) {
	var loginResponse map[string]interface{}
	var login = []string{"auth.login", user, pass}
	err := api.request(&login, &loginResponse)
	if err != nil {
		return "", err
	}
	if loginResponse["error"] != nil {
		if loginResponse["error"].(bool) {
			return "", errors.New(fmt.Sprintf("%#v, %#v",
				loginResponse["error_class"].(string),
				loginResponse["error_message"].(string),
			))
		}
	}
	token := fmt.Sprintf("%v", loginResponse["token"])
	return token, nil
}

func (api *API) AuthTokenAdd(token string) error {
	var response map[string]string
	var request = []string{"auth.token_add", api.Token, token}
	err := api.request(&request, &response)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) AuthTokenGenerate() (string, error) {
	var response map[string]string
	var request = []string{"auth.token_generate", api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return "", err
	}
	token := response["token"]
	return token, nil
}

func (api *API) AuthTokenList() ([]string, error) {
	var response map[string][]string
	var request = []string{"auth.token_list", api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return []string{}, err
	}
	tokens := response["tokens"]
	return tokens, nil
}

func (api *API) AuthTokenRemove(token string) error {
	var response map[string]string
	var request = []string{"auth.token_remove", api.Token, token}
	err := api.request(&request, &response)
	if err != nil {
		return err
	}
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
