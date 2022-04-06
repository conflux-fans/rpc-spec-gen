package openrpc

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/conflux-fans/rpc-spec-gen/code_gen/openrpc/types"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func TestMarshalStructPtr(t *testing.T) {
	tag := types.Tag{}
	j, _ := json.MarshalIndent(tag, "", "  ")
	fmt.Printf("%s\n", j)

	m := types.Method{}
	j, _ = json.MarshalIndent(m, "", "  ")
	fmt.Printf("%s\n", j)
}

func TestLogrus(t *testing.T) {
	c := types.Contact{
		Name:  "SOPHIA",
		Email: "SOPHIA@163.com",
	}
	j, _ := json.MarshalIndent(c, "", "  ")

	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &prefixed.TextFormatter{
			// DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}

	logger.Printf("Log message")
	// logrus.WithField("content", string(j)).Info("demo content")
	logger.WithField("content", string(j)).Debug("demo content")
}
