package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if version == "" {
		t.Error("Version should not be empty")
	}
}
