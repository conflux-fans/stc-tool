package core

import (
	"testing"

	"github.com/conflux-fans/storage-cli/config"
)

func TestGetOwnerHistory(t *testing.T) {
	config.SetConfigFile("/Users/dayong/myspace/mywork/storage-tool/config.yaml")
	config.Init()
	Init()

	history, err := DefaultOwnerOperator().GetOwnerHistory("content1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(history)
}
