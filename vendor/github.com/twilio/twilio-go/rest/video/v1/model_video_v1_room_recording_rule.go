/*
 * Twilio - Video
 *
 * This is the public Twilio REST API.
 *
 * API version: 1.22.0
 * Contact: support@twilio.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"time"
)

// VideoV1RoomRecordingRule struct for VideoV1RoomRecordingRule
type VideoV1RoomRecordingRule struct {
	// The ISO 8601 date and time in GMT when the resource was created
	DateCreated *time.Time `json:"date_created,omitempty"`
	// The ISO 8601 date and time in GMT when the resource was last updated
	DateUpdated *time.Time `json:"date_updated,omitempty"`
	// The SID of the Room resource for the Recording Rules
	RoomSid *string `json:"room_sid,omitempty"`
	// A collection of recording Rules that describe how to include or exclude matching tracks for recording
	Rules *[]VideoV1RoomRoomRecordingRuleRules `json:"rules,omitempty"`
}
