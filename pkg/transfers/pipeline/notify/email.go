// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package notify

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/moov-io/ach"
	"github.com/moov-io/paygate/pkg/config"

	gomail "github.com/ory/mail/v3"
)

type Email struct {
	cfg    *config.Email
	dialer *gomail.Dialer
}

type EmailTemplateData struct {
	CompanyName string // e.g. Moov
	Verb        string // e.g. uploaded, downloaded
	Filename    string // e.g. 20200529-131400.ach

	DebitTotal  int
	CreditTotal int

	BatchCount int
	EntryCount int
}

var (
	// Ensure the default template validates against our data struct
	_ = config.DefaultEmailTemplate.Execute(ioutil.Discard, EmailTemplateData{})
)

func NewEmail(cfg *config.Email) (*Email, error) {
	dialer, err := setupGoMailClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Email{
		cfg:    cfg,
		dialer: dialer,
	}, nil
}

func setupGoMailClient(cfg *config.Email) (*gomail.Dialer, error) {
	uri, err := url.Parse(cfg.ConnectionURI)
	if err != nil {
		return nil, err
	}
	password, _ := uri.User.Password()
	port, _ := strconv.ParseInt(uri.Port(), 10, 64)

	tlsConfig := &tls.Config{}

	ssl := strings.EqualFold(uri.Scheme, "smtps")
	if ssl {
		sslSkipVerify, _ := strconv.ParseBool(uri.Query().Get("skip_ssl_verify"))
		tlsConfig = &tls.Config{InsecureSkipVerify: sslSkipVerify}
	}

	return &gomail.Dialer{
		TLSConfig:    tlsConfig,
		Host:         uri.Hostname(),
		Port:         int(port),
		Username:     uri.User.Username(),
		Password:     password,
		SSL:          ssl,
		Timeout:      time.Second * 10,
		RetryFailure: true,
	}, nil
}

func (mailer *Email) Info(msg *Message) error {
	contents, err := marshalEmail(mailer.cfg, msg)
	if err != nil {
		return err
	}
	return sendEmail(mailer.cfg, mailer.dialer, msg.Filename, contents)
}

func (mailer *Email) Critical(msg *Message) error {
	contents, err := marshalEmail(mailer.cfg, msg)
	if err != nil {
		return err
	}
	return sendEmail(mailer.cfg, mailer.dialer, msg.Filename, contents)
}

func marshalEmail(cfg *config.Email, msg *Message) (string, error) {
	data := EmailTemplateData{
		CompanyName: cfg.CompanyName,
		Verb:        string(msg.Direction),
		Filename:    msg.Filename,
		DebitTotal:  msg.File.Control.TotalDebitEntryDollarAmountInFile,
		CreditTotal: msg.File.Control.TotalCreditEntryDollarAmountInFile,
		BatchCount:  msg.File.Control.BatchCount,
		EntryCount:  countEntries(msg.File),
	}

	var buf bytes.Buffer
	if err := cfg.Tmpl().Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func sendEmail(cfg *config.Email, dialer *gomail.Dialer, filename, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", cfg.To...)
	m.SetHeader("Subject", fmt.Sprintf("%s uploaded by %s", filename, cfg.CompanyName))
	m.SetBody("text/plain", body)

	if err := dialer.DialAndSend(context.Background(), m); err != nil {
		fmt.Printf("m=%#v\n", m)
		fmt.Printf("dialer=%#v\n", dialer)
		return err
	}
	return nil
}

func countEntries(file *ach.File) int {
	var total int
	for i := range file.Batches {
		total += len(file.Batches[i].GetEntries())
	}
	return total
}
