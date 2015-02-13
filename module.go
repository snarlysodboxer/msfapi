package msfapi

import (
	"errors"
	"fmt"
)

func (api *API) ModuleExecute(mType, name string, mapp map[string]interface{}) (int64, error) {
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
