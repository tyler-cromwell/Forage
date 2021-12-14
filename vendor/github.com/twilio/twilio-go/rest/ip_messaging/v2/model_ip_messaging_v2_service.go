/*
 * Twilio - Ip_messaging
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

// IpMessagingV2Service struct for IpMessagingV2Service
type IpMessagingV2Service struct {
	AccountSid                   *string                 `json:"account_sid,omitempty"`
	ConsumptionReportInterval    *int                    `json:"consumption_report_interval,omitempty"`
	DateCreated                  *time.Time              `json:"date_created,omitempty"`
	DateUpdated                  *time.Time              `json:"date_updated,omitempty"`
	DefaultChannelCreatorRoleSid *string                 `json:"default_channel_creator_role_sid,omitempty"`
	DefaultChannelRoleSid        *string                 `json:"default_channel_role_sid,omitempty"`
	DefaultServiceRoleSid        *string                 `json:"default_service_role_sid,omitempty"`
	FriendlyName                 *string                 `json:"friendly_name,omitempty"`
	Limits                       *map[string]interface{} `json:"limits,omitempty"`
	Links                        *map[string]interface{} `json:"links,omitempty"`
	Media                        *map[string]interface{} `json:"media,omitempty"`
	Notifications                *map[string]interface{} `json:"notifications,omitempty"`
	PostWebhookRetryCount        *int                    `json:"post_webhook_retry_count,omitempty"`
	PostWebhookUrl               *string                 `json:"post_webhook_url,omitempty"`
	PreWebhookRetryCount         *int                    `json:"pre_webhook_retry_count,omitempty"`
	PreWebhookUrl                *string                 `json:"pre_webhook_url,omitempty"`
	ReachabilityEnabled          *bool                   `json:"reachability_enabled,omitempty"`
	ReadStatusEnabled            *bool                   `json:"read_status_enabled,omitempty"`
	Sid                          *string                 `json:"sid,omitempty"`
	TypingIndicatorTimeout       *int                    `json:"typing_indicator_timeout,omitempty"`
	Url                          *string                 `json:"url,omitempty"`
	WebhookFilters               *[]string               `json:"webhook_filters,omitempty"`
	WebhookMethod                *string                 `json:"webhook_method,omitempty"`
}
