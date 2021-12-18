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

// Optional parameters for the method 'CreateAccount'
type CreateAccountParams struct {
	// A human readable description of the account to create, defaults to `SubAccount Created at {YYYY-MM-DD HH:MM meridian}`
	FriendlyName *string `json:"FriendlyName,omitempty"`
}

func (params *CreateAccountParams) SetFriendlyName(FriendlyName string) *CreateAccountParams {
	params.FriendlyName = &FriendlyName
	return params
}

// Create a new Twilio Subaccount from the account making the request
func (c *ApiService) CreateAccount(params *CreateAccountParams) (*ApiV2010Account, error) {
	path := "/2010-04-01/Accounts.json"

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.FriendlyName != nil {
		data.Set("FriendlyName", *params.FriendlyName)
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ApiV2010Account{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Fetch the account specified by the provided Account Sid
func (c *ApiService) FetchAccount(Sid string) (*ApiV2010Account, error) {
	path := "/2010-04-01/Accounts/{Sid}.json"
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ApiV2010Account{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Optional parameters for the method 'ListAccount'
type ListAccountParams struct {
	// Only return the Account resources with friendly names that exactly match this name.
	FriendlyName *string `json:"FriendlyName,omitempty"`
	// Only return Account resources with the given status. Can be `closed`, `suspended` or `active`.
	Status *string `json:"Status,omitempty"`
	// How many resources to return in each list page. The default is 50, and the maximum is 1000.
	PageSize *int `json:"PageSize,omitempty"`
	// Max number of records to return.
	Limit *int `json:"limit,omitempty"`
}

func (params *ListAccountParams) SetFriendlyName(FriendlyName string) *ListAccountParams {
	params.FriendlyName = &FriendlyName
	return params
}
func (params *ListAccountParams) SetStatus(Status string) *ListAccountParams {
	params.Status = &Status
	return params
}
func (params *ListAccountParams) SetPageSize(PageSize int) *ListAccountParams {
	params.PageSize = &PageSize
	return params
}
func (params *ListAccountParams) SetLimit(Limit int) *ListAccountParams {
	params.Limit = &Limit
	return params
}

// Retrieve a single page of Account records from the API. Request is executed immediately.
func (c *ApiService) PageAccount(params *ListAccountParams, pageToken, pageNumber string) (*ListAccountResponse, error) {
	path := "/2010-04-01/Accounts.json"

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.FriendlyName != nil {
		data.Set("FriendlyName", *params.FriendlyName)
	}
	if params != nil && params.Status != nil {
		data.Set("Status", *params.Status)
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

	ps := &ListAccountResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Lists Account records from the API as a list. Unlike stream, this operation is eager and loads 'limit' records into memory before returning.
func (c *ApiService) ListAccount(params *ListAccountParams) ([]ApiV2010Account, error) {
	if params == nil {
		params = &ListAccountParams{}
	}
	params.SetPageSize(client.ReadLimits(params.PageSize, params.Limit))

	response, err := c.PageAccount(params, "", "")
	if err != nil {
		return nil, err
	}

	curRecord := 0
	var records []ApiV2010Account

	for response != nil {
		records = append(records, response.Accounts...)

		var record interface{}
		if record, err = client.GetNext(c.baseURL, response, &curRecord, params.Limit, c.getNextListAccountResponse); record == nil || err != nil {
			return records, err
		}

		response = record.(*ListAccountResponse)
	}

	return records, err
}

// Streams Account records from the API as a channel stream. This operation lazily loads records as efficiently as possible until the limit is reached.
func (c *ApiService) StreamAccount(params *ListAccountParams) (chan ApiV2010Account, error) {
	if params == nil {
		params = &ListAccountParams{}
	}
	params.SetPageSize(client.ReadLimits(params.PageSize, params.Limit))

	response, err := c.PageAccount(params, "", "")
	if err != nil {
		return nil, err
	}

	curRecord := 0
	//set buffer size of the channel to 1
	channel := make(chan ApiV2010Account, 1)

	go func() {
		for response != nil {
			for item := range response.Accounts {
				channel <- response.Accounts[item]
			}

			var record interface{}
			if record, err = client.GetNext(c.baseURL, response, &curRecord, params.Limit, c.getNextListAccountResponse); record == nil || err != nil {
				close(channel)
				return
			}

			response = record.(*ListAccountResponse)
		}
		close(channel)
	}()

	return channel, err
}

func (c *ApiService) getNextListAccountResponse(nextPageUrl string) (interface{}, error) {
	if nextPageUrl == "" {
		return nil, nil
	}
	resp, err := c.requestHandler.Get(nextPageUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ListAccountResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}
	return ps, nil
}

// Optional parameters for the method 'UpdateAccount'
type UpdateAccountParams struct {
	// Update the human-readable description of this Account
	FriendlyName *string `json:"FriendlyName,omitempty"`
	// Alter the status of this account: use `closed` to irreversibly close this account, `suspended` to temporarily suspend it, or `active` to reactivate it.
	Status *string `json:"Status,omitempty"`
}

func (params *UpdateAccountParams) SetFriendlyName(FriendlyName string) *UpdateAccountParams {
	params.FriendlyName = &FriendlyName
	return params
}
func (params *UpdateAccountParams) SetStatus(Status string) *UpdateAccountParams {
	params.Status = &Status
	return params
}

// Modify the properties of a given Account
func (c *ApiService) UpdateAccount(Sid string, params *UpdateAccountParams) (*ApiV2010Account, error) {
	path := "/2010-04-01/Accounts/{Sid}.json"
	path = strings.Replace(path, "{"+"Sid"+"}", Sid, -1)

	data := url.Values{}
	headers := make(map[string]interface{})

	if params != nil && params.FriendlyName != nil {
		data.Set("FriendlyName", *params.FriendlyName)
	}
	if params != nil && params.Status != nil {
		data.Set("Status", *params.Status)
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ApiV2010Account{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}