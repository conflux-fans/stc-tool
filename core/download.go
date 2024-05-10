package core

import (
	"fmt"
	"os"
	"strconv"

	"github.com/conflux-fans/storage-cli/logger"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
)

var (
	ERR_UNEXIST_CONTENT = errors.New("Unexists content name")
)

func DownloadFile(root string, savePath string) {
	downloader := transfer.NewDownloader(nodeClients...)
	if err := downloader.Download(root, root, false); err != nil {
		logger.Get().WithField("root", root).WithError(err).Fatal("Failed to download file")
	}
	// rename file
	if err := os.Rename(root, savePath); err != nil {
		logger.Get().WithField("root", root).WithError(err).Fatal("Failed to rename file")
	}
}

func DownloadDataByKv(name string) error {
	logger.Get().WithField("name", name).Info("Start download content")

	meta, err := GetContentMetadata(name)
	if err != nil {
		return err
	}

	logger.Get().WithField("metadata", meta).Info("Get content metadata")

	f, err := os.OpenFile(name+".zg", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, k := range meta.LineKeys {
		val, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(k))
		if err != nil {
			return err
		}
		logger.Get().WithField("key", k).WithField("val", val).Debug("Get line value")
		_, err = f.Write(val.Data)
		if err != nil {
			return err
		}
	}
	logger.Get().Info(fmt.Sprintf("Download data %s to file %s.zg completed ", name, name))
	return nil
}

type ContentMetadata struct {
	LineSizeKey string
	LineKeys    []string
	LineSize    int
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

	lineKeys := lo.Map(make([]int, lineCount), func(v int, index int) string {
		return keyLineIndex(name, index)
	})

	return &ContentMetadata{
		LineSizeKey: lineSizeKey,
		LineKeys:    lineKeys,
		LineSize:    lineCount,
	}, nil
}
