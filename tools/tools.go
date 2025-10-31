//go:build tools
// +build tools

package tools

import (
	_ "github.com/daixiang0/gci"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "mvdan.cc/gofumpt"
)
