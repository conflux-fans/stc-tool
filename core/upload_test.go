package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zero-gravity-labs/zerog-storage-tool/config"
)

func TestUploadStream(t *testing.T) {
	config.SetConfigFile("/Users/dayong/myspace/mywork/zerog-storage-tool/config.yaml")
	config.Init()

	Init()
	// put
	batcher := kvClientForPut.Batcher()
	batcher.Set(STREAM_FILE,
		[]byte("KEY0"),
		[]byte("hello world"),
	)

	err := batcher.Exec()
	assert.NoError(t, err)

	// query
	iter := kvClientForIterator.NewIterator(STREAM_FILE)
	err = iter.SeekToFirst()
	assert.NoError(t, err)

	assert.True(t, iter.Valid())
}
