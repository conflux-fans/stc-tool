package core

import (
	"context"
	"fmt"
	"os"

	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/conflux-fans/storage-cli/constants/enums"
	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ERR_UNEXIST_CONTENT = errors.New("Unexists content name")
	downloader          Downloader
)

type Downloader struct {
	zgDownloader *transfer.Downloader
}

func InitDefaultDownloader() {
	_zgDownloader, err := transfer.NewDownloader(zgNodeClients, zgLogOpt)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create zg downloader")
	}
	downloader.zgDownloader = _zgDownloader
}

func DefaultDownloader() *Downloader {
	return &downloader
}

func (d *Downloader) DownloadFile(root string, savePath string) {
	if err := d.zgDownloader.Download(context.Background(), root, root, false); err != nil {
		logger.Get().WithField("root", root).WithError(err).Fatal("Failed to download file")
	}
	// rename file
	if err := os.Rename(root, savePath); err != nil {
		logger.Get().WithField("root", root).WithError(err).Fatal("Failed to rename file")
	}
}

func (d *Downloader) DownloadExtend(name string, showMetadata, outputToConsole bool) error {
	logger.Get().WithField("name", name).Info("Start download content")

	meta, err := GetContentMetadata(name)
	if err != nil {
		return errors.WithMessage(err, "Failed to get content metadata")
	}
	logger.Get().WithField("metadata", meta).Info("Get content metadata")

	// 如果是 POINTER，读取第一样为 root，然后根据 root获取文件，写到 f
	// 如果是 TEXT，读取所有行写到 f
	if meta.ExtendDataType == enums.EXTEND_DATA_POINTER {
		if err := d.downloadToFileByPointer(meta); err != nil {
			return errors.WithMessage(err, "Failed to download file")
		}
	} else {
		if err := d.downloadToFileByText(meta); err != nil {
			return errors.WithMessage(err, "Failed to download file")
		}
	}

	// 如果文件长度大于 1k，则提示文本大于 1k，不输出到控制台
	if outputToConsole {
		if err := d.displayOnConsole(meta); err != nil {
			return errors.WithMessage(err, "Failed to display on console")
		}
	}

	logger.Get().Info(fmt.Sprintf("Download data %s to file %s completed ", meta.Name, meta.SaveFile()))

	return nil
}

func (d *Downloader) downloadToFileByPointer(meta *ContentMetadata) error {
	// Read the first line as file root
	root, err := kvClientForIterator.GetValue(context.Background(), kvStreamId, []byte(meta.LineKeys()[0]))
	if err != nil {
		return errors.WithMessage(err, "Failed to get root value")
	}
	logger.Get().WithField("root", fmt.Sprintf("%x", root.Data[:32])).WithField("data", fmt.Sprintf("%x", root.Data)).Info("Got root")

	d.DownloadFile(common.BytesToHash(root.Data[:32]).Hex(), meta.SaveFile())
	logger.Get().Info("Downloaded POINTER type file")
	return nil
}

func (d *Downloader) downloadToFileByText(meta *ContentMetadata) error {
	f, err := os.OpenFile(meta.SaveFile(), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return errors.WithMessage(err, "Failed to open file")
	}
	defer f.Close()
	// TEXT 类型继续执行后续代码
	for _, k := range meta.LineKeys() {
		val, err := kvClientForIterator.GetValue(context.Background(), kvStreamId, []byte(k))
		if err != nil {
			return errors.WithMessage(err, "Failed to get line value")
		}
		logger.Get().WithField("key", k).WithField("val", val).Debug("Get line value")
		_, err = f.Write(val.Data)
		if err != nil {
			return errors.WithMessage(err, "Failed to write file")
		}
	}
	logger.Get().Info(fmt.Sprintf("Download data %s to file %s completed ", meta.Name, meta.SaveFile()))
	return nil
}

func (d *Downloader) displayOnConsole(meta *ContentMetadata) error {
	fileInfo, err := os.Stat(meta.SaveFile())
	if err != nil {
		return errors.WithMessage(err, "Failed to get file information")
	}

	if fileInfo.Size() > 1024 {
		logger.Get().Warn("File size exceeds 1k, not displaying on console")
	} else {
		content, err := os.ReadFile(meta.SaveFile())
		if err != nil {
			return errors.WithMessage(err, "Failed to read file")
		}
		metaMap := meta.ToMap()
		metaMap["content"] = string(content)
		logger.SuccessfWithParams(metaMap, "Download content completed")
	}
	return nil
}
