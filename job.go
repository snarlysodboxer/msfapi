package msfapi

import ()

func (api *API) JobList() (map[string]string, error) {
	var response map[string]string
	var request = []string{"job.list", api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return map[string]string{}, err
	}
	return response, nil
}
