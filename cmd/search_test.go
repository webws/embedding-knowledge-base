package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webws/embedding-knowledge-base/options"
)

func TestNewSearchCmd(t *testing.T) {
	// first, you must start the vector database qdrant
	// docker run --rm -p 6334:6334 qdrant/qdrant
	c := options.NewConfigFlags()
	c.ApiKey = "you apiKey" // or env set apiKey
	scmd := NewSearchCmd(*c)
	scmd.Flags().Set(flagmsg, "k8s")
	err := scmd.RunE(scmd, []string{})
	assert.NoError(t, err)
}
