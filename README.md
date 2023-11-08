# asyncsqs

[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go)](https://pkg.go.dev/github.com/ngoyal16/asyncsqs?tab=doc)
[![Build, Unit Tests, Linters Status](https://github.com/ngoyal16/asyncsqs/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/ngoyal16/asyncsqs/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ngoyal16/asyncsqs/branch/master/graph/badge.svg)](https://codecov.io/gh/ngoyal16/asyncsqs)
[![Go Report Card](https://goreportcard.com/badge/github.com/ngoyal16/asyncsqs?clear_cache=2)](https://goreportcard.com/report/github.com/ngoyal16/asyncsqs)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

asyncsqs wraps around [SQS client](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sqs#Client)
from [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) to provide an async
buffered client which batches send message and delete message requests to
**optimise AWS costs**.

Messages can be scheduled to be sent and deleted. Requests will be dispatched
when

* either batch becomes full
* or waiting period exhausts (if configured)
* or the batch total body size becomes grater than or equal to 256 kb. (addition)

...**whichever occurs earlier**.

## Getting started

###### Add dependency

asyncsqs requires a Go version with [modules](https://github.com/golang/go/wiki/Modules)
support. If you're starting a new project, make sure to initialise a Go module:

```sh
$ mkdir ~/hellosqs
$ cd ~/hellosqs
$ go mod init github.com/my/hellosqs
```

And then add asyncsqs as a dependency to your existing or new project:

```sh
$ go get github.com/ngoyal16/asyncsqs
```

###### Write Code

please follow the demo code in the example folder