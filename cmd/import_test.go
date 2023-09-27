package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webws/embedding-knowledge-base/options"
)

func TestNewImportCmd(t *testing.T) {
	// first, you must start the vector database qdrant
	// docker run --rm -p 6334:6334 qdrant/qdrant
	c := options.NewConfigFlags()
	c.ApiKey = "you apiKey" // or env set apiKey
	cmd := NewImportCmd(*c)
	cmd.Flags().Set(FlagdataFile, "../example/data.json")
	err := cmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
