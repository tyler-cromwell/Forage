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
	"fmt"
	"net/url"

	"strings"

	"github.com/twilio/twilio-go/client"
)

// Optional parameters for the method 'CreateQuery'
type CreateQueryParams struct {
	// The [ISO language-country](https://docs.oracle.com/cd/E13214_01/wli/docs92/xref/xqisocodes.html) string that specifies the language used for the new query. For example: `en-US`.
	Language *string `json:"Language,omitempty"`
	// The SID or unique name of the [Model Build](https://www.twilio.com/docs/autopilot/api/model-build) to be queried.
	ModelBuild *string `json:"ModelBuild,omitempty"`
	// The end-user's natural language input. It can be up to 2048 characters long.
	Query *string `json:"Query,omitempty"`
	// The list of tasks to limit the new query to. Tasks are expressed as a comma-separated list of task `unique_name` values. For example, `task-unique_name-1, task-unique_name-2`. Listing specific tasks is useful to constrain the paths that a user can take.
	Tasks *string `json:"Tasks,omitempty"`
}

func (params *CreateQueryParams) SetLanguage(Language string) *CreateQueryParams {
	params.Language = &Language
	return params
}
func (params *CreateQueryParams) SetModelBuild(ModelBuild string) *CreateQueryParams {
	params.ModelBuild = &ModelBuild
	return params
}
func (params *CreateQueryParams) SetQuery(Query string) *CreateQueryParams {
	params.Query = &Query
	return params
}
func (params *CreateQueryParams) SetTasks(Tasks string) *CreateQueryParams {
	params.Tasks = &Tasks
	return params
}

func (c *ApiService) CreateQuery(AssistantSid string, params *CreateQueryParams) (*AutopilotV1Query, error) {
	path := "/v1/Assistants/{AssistantSid}/Queries"
	path = strings.Replace(path, "{"+"AssistantSid"+"}", AssistantSid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.Language != nil {
		data.Set("Language", *params.Language)
	}
	if params != nil && params.ModelBuild != nil {
		data.Set("ModelBuild", *params.ModelBuild)
	}
	if params != nil && params.Query != nil {
		data.Set("Query", *params.Query)
	}
	if params != nil && params.Tasks != nil {
		data.Set("Tasks", *params.Tasks)
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &AutopilotV1Query{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

func (c *ApiService) DeleteQuery(AssistantSid string, Sid string) error {
	path := "/v1/Assistants/{AssistantSid}/Queries/{Sid}"
	path = strings.Replace(path, "{"+"AssistantSid"+"}", AssistantSid, -1)
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	resp, err := c.requestHandler.Delete(c.baseURL+path, data, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (c *ApiService) FetchQuery(AssistantSid string, Sid string) (*AutopilotV1Query, error) {
	path := "/v1/Assistants/{AssistantSid}/Queries/{Sid}"
	path = strings.Replace(path, "{"+"AssistantSid"+"}", AssistantSid, -1)
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &AutopilotV1Query{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Optional parameters for the method 'ListQuery'
type ListQueryParams struct {
	// The [ISO language-country](https://docs.oracle.com/cd/E13214_01/wli/docs92/xref/xqisocodes.html) string that specifies the language used by the Query resources to read. For example: `en-US`.
	Language *string `json:"Language,omitempty"`
	// The SID or unique name of the [Model Build](https://www.twilio.com/docs/autopilot/api/model-build) to be queried.
	ModelBuild *string `json:"ModelBuild,omitempty"`
	// The status of the resources to read. Can be: `pending-review`, `reviewed`, or `discarded`
	Status *string `json:"Status,omitempty"`
	// The SID of the [Dialogue](https://www.twilio.com/docs/autopilot/api/dialogue).
	DialogueSid *string `json:"DialogueSid,omitempty"`
	// How many resources to return in each list page. The default is 50, and the maximum is 1000.
	PageSize *int `json:"PageSize,omitempty"`
	// Max number of records to return.
	Limit *int `json:"limit,omitempty"`
}

func (params *ListQueryParams) SetLanguage(Language string) *ListQueryParams {
	params.Language = &Language
	return params
}
func (params *ListQueryParams) SetModelBuild(ModelBuild string) *ListQueryParams {
	params.ModelBuild = &ModelBuild
	return params
}
func (params *ListQueryParams) SetStatus(Status string) *ListQueryParams {
	params.Status = &Status
	return params
}
func (params *ListQueryParams) SetDialogueSid(DialogueSid string) *ListQueryParams {
	params.DialogueSid = &DialogueSid
	return params
}
func (params *ListQueryParams) SetPageSize(PageSize int) *ListQueryParams {
	params.PageSize = &PageSize
	return params
}
func (params *ListQueryParams) SetLimit(Limit int) *ListQueryParams {
	params.Limit = &Limit
	return params
}

// Retrieve a single page of Query records from the API. Request is executed immediately.
func (c *ApiService) PageQuery(AssistantSid string, params *ListQueryParams, pageToken, pageNumber string) (*ListQueryResponse, error) {
	path := "/v1/Assistants/{AssistantSid}/Queries"

	path = strings.Replace(path, "{"+"AssistantSid"+"}", AssistantSid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.Language != nil {
		data.Set("Language", *params.Language)
	}
	if params != nil && params.ModelBuild != nil {
		data.Set("ModelBuild", *params.ModelBuild)
	}
	if params != nil && params.Status != nil {
		data.Set("Status", *params.Status)
	}
	if params != nil && params.DialogueSid != nil {
		data.Set("DialogueSid", *params.DialogueSid)
	}
	if params != nil && params.PageSize != nil {
		data.Set("PageSize", fmt.Sprint(*params.PageSize))
	}

	if pageToken != "" {
		data.Set("PageToken", pageToken)
	}
	if pageNumber != "" {
		data.Set("Page", pageNumber)
	}

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ListQueryResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Lists Query records from the API as a list. Unlike stream, this operation is eager and loads 'limit' records into memory before returning.
func (c *ApiService) ListQuery(AssistantSid string, params *ListQueryParams) ([]AutopilotV1Query, error) {
	if params == nil {
		params = &ListQueryParams{}
	}
	params.SetPageSize(client.ReadLimits(params.PageSize, params.Limit))

	response, err := c.PageQuery(AssistantSid, params, "", "")
	if err != nil {
		return nil, err
	}

	curRecord := 0
	var records []AutopilotV1Query

	for response != nil {
		records = append(records, response.Queries...)

		var record interface{}
		if record, err = client.GetNext(c.baseURL, response, &curRecord, params.Limit, c.getNextListQueryResponse); record == nil || err != nil {
			return records, err
		}

		response = record.(*ListQueryResponse)
	}

	return records, err
}

// Streams Query records from the API as a channel stream. This operation lazily loads records as efficiently as possible until the limit is reached.
func (c *ApiService) StreamQuery(AssistantSid string, params *ListQueryParams) (chan AutopilotV1Query, error) {
	if params == nil {
		params = &ListQueryParams{}
	}
	params.SetPageSize(client.ReadLimits(params.PageSize, params.Limit))

	response, err := c.PageQuery(AssistantSid, params, "", "")
	if err != nil {
		return nil, err
	}

	curRecord := 0
	//set buffer size of the channel to 1
	channel := make(chan AutopilotV1Query, 1)

	go func() {
		for response != nil {
			for item := range response.Queries {
				channel <- response.Queries[item]
			}

			var record interface{}
			if record, err = client.GetNext(c.baseURL, response, &curRecord, params.Limit, c.getNextListQueryResponse); record == nil || err != nil {
				close(channel)
				return
			}

			response = record.(*ListQueryResponse)
		}
		close(channel)
	}()

	return channel, err
}

func (c *ApiService) getNextListQueryResponse(nextPageUrl string) (interface{}, error) {
	if nextPageUrl == "" {
		return nil, nil
	}
	resp, err := c.requestHandler.Get(nextPageUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ListQueryResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}
	return ps, nil
}

// Optional parameters for the method 'UpdateQuery'
type UpdateQueryParams struct {
	// The SID of an optional reference to the [Sample](https://www.twilio.com/docs/autopilot/api/task-sample) created from the query.
	SampleSid *string `json:"SampleSid,omitempty"`
	// The new status of the resource. Can be: `pending-review`, `reviewed`, or `discarded`
	Status *string `json:"Status,omitempty"`
}

func (params *UpdateQueryParams) SetSampleSid(SampleSid string) *UpdateQueryParams {
	params.SampleSid = &SampleSid
	return params
}
func (params *UpdateQueryParams) SetStatus(Status string) *UpdateQueryParams {
	params.Status = &Status
	return params
}

func (c *ApiService) UpdateQuery(AssistantSid string, Sid string, params *UpdateQueryParams) (*AutopilotV1Query, error) {
	path := "/v1/Assistants/{AssistantSid}/Queries/{Sid}"
	path = strings.Replace(path, "{"+"AssistantSid"+"}", AssistantSid, -1)
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.SampleSid != nil {
		data.Set("SampleSid", *params.SampleSid)
	}
	if params != nil && params.Status != nil {
		data.Set("Status", *params.Status)
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &AutopilotV1Query{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}