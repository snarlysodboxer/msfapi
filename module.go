package msfapi

import ()

func (api *API) ModuleExecute(mType, name string, mapp map[string]interface{}) (int64, error) {
	var response map[string]interface{}
	var request = []interface{}{"module.execute", api.Token, mType, name, mapp}
	err := api.request(&request, &response)
	if err != nil {
		return 0, err
	}
	jobID := int64ify(response["job_id"])
	return jobID, nil
}
