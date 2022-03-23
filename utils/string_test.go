package utils

import (
	"fmt"
	"regexp"
	"testing"
)

func TestCC2USC(t *testing.T) {
	s := CamelCase2UnderScoreCase("BlockNumber")
	if s != "block_number" {
		t.Errorf("expected block, got %s", s)
	}
}

func TestUSC2CC(t *testing.T) {
	s := UnderScoreCase2CamelCase("block_number", true)
	if s != "BlockNumber" {
		t.Errorf("expected block, got %s", s)
	}
	s = UnderScoreCase2CamelCase("block_number", false)
	if s != "blockNumber" {
		t.Errorf("expected block, got %s", s)
	}
	s = UnderScoreCase2CamelCase("H256", false)
	if s != "h256" {
		t.Errorf("expected h256, got %s", s)
	}
	s = UnderScoreCase2CamelCase("H256", true)
	if s != "H256" {
		t.Errorf("expected H256, got %s", s)
	}
}

func TestA(t *testing.T) {
	sampleRegexp := regexp.MustCompile(`(?P\w+):(?P[0-9]\d{1,3})`)

	input := "The names are John:21, Simon:23, Mike:19"

	result := sampleRegexp.ReplaceAllString(input, "$Age:$Name")
	fmt.Println(string(result))
}
