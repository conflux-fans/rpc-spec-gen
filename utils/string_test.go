package utils

import (
	"fmt"
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

func TestCleanComment(t *testing.T) {
	c := "Block representation\n#[derive(Debug, Serialize)]\n#[serde(rename_all = \"camelCase\")]"
	c = CleanComment(c)
	fmt.Printf("%s\n", c)
}
