package core

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/conflux-fans/storage-cli/config"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/stretchr/testify/assert"
)

func TestUploadStream(t *testing.T) {
	config.SetConfigFile("/Users/dayong/myspace/mywork/storage-cli/config.yaml")
	config.Init()

	Init()
	// put
	batcher := adminBatcher
	key := []byte(fmt.Sprintf("TEST-KEY-%d", time.Now().Unix()))

	logger.Get().WithField("key", string(key)).Info("Start put")

	batcher.Set(STREAM_FILE, key, []byte("hello world")).
		SetKeyToSpecial(STREAM_FILE, key).
		GrantSpecialWriteRole(STREAM_FILE, key, defaultAccount)

	_, err := batcher.Exec(context.Background())
	assert.NoError(t, err)

	// query
	time.Sleep(10 * time.Second)
	val, err := kvClientForIterator.GetValue(context.Background(), STREAM_FILE, key)
	assert.NoError(t, err)
	assert.True(t, val.Size > 0)

	isSpecialKey, err := kvClientForIterator.IsSpecialKey(context.Background(), STREAM_FILE, key)
	assert.NoError(t, err)
	assert.True(t, isSpecialKey)

	isWriter, err := kvClientForIterator.IsWriterOfKey(context.Background(), defaultAccount, STREAM_FILE, key)
	assert.NoError(t, err)
	assert.True(t, isWriter)
}
