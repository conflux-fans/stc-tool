package core

import (
	ccore "github.com/0glabs/0g-storage-client/core"
	"github.com/conflux-fans/storage-cli/constants/enums"
)

type ExtendDataConverter struct {
}

func (p *ExtendDataConverter) ByContent(data []byte) (enums.ExtendDataType, ccore.IterableData, error) {
	_data, err := ccore.NewDataInMemory(data)
	if err != nil {
		return enums.ExtendDataType(-1), nil, err
	}
	return enums.EXTEND_DATA_TEXT, _data, nil
}

func (p *ExtendDataConverter) ByFile(filePath string) (enums.ExtendDataType, ccore.IterableData, error) {
	f, err := ccore.Open(filePath)
	if err != nil {
		return enums.ExtendDataType(-1), nil, err
	}
	return enums.EXTEND_DATA_POINTER, f, nil
}
