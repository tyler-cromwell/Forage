/*
 * Twilio - Trusthub
 *
 * This is the public Twilio REST API.
 *
 * API version: 1.22.0
 * Contact: support@twilio.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

// ListTrustProductChannelEndpointAssignmentResponse struct for ListTrustProductChannelEndpointAssignmentResponse
type ListTrustProductChannelEndpointAssignmentResponse struct {
	Meta    ListCustomerProfileResponseMeta                   `json:"meta,omitempty"`
	Results []TrusthubV1TrustProductChannelEndpointAssignment `json:"results,omitempty"`
}