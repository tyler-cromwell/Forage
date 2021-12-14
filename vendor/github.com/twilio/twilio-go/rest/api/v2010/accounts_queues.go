/*
 * Twilio - Api
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

// Optional parameters for the method 'CreateQueue'
type CreateQueueParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that will create the resource.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
	// A descriptive string that you created to describe this resource. It can be up to 64 characters long.
	FriendlyName *string `json:"FriendlyName,omitempty"`
	// The maximum number of calls allowed to be in the queue. The default is 100. The maximum is 5000.
	MaxSize *int `json:"MaxSize,omitempty"`
}

func (params *CreateQueueParams) SetPathAccountSid(PathAccountSid string) *CreateQueueParams {
	params.PathAccountSid = &PathAccountSid
	return params
}
func (params *CreateQueueParams) SetFriendlyName(FriendlyName string) *CreateQueueParams {
	params.FriendlyName = &FriendlyName
	return params
}
func (params *CreateQueueParams) SetMaxSize(MaxSize int) *CreateQueueParams {
	params.MaxSize = &MaxSize
	return params
}

// Create a queue
func (c *ApiService) CreateQueue(params *CreateQueueParams) (*ApiV2010Queue, error) {
	path := "/2010-04-01/Accounts/{AccountSid}/Queues.json"
	if params != nil && params.PathAccountSid != nil {
		path = strings.Replace(path, "{"+"AccountSid"+"}", *params.PathAccountSid, -1)
	} else {
		path = strings.Replace(path, "{"+"AccountSid"+"}", c.requestHandler.Client.AccountSid(), -1)
	}

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.FriendlyName != nil {
		data.Set("FriendlyName", *params.FriendlyName)
	}
	if params != nil && params.MaxSize != nil {
		data.Set("MaxSize", fmt.Sprint(*params.MaxSize))
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ApiV2010Queue{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Optional parameters for the method 'DeleteQueue'
type DeleteQueueParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Queue resource to delete.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
}

func (params *DeleteQueueParams) SetPathAccountSid(PathAccountSid string) *DeleteQueueParams {
	params.PathAccountSid = &PathAccountSid
	return params
}

// Remove an empty queue
func (c *ApiService) DeleteQueue(Sid string, params *DeleteQueueParams) error {
	path := "/2010-04-01/Accounts/{AccountSid}/Queues/{Sid}.json"
	if params != nil && params.PathAccountSid != nil {
		path = strings.Replace(path, "{"+"AccountSid"+"}", *params.PathAccountSid, -1)
	} else {
		path = strings.Replace(path, "{"+"AccountSid"+"}", c.requestHandler.Client.AccountSid(), -1)
	}
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

// Optional parameters for the method 'FetchQueue'
type FetchQueueParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Queue resource to fetch.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
}

func (params *FetchQueueParams) SetPathAccountSid(PathAccountSid string) *FetchQueueParams {
	params.PathAccountSid = &PathAccountSid
	return params
}

// Fetch an instance of a queue identified by the QueueSid
func (c *ApiService) FetchQueue(Sid string, params *FetchQueueParams) (*ApiV2010Queue, error) {
	path := "/2010-04-01/Accounts/{AccountSid}/Queues/{Sid}.json"
	if params != nil && params.PathAccountSid != nil {
		path = strings.Replace(path, "{"+"AccountSid"+"}", *params.PathAccountSid, -1)
	} else {
		path = strings.Replace(path, "{"+"AccountSid"+"}", c.requestHandler.Client.AccountSid(), -1)
	}
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ApiV2010Queue{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Optional parameters for the method 'ListQueue'
type ListQueueParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Queue resources to read.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
	// How many resources to return in each list page. The default is 50, and the maximum is 1000.
	PageSize *int `json:"PageSize,omitempty"`
	// Max number of records to return.
	Limit *int `json:"limit,omitempty"`
}

func (params *ListQueueParams) SetPathAccountSid(PathAccountSid string) *ListQueueParams {
	params.PathAccountSid = &PathAccountSid
	return params
}
func (params *ListQueueParams) SetPageSize(PageSize int) *ListQueueParams {
	params.PageSize = &PageSize
	return params
}
func (params *ListQueueParams) SetLimit(Limit int) *ListQueueParams {
	params.Limit = &Limit
	return params
}

// Retrieve a single page of Queue records from the API. Request is executed immediately.
func (c *ApiService) PageQueue(params *ListQueueParams, pageToken, pageNumber string) (*ListQueueResponse, error) {
	path := "/2010-04-01/Accounts/{AccountSid}/Queues.json"

	if params != nil && params.PathAccountSid != nil {
		path = strings.Replace(path, "{"+"AccountSid"+"}", *params.PathAccountSid, -1)
	} else {
		path = strings.Replace(path, "{"+"AccountSid"+"}", c.requestHandler.Client.AccountSid(), -1)
	}

	data := url.Values{}
	headers := make(map[string]interface{})

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

	ps := &ListQueueResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Lists Queue records from the API as a list. Unlike stream, this operation is eager and loads 'limit' records into memory before returning.
func (c *ApiService) ListQueue(params *ListQueueParams) ([]ApiV2010Queue, error) {
	if params == nil {
		params = &ListQueueParams{}
	}
	params.SetPageSize(client.ReadLimits(params.PageSize, params.Limit))

	response, err := c.PageQueue(params, "", "")
	if err != nil {
		return nil, err
	}

	curRecord := 0
	var records []ApiV2010Queue

	for response != nil {
		records = append(records, response.Queues...)

		var record interface{}
		if record, err = client.GetNext(c.baseURL, response, &curRecord, params.Limit, c.getNextListQueueResponse); record == nil || err != nil {
			return records, err
		}

		response = record.(*ListQueueResponse)
	}

	return records, err
}

// Streams Queue records from the API as a channel stream. This operation lazily loads records as efficiently as possible until the limit is reached.
func (c *ApiService) StreamQueue(params *ListQueueParams) (chan ApiV2010Queue, error) {
	if params == nil {
		params = &ListQueueParams{}
	}
	params.SetPageSize(client.ReadLimits(params.PageSize, params.Limit))

	response, err := c.PageQueue(params, "", "")
	if err != nil {
		return nil, err
	}

	curRecord := 0
	//set buffer size of the channel to 1
	channel := make(chan ApiV2010Queue, 1)

	go func() {
		for response != nil {
			for item := range response.Queues {
				channel <- response.Queues[item]
			}

			var record interface{}
			if record, err = client.GetNext(c.baseURL, response, &curRecord, params.Limit, c.getNextListQueueResponse); record == nil || err != nil {
				close(channel)
				return
			}

			response = record.(*ListQueueResponse)
		}
		close(channel)
	}()

	return channel, err
}

func (c *ApiService) getNextListQueueResponse(nextPageUrl string) (interface{}, error) {
	if nextPageUrl == "" {
		return nil, nil
	}
	resp, err := c.requestHandler.Get(nextPageUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ListQueueResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}
	return ps, nil
}

// Optional parameters for the method 'UpdateQueue'
type UpdateQueueParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Queue resource to update.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
	// A descriptive string that you created to describe this resource. It can be up to 64 characters long.
	FriendlyName *string `json:"FriendlyName,omitempty"`
	// The maximum number of calls allowed to be in the queue. The default is 100. The maximum is 5000.
	MaxSize *int `json:"MaxSize,omitempty"`
}

func (params *UpdateQueueParams) SetPathAccountSid(PathAccountSid string) *UpdateQueueParams {
	params.PathAccountSid = &PathAccountSid
	return params
}
func (params *UpdateQueueParams) SetFriendlyName(FriendlyName string) *UpdateQueueParams {
	params.FriendlyName = &FriendlyName
	return params
}
func (params *UpdateQueueParams) SetMaxSize(MaxSize int) *UpdateQueueParams {
	params.MaxSize = &MaxSize
	return params
}

// Update the queue with the new parameters
func (c *ApiService) UpdateQueue(Sid string, params *UpdateQueueParams) (*ApiV2010Queue, error) {
	path := "/2010-04-01/Accounts/{AccountSid}/Queues/{Sid}.json"
	if params != nil && params.PathAccountSid != nil {
		path = strings.Replace(path, "{"+"AccountSid"+"}", *params.PathAccountSid, -1)
	} else {
		path = strings.Replace(path, "{"+"AccountSid"+"}", c.requestHandler.Client.AccountSid(), -1)
	}
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.FriendlyName != nil {
		data.Set("FriendlyName", *params.FriendlyName)
	}
	if params != nil && params.MaxSize != nil {
		data.Set("MaxSize", fmt.Sprint(*params.MaxSize))
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ApiV2010Queue{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}
