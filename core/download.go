package core

import (
	"context"
	"fmt"
	"os"

	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/pkg/errors"
)

var (
	ERR_UNEXIST_CONTENT = errors.New("Unexists content name")
)

func DownloadFile(root string, savePath string) {
	downloader, err := transfer.NewDownloader(zgNodeClients)
	if err != nil {
		logger.Get().WithField("root", root).WithError(err).Fatal("Failed to create downloader")
	}

	if err := downloader.Download(context.Background(), root, root, false); err != nil {
		logger.Get().WithField("root", root).WithError(err).Fatal("Failed to download file")
	}
	// rename file
	if err := os.Rename(root, savePath); err != nil {
		logger.Get().WithField("root", root).WithError(err).Fatal("Failed to rename file")
	}
}

func DownloadExtend(name string, showMetadata, outputToConsole bool) error {
	logger.Get().WithField("name", name).Info("Start download content")

	meta, err := GetContentMetadata(name)
	if err != nil {
		return errors.WithMessage(err, "Failed to get content metadata")
	}
	logger.Get().WithField("metadata", meta).Info("Get content metadata")

	f, err := os.OpenFile(name+".zg", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return errors.WithMessage(err, "Failed to open file")
	}
	defer f.Close()

	for _, k := range meta.LineKeys() {
		val, err := kvClientForIterator.GetValue(context.Background(), STREAM_FILE, []byte(k))
		if err != nil {
			return errors.WithMessage(err, "Failed to get line value")
		}
		logger.Get().WithField("key", k).WithField("val", val).Debug("Get line value")
		_, err = f.Write(val.Data)
		if err != nil {
			return errors.WithMessage(err, "Failed to write file")
		}
	}
	logger.Get().Info(fmt.Sprintf("Download data %s to file %s.zg completed ", name, name))

	if outputToConsole {
		content, err := os.ReadFile(name + ".zg")
		if err != nil {
			return errors.WithMessage(err, "Failed to read file")
		}
		metaMap := meta.ToMap()
		metaMap["Content"] = string(content)
		logger.SuccessfWithParams(metaMap, "Download content completed")
	}

	return nil
}
