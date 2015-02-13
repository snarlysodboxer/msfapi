package msfapi

import ()

type CoreVersionResponse struct {
	Version string
	Ruby    string
	API     string
}

func (api *API) CoreVersion() (CoreVersionResponse, error) {
	var response map[string]string
	var request = []string{"core.version", api.Token}
	err := api.request(&request, &response)
	if err != nil {
		return CoreVersionResponse{}, err
	}
	versionInfo := CoreVersionResponse{
		Version: response["version"],
		Ruby:    response["ruby"],
		API:     response["api"],
	}
	return versionInfo, nil
}
