/*
 * Paygate API
 *
 * PayGate is a RESTful API enabling first-party Automated Clearing House ([ACH](https://en.wikipedia.org/wiki/Automated_Clearing_House)) transfers to be created without a deep understanding of a full NACHA file specification. First-party transfers initiate at an Originating Depository Financial Institution (ODFI) and are sent off to other Financial Institutions.  Tenants are the largest grouping in PayGate and are typically a vendor who is reselling ACH services or a company making ACH payments themselves. A legal entity is linked off a Tenant as the primary Customer used to KYC and in transfers with the Tenant itself.  An Organization is a grouping within a Tenant which typically represents an entity making ACH transfers. These include clients of an ACH reseller or business accepting payments over ACH. A legal entity is linked off an Organization as the primary Customer used to KYC and in transfers with the Organization itself.  ![](https://raw.githubusercontent.com/moov-io/paygate/master/docs/images/tenant-in-paygate.png)
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package client

import (
	"time"
)

// MicroDeposits struct for MicroDeposits
type MicroDeposits struct {
	// A microDepositID to identify this set of credits to an external account
	MicroDepositID string `json:"microDepositID,omitempty"`
	// An array of transferID values created from this micro-deposit
	TransferIDs []string       `json:"transferIDs,omitempty"`
	Destination Destination    `json:"destination,omitempty"`
	Amounts     []string       `json:"amounts,omitempty"`
	Status      TransferStatus `json:"status,omitempty"`
	ReturnCode  ReturnCode     `json:"returnCode,omitempty"`
	Created     time.Time      `json:"created,omitempty"`
}
