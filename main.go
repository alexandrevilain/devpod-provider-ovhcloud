package main

import (
	"math/rand"
	"time"

	"github.com/alexandrevilain/devpod-provider-ovhcloud/cmd"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //nolint:staticcheck
	cmd.Execute()
}
