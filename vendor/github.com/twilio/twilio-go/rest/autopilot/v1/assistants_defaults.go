/*
 * Twilio - Autopilot
 *
 * This is the public Twilio REST API.
 *
 * API version: 1.22.0
 * Contact: support@twilio.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"net/url"

	"strings"
)

func (c *ApiService) FetchDefaults(AssistantSid string) (*AutopilotV1Defaults, error) {
	path := "/v1/Assistants/{AssistantSid}/Defaults"
	path = strings.Replace(path, "{"+"AssistantSid"+"}", AssistantSid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &AutopilotV1Defaults{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Optional parameters for the method 'UpdateDefaults'
type UpdateDefaultsParams struct {
	// A JSON string that describes the default task links for the `assistant_initiation`, `collect`, and `fallback` situations.
	Defaults *map[string]interface{} `json:"Defaults,omitempty"`
}

func (params *UpdateDefaultsParams) SetDefaults(Defaults map[string]interface{}) *UpdateDefaultsParams {
	params.Defaults = &Defaults
	return params
}

func (c *ApiService) UpdateDefaults(AssistantSid string, params *UpdateDefaultsParams) (*AutopilotV1Defaults, error) {
	path := "/v1/Assistants/{AssistantSid}/Defaults"
	path = strings.Replace(path, "{"+"AssistantSid"+"}", AssistantSid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.Defaults != nil {
		v, err := json.Marshal(params.Defaults)

		if err != nil {
			return nil, err
		}

		data.Set("Defaults", string(v))
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &AutopilotV1Defaults{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}
