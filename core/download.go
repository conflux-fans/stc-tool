package core

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
)

func DownloadFile(root string, savePath string) {
	downloader := transfer.NewDownloader(nodeClients...)
	if err := downloader.Download(root, root, false); err != nil {
		logrus.WithField("root", root).WithError(err).Fatal("Failed to download file")
	}
}

func DownloadDataByKv(name string) error {
	// query size
	v, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(fmt.Sprintf("%s:line", name)))
	if err != nil {
		return errors.WithMessage(err, "Failed to get file line size")
	}

	if v.Size == 0 {
		return errors.New("Unexists name")
	}

	lineCountInStr := string(v.Data)
	lineCount, _ := strconv.Atoi(lineCountInStr)

	f, err := os.OpenFile(name+".zg", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := 0; i < lineCount; i++ {
		val, err := kvClientForIterator.GetValue(STREAM_FILE, []byte(fmt.Sprintf("%s:%d", name, i)))
		if err != nil {
			return err
		}
		_, err = f.Write(val.Data)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Download data %s to file %s.zg completed ", name, name)
	return nil
}
