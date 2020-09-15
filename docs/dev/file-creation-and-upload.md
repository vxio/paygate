## How are ACH files uploaded?

Router, fundflow
needs source, destination, originate, save trace numbers, publish
https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/router.go#L135-L213

https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/fundflow/first_party.go#L43

https://github.com/moov-io/paygate/tree/v0.8.0/pkg/transfers/pipeline

https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/pipeline/publisher.go#L38

https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/pipeline/publisher_stream.go#L23

<-await() into handleMessage (reads Xfer or Cancel)
https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/pipeline/aggregate.go#L267

handoff to `filesystemMerging` which writes file to disk
later at each cutoff window on `WithEachFile` to merge all files together
https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/pipeline/merging.go

`runTransformers` to transform (encrypt, etc) and finally upload
https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/pipeline/aggregate.go#L170-L176

output: https://github.com/moov-io/paygate/tree/v0.8.0/pkg/transfers/pipeline/output
encrypt: https://github.com/moov-io/paygate/tree/v0.8.0/pkg/transfers/pipeline/transform

`uploadFile`: (render filename, format, audit storage, FTP/SFTP upload)
https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/pipeline/aggregate.go#L211

`notifyAfterUpload` (slack, email, PD) on Info or Critical
https://github.com/moov-io/paygate/blob/v0.8.0/pkg/transfers/pipeline/aggregate.go#L247


---
(to Pavel in Slack)

it starts with these two models
11:27
https://github.com/moov-io/paygate/blob/master/pkg/transfers/pipeline/pipeline.go

pkg/transfers/pipeline/pipeline.go
// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package pipeline
Show more
<https://github.com/moov-io/paygate|moov-io/paygate>moov-io/paygate | Added by GitHub
11:27
in the POST /transfers route the happy result is to create an ach.File and publish it to the “pipeline” here
11:28
so, client.Transfer is the stored DB model, and ach.File is what we intend to send off to the ODFI
11:28
an Xfer gets published - https://github.com/moov-io/paygate/blob/master/pkg/transfers/pipeline/publisher.go (edited)
11:28
(there are a few different implementations, by default it’s inmem, but could be kafka)
11:29
https://github.com/moov-io/paygate/blob/master/pkg/transfers/pipeline/aggregate.go#L45

pkg/transfers/pipeline/aggregate.go:45
    subscription *pubsub.Subscription
<https://github.com/moov-io/paygate|moov-io/paygate>moov-io/paygate | Added by GitHub
11:29
XferAggregation takes the receiving end of that pubsub (by default inmem, could be kafka)
11:30
and with each incoming Xfer, writes it to local disk in preparation for merging and upload
11:30
https://github.com/moov-io/paygate/blob/master/pkg/transfers/pipeline/aggregate.go#L283-L287

pkg/transfers/pipeline/aggregate.go:283-287
    var xfer Xfer
    err := json.NewDecoder(bytes.NewReader(msg.Body)).Decode(&xfer)
    if err == nil && xfer.Transfer != nil && xfer.File != nil {
        // Handle the Xfer after decoding it.
        if err := merger.HandleXfer(xfer); err != nil {
Show more
<https://github.com/moov-io/paygate|moov-io/paygate>moov-io/paygate | Added by GitHub
11:31
https://github.com/moov-io/paygate/blob/master/pkg/transfers/pipeline/merging.go#L62

pkg/transfers/pipeline/merging.go:62
func (m *filesystemMerging) HandleXfer(xfer Xfer) error {
<https://github.com/moov-io/paygate|moov-io/paygate>moov-io/paygate | Added by GitHub
11:31
(I don’t expect this to make sense right now, just typing up the flow)
11:31
okay, so over time we collect N transfers that intend to go out by our cutoff window, say 4pm in NYC
11:33
(this is a bit sloppy, there’s no estimation done to start this process before 4pm - we’re going to configure it to run at 3:50pm for now)
11:33
but at that cutoff trigger we start processing each pending file
11:33
https://github.com/moov-io/paygate/blob/master/pkg/transfers/pipeline/aggregate.go#L135-L139

pkg/transfers/pipeline/aggregate.go:135-139
        case tt := <-cutoffs.C:
            if err := xfagg.processCutoffCallbacks(); err != nil {
                xfagg.logger.Log("aggregate", fmt.Sprintf("ERROR with cutoff callbacks: %v", err))
            }
            xfagg.withEachFile(tt)
Show more
<https://github.com/moov-io/paygate|moov-io/paygate>moov-io/paygate | Added by GitHub
11:33
those pending files get merged and finally uploaded to the ODFI
11:34
only at that point of upload do we “mark the transfers as processed”
  https://github.com/moov-io/paygate/blob/master/pkg/transfers/pipeline/aggregate.go#L200-L204

pkg/transfers/pipeline/aggregate.go:200-204
    if processed, err := xfagg.merger.WithEachMerged(xfagg.runTransformers); err != nil {
        xfagg.logger.Log("aggregate", fmt.Sprintf("ERROR inside WithEachMerged: %v", err))
    } else {
        if err := xfagg.repo.MarkTransfersAsProcessed(processed.transferIDs); err != nil {
            xfagg.logger.Log("aggregate", fmt.Sprintf("ERROR marking %d transfers as processed: %v", len(processed.transferIDs), err))
Show more
<https://github.com/moov-io/paygate|moov-io/paygate>moov-io/paygate | Added by GitHub
11:35
then we archive that uploaded file, notify ourselves (slack, email group, PagerDuty) and continue to the next file
11:35
easy, right? hahaha
