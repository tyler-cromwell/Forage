/*
 * Twilio - Voice
 *
 * This is the public Twilio REST API.
 *
 * API version: 1.22.0
 * Contact: support@twilio.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

// ListDialingPermissionsCountryResponse struct for ListDialingPermissionsCountryResponse
type ListDialingPermissionsCountryResponse struct {
	Content []VoiceV1DialingPermissionsCountry `json:"content,omitempty"`
	Meta    ListByocTrunkResponseMeta          `json:"meta,omitempty"`
}