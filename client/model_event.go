/*
 * Paygate API
 *
 * Paygate is a RESTful API enabling Automated Clearing House ([ACH](https://en.wikipedia.org/wiki/Automated_Clearing_House)) transactions to be submitted and received without a deep understanding of a full NACHA file specification.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"time"
)

// Event struct for Event
type Event struct {
	// ID to uniquely identify a event
	ID string `json:"ID,omitempty"`
	// Type of event
	Topic string `json:"topic,omitempty"`
	// A human readable description of the topic
	Message string `json:"message,omitempty"`
	Type    string `json:"type,omitempty"`
	// ID of the resource type the event was generated on behalf of.
	Resource string    `json:"resource,omitempty"`
	Created  time.Time `json:"created,omitempty"`
}