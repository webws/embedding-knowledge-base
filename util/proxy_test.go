package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyTest(t *testing.T) {
	b := ProxyCheck(nil, "socks5://127.0.0.1:1080", "https://google.com", 10)
	fmt.Println(b)
	assert.True(t, b)
	// assert.True(t, ProxyCheck(nil, "socks5://127.0.0.1:1080", "https://google.com", 5))
}
