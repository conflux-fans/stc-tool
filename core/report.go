package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type uploadRecord struct {
	Content string
	Root    common.Hash
}

type BatchUploadReport struct {
	StartTime time.Time
	EndTime   time.Time
	Hash      common.Hash
	Records   []uploadRecord
}

func (b *BatchUploadReport) Save(file string) error {
	dirPath := filepath.Dir(file)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return errors.WithMessage(err, "Failed to create directory")
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.WithMessage(err, "Failed to open file")
	}

	_, err = f.WriteString(b.String())
	return errors.WithMessage(err, "Failed to write report to file")
}

func (b *BatchUploadReport) String() string {
	j, _ := json.Marshal(b)
	return string(j)
}
