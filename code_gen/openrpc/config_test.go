package openrpc

import (
	"fmt"
	"testing"
)

func TestRelateivePath(t *testing.T) {
	fmt.Println(getSchemaSaveRelativePath("cfx_space", "a::b::c"))
}

func TestGetSchemaSavePath(t *testing.T) {
	fmt.Println(getSchemaSavePath("cfx_space", "a::b::c"))
}
