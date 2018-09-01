package controllers

import (
	"testing"
)

func TestDataCaching(t *testing.T) {
	fetch()

	if cache.Time == "" {
		t.Error("Data caching has failed")
	}
}
