// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package api

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ChainSafe/sygma-fee-oracle/config"
	"github.com/ChainSafe/sygma-fee-oracle/util"
	"github.com/pkg/errors"
)

func (h *Handler) rateSignature(result *FetchRateResp, fromDomainID int, resourceTokenAddr, resourceID string) (string, error) {
	fromDomainBaseCurrency := h.conf.GetRegisteredResources(config.ResourceIDBuilder(config.NativeCurrencyAddr, fromDomainID))
	if fromDomainBaseCurrency == nil {
		return "", errors.New("failed to find the registered resource for native currency of from domain")
	}
	baseRate, err := util.Large2SmallUnitConverter(result.BaseRate, uint(fromDomainBaseCurrency.Decimal))
	if err != nil {
		return "", errors.Wrap(err, "failed to convert BaseRate")
	}
	finalBaseEffectiveRate := util.PaddingZero(baseRate.Bytes(), 32)

	tokenRateCurrency := h.conf.GetRegisteredResources(resourceID)
	if tokenRateCurrency == nil {
		return "", errors.New("failed to find the registered resource for given address and domainId")
	}
	tokenRate, err := util.Large2SmallUnitConverter(result.TokenRate, uint(tokenRateCurrency.Decimal))
	if err != nil {
		return "", errors.Wrap(err, "failed to convert TokenRate")
	}
	finalTokenEffectiveRate := util.PaddingZero(tokenRate.Bytes(), 32)

	gasPrice, err := util.Str2BigInt(result.DestinationChainGasPrice)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert DestinationChainGasPrice")
	}
	finalGasPrice := util.PaddingZero(gasPrice.Bytes(), 32)

	finalTimestamp := fmt.Sprintf("%064x", result.ExpirationTimestamp)
	finalFromDomainId := util.PaddingZero([]byte{uint8(result.FromDomainID)}, 32)
	finalToDomainId := util.PaddingZero([]byte{uint8(result.ToDomainID)}, 32)

	finalResourceId, err := hex.DecodeString(resourceID[len(h.conf.GetRegisteredDomains(fromDomainID).AddressPrefix):])
	if err != nil {
		return "", errors.Wrap(err, "failed to decode resourceID")
	}
	
	finalTimestampBytes, err := hex.DecodeString(finalTimestamp)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode timestamp")
	}

	feeDataMessageByte := bytes.Buffer{}
	feeDataMessageByte.Write(finalBaseEffectiveRate)
	feeDataMessageByte.Write(finalTokenEffectiveRate)
	feeDataMessageByte.Write(finalGasPrice)
	feeDataMessageByte.Write(finalTimestampBytes)
	feeDataMessageByte.Write(finalFromDomainId)
	feeDataMessageByte.Write(finalToDomainId)
	feeDataMessageByte.Write(finalResourceId)
	feeDataRaw := feeDataMessageByte.Bytes()

	signature, err := h.identity.Sign(feeDataRaw)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign")
	}

	// modify v
	// openzepplin is verifying if v is 27/28, need to manually add 27 to v
	sigb := bytes.Buffer{}
	sigb.Write(signature[:64])
	sigb.WriteByte(byte(int8(signature[64]) + 27))
	signature = sigb.Bytes()

	return hex.EncodeToString(signature), nil
}

// TODO: this is the placeholder for the algorithm of calculating the data expiration
func (h *Handler) dataExpirationManager(baseTimestamp int64) int64 {
	return baseTimestamp + h.conf.DataValidIntervalConfig()
}
