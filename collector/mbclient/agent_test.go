package mbclient

import (
	"testing"

	"github.com/mythay/spider"

	"github.com/stretchr/testify/assert"
)

func Test_mbRoot_Once(t *testing.T) {
	assert := assert.New(t)

	config, err := spider.LoadCfgModbus("testtcp.toml")
	assert.Empty(err)
	root, err := newRoot(*config)
	assert.Empty(err)
	root.Once()

}

func Benchmark_mbTcp(b *testing.B) {
	assert := assert.New(b)

	b.StopTimer()
	config, err := spider.LoadCfgModbus("testtcp.toml")
	assert.Empty(err)
	root, err := newRoot(*config)

	assert.Empty(err)
	defer root.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		root.Once()
	}

}
