package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func TransferOwner(streamName string, newAdmin common.Address) error {
	streamId := streamIdByName(streamName)
	isAdmin, err := kvClientForPut.IsAdmin(defaultAccount, streamId)
	if err != nil {
		return errors.WithMessage(err, "Failed to check if admin")
	}

	if !isAdmin {
		return errors.New("You are not admin")
	}

	batcher := kvClientForPut.Batcher()
	batcher.GrantAdminRole(streamId, newAdmin)
	batcher.RenounceAdminRole(streamId)
	return batcher.Exec()
}
