package util

import (
	"fmt"
	"testing"
	"time"
)

func TestLongTimestamps(t *testing.T) {

	ts := time.Now().UTC()

	formattedTs := FormatLongTimeStamp(ts)

	fmt.Println("Formatted", formattedTs)

	if len(formattedTs) == 0 {
		t.Error("Empty result")
	}

	parsedTs, err := ParseLongTimeStamp(formattedTs, time.UTC)

	if err != nil {
		t.Error(err)
	}

	if !TimeClose(ts, parsedTs, time.Second) {
		t.Error("Comparison didn't match")
	}

}

func TestShortTimestamps(t *testing.T) {

	ts := time.Now()

	formattedTs := FormatShortTimeStamp(ts)

	fmt.Println("Formatted", formattedTs)

	if len(formattedTs) == 0 {
		t.Error("Empty result")
	}

	parsedTs, err := ParseShortTimeStamp(formattedTs, time.Local)

	if err != nil {
		t.Error(err)
	}

	if !TimeClose(ts, parsedTs, time.Second) {
		t.Error("Comparison didn't match")
	}

}
