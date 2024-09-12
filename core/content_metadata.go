package core

import (
	"fmt"
	"strconv"

	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type ContentMetadata struct {
	LineCountKey string
	LineCount    int
	ExtendData   enums.ExtendDataType
	OwnerTokenID uint64
}

func GetContentMetadata(name string) (*ContentMetadata, error) {
	// query size
	lineSizeKey := keyLineCount(name)
	v, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(lineSizeKey))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get file line size")
	}

	if v.Size == 0 {
		return nil, ERR_UNEXIST_CONTENT
	}

	lineCountInStr := string(v.Data)
	lineCount, err := strconv.Atoi(lineCountInStr)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to convert")
	}

	return &ContentMetadata{
		LineCountKey: lineSizeKey,
		LineCount:    lineCount,
	}, nil
}

func (m *ContentMetadata) LineKeys() []string {
	lineKeys := lo.Map(make([]int, m.LineCount), func(v int, index int) string {
		return keyLineIndex(m.LineCountKey, index)
	})
	return lineKeys
}

func keyLineCount(name string) string {
	return fmt.Sprintf("%s:line", name)
}

func keyLineIndex(name string, index int) string {
	return fmt.Sprintf("%s:%d", name, index)
}
