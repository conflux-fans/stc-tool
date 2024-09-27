package core

import (
	"context"
	"fmt"
	"strconv"

	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/pkg/utils/commonutils"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type ContentMetadata struct {
	Name           string
	LineTotal      int
	ExtendDataType enums.ExtendDataType
	OwnerTokenID   string
}

func GetContentMetadata(name string) (*ContentMetadata, error) {
	m := &ContentMetadata{
		Name: name,
	}
	// query size
	v, err := kvClientForIterator.GetValue(context.Background(), kvStreamId, []byte(m.LineTotalKey()))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get file line size")
	}
	if v.Size == 0 {
		return nil, ERR_UNEXIST_CONTENT
	}

	m.LineTotal, err = strconv.Atoi(string(v.Data))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to convert")
	}

	v, err = kvClientForIterator.GetValue(context.Background(), kvStreamId, []byte(m.ExtendDataTypeKey()))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get file extend data type")
	}
	m.ExtendDataType, err = enums.ParseExtendDataType(string(v.Data))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to parse extend data type")
	}

	v, err = kvClientForIterator.GetValue(context.Background(), kvStreamId, []byte(m.ExtendDataOwnerTokenIDKey()))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get file extend data owner token id")
	}
	m.OwnerTokenID = string(v.Data)
	return m, nil
}

func (m *ContentMetadata) LineKeys() []string {
	lineKeys := lo.Map(make([]int, m.LineTotal), func(v int, index int) string {
		return m.LineIndexKey(index)
	})
	return lineKeys
}

func (m *ContentMetadata) LineTotalKey() string {
	return fmt.Sprintf("%s:line", m.Name)
}

func (m *ContentMetadata) LineIndexKey(index int) string {
	return fmt.Sprintf("%s:%d", m.Name, index)
}

func (m *ContentMetadata) ExtendDataTypeKey() string {
	return fmt.Sprintf("%s:type", m.Name)
}

func (m *ContentMetadata) ExtendDataOwnerTokenIDKey() string {
	return fmt.Sprintf("%s:owner_token_id", m.Name)
}

func (m *ContentMetadata) AllKeys() []string {
	keys := append(m.LineKeys(), m.LineTotalKey(), m.ExtendDataTypeKey(), m.ExtendDataOwnerTokenIDKey())
	return keys
}

func (m *ContentMetadata) ToMap() map[string]string {
	return commonutils.StructToStringMap(m)
}

func (m *ContentMetadata) SaveFile() string {
	return m.Name + ".zg"
}
