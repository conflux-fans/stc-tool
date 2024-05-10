package core

import (
	"os"
	"strings"

	"github.com/conflux-fans/storage-cli/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

func CreateTemplate(name string, keys []string) error {
	logger.Get().WithField("name", name).WithField("keys", keys).Info("create template")
	tx, err := templates.AddTemplate(&bind.TransactOpts{Signer: signerFn, From: defaultAccount}, name, keys)
	if err != nil {
		return err
	}
	logger.Get().WithField("tx", tx.Hash()).Info("template created")
	return nil
}

// gen csv
func DownloadTemplate(name string) (string, error) {
	header := "field,content"
	keys, err := templates.GetTemplate(nil, name)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to get template from remote")
	}
	logger.Get().WithField("keys", keys).Info("get template keys")

	lines := append([]string{header}, keys...)
	content := strings.Join(lines, "\n")

	filePath := "./" + name + ".template.csv"
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to create template file")
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to write content to file")
	}
	return filePath, nil
}

func DeleteTemplate(name string) error {
	tx, err := templates.DeleteTemplate(&bind.TransactOpts{Signer: signerFn, From: defaultAccount}, name)
	if err != nil {
		return err
	}
	logger.Get().WithField("tx", tx.Hash()).Info("template deleted")
	return nil
}

func ListTemplate() ([]string, error) {
	return templates.ListTemplates(nil)
}
