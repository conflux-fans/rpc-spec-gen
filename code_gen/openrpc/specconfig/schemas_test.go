package specconfig

import (
	"fmt"
	"testing"
)

func TestRelateivePath(t *testing.T) {
	fmt.Println(GetSchemaSaveRelativePath("cfx_space", "a::b::c"))
}

func TestGetSchemaSavePath(t *testing.T) {
	fmt.Println(GetSchemaSavePath("cfx_space", "a::b::c"))
}
