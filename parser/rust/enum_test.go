package rust

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseEnum(t *testing.T) {
	var e Enum = `
pub enum BlockNumber {
	/// Number
	Num(U64),
	/// Earliest block (true genesis)
	Earliest,
	/// The latest committed
	LatestCommitted,
	/// The latest voted
	LatestVoted,
}`

	e = `#[derive(Debug, PartialEq, Clone, Hash, Eq)]
pub enum BlockHashOrEpochNumber {
    BlockHash(H256),
    EpochNumber(EpochNumber),
}`

	r := e.Parse()
	j, _ := json.MarshalIndent(r, "", "  ")
	fmt.Printf("%s\n", j)
}
