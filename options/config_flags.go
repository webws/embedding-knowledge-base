package options

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

const (
	DefaultQdrantAddr = "127.0.0.1:6334"
	DefaultSocksProxy = "socks5://127.0.0.1:1080"
)

var (
	flagQdrant = "qdrant" // qdrant address
	// flagDataFile = "data_file" // cmd import flag
	// msgFlag      = "msg"       // cmd ask flag  question
	flagApiKey     = "apiKey" // open apikey
	flagProxy      = "proxy"  // openai http proxy
	flagCollection = "collection"
	flagVectorSize = "vectorSize"
)

type ConfigFlags struct {
	Qdrant     string
	ApiKey     string
	Proxy      string
	Collection string
	VectorSize uint64
}

// NewConfigFlags
func NewConfigFlags() *ConfigFlags {
	return &ConfigFlags{
		Qdrant:     DefaultQdrantAddr,
		ApiKey:     "",
		Proxy:      DefaultSocksProxy,
		Collection: "kubernetes",
		VectorSize: 1536,
	}
}

// AddFlags binds client configuration flags to a given flagset
func (cf *ConfigFlags) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cf.Qdrant, flagQdrant, DefaultQdrantAddr, fmt.Sprintf("qdrant address default: %s", DefaultQdrantAddr))
	fs.StringVar(&cf.Proxy, flagProxy, DefaultSocksProxy, fmt.Sprintf("http client proxy default:%s ", DefaultSocksProxy))
	apiKey := os.Getenv(flagApiKey)
	fs.StringVar(&cf.ApiKey, flagApiKey, apiKey, "openai apikey:default from env "+flagApiKey)
	fs.StringVar(&cf.Collection, flagCollection, cf.Collection, "qdrant collection name default: "+cf.Collection)
	fs.Uint64Var(&cf.VectorSize, flagVectorSize, cf.VectorSize, "qdrant vector size default: "+fmt.Sprintf("%d", cf.VectorSize))
}
