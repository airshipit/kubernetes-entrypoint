//go:build tools
// +build tools

package tools

import (
	// These imports are all tools used in the building and testing process
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
