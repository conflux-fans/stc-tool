package enums

import "github.com/conflux-fans/storage-cli/pkg/utils/enumutils"

type ExtendDataType int

const (
	EXTEND_DATA_POINTER ExtendDataType = iota + 1
	EXTEND_DATA_TEXT
)

var extendDataTypeEb enumutils.EnumBase[ExtendDataType]

func init() {
	extendDataTypeEb = enumutils.NewEnumBase("EXTEND_DATA", map[ExtendDataType]string{
		EXTEND_DATA_POINTER: "POINTER",
		EXTEND_DATA_TEXT:    "TEXT",
	})
}

func (t ExtendDataType) MarshalText() ([]byte, error) {
	return extendDataTypeEb.MarshalText(t)
}

func (t *ExtendDataType) UnmarshalText(data []byte) error {
	val, err := extendDataTypeEb.UnmarshalText(data)
	if err != nil {
		return err
	}
	*t = val
	return nil
}

func (t ExtendDataType) String() string {
	return extendDataTypeEb.String(t)
}

func ParseExtendDataType(s string) (ExtendDataType, error) {
	return extendDataTypeEb.Parse(s)
}
