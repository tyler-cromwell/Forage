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

// AutopilotV1Defaults struct for AutopilotV1Defaults
type AutopilotV1Defaults struct {
	// The SID of the Account that created the resource
	AccountSid *string `json:"account_sid,omitempty"`
	// The SID of the Assistant that is the parent of the resource
	AssistantSid *string `json:"assistant_sid,omitempty"`
	// The JSON string that describes the default task links
	Data *map[string]interface{} `json:"data,omitempty"`
	// The absolute URL of the Defaults resource
	Url *string `json:"url,omitempty"`
}
