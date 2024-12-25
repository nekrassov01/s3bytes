package main

import (
	"context"

	"github.com/nekrassov01/s3bytes"
)

func main() {
	s3bytes.CLI(context.Background())
}
