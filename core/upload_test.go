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

	batcher.Set(kvStreamId, key, []byte("hello world")).
		SetKeyToSpecial(kvStreamId, key).
		GrantSpecialWriteRole(kvStreamId, key, defaultAccount)

	_, err := batcher.Exec(context.Background())
	assert.NoError(t, err)

	// query
	time.Sleep(10 * time.Second)
	val, err := kvClientForIterator.GetValue(context.Background(), kvStreamId, key)
	assert.NoError(t, err)
	assert.True(t, val.Size > 0)

	isSpecialKey, err := kvClientForIterator.IsSpecialKey(context.Background(), kvStreamId, key)
	assert.NoError(t, err)
	assert.True(t, isSpecialKey)

	isWriter, err := kvClientForIterator.IsWriterOfKey(context.Background(), defaultAccount, kvStreamId, key)
	assert.NoError(t, err)
	assert.True(t, isWriter)
}
