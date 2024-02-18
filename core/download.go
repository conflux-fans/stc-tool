package core

import (
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
)

func Download(root string) {
	downloader := transfer.NewDownloader(nodeClients...)
	if err := downloader.Download(root, root, false); err != nil {
		logrus.WithError(err).Fatal("Failed to download file")
	}
}
