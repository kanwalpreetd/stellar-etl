package transform

import (
	"fmt"
	"strconv"

	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
	"github.com/stellar/stellar-etl/internal/utils"
)

// TransformConfigSetting converts an config setting ledger change entry into a form suitable for BigQuery
func TransformConfigSetting(ledgerChange ingest.Change) (ConfigSettingOutput, error) {
	ledgerEntry, changeType, outputDeleted, err := utils.ExtractEntryFromChange(ledgerChange)
	if err != nil {
		return ConfigSettingOutput{}, err
	}

	configSetting, ok := ledgerEntry.Data.GetConfigSetting()
	if !ok {
		return ConfigSettingOutput{}, fmt.Errorf("Could not extract config setting from ledger entry; actual type is %s", ledgerEntry.Data.Type)
	}

	// LedgerEntryChange must contain a config setting state change to be parsed, otherwise skip
	if ledgerEntry.Data.Type != xdr.LedgerEntryTypeConfigSetting {
		return ConfigSettingOutput{}, nil
	}

	configSettingId := configSetting.ConfigSettingId

	contractMaxSizeBytes, _ := configSetting.GetContractMaxSizeBytes()

	contractCompute, _ := configSetting.GetContractCompute()
	ledgerMaxInstructions := contractCompute.LedgerMaxInstructions
	txMaxInstructions := contractCompute.TxMaxInstructions
	feeRatePerInstructionsIncrement := contractCompute.FeeRatePerInstructionsIncrement
	txMemoryLimit := contractCompute.TxMemoryLimit

	contractLedgerCost, _ := configSetting.GetContractLedgerCost()
	ledgerMaxReadLedgerEntries := contractLedgerCost.LedgerMaxReadLedgerEntries
	ledgerMaxReadBytes := contractLedgerCost.LedgerMaxReadBytes
	ledgerMaxWriteLedgerEntries := contractLedgerCost.LedgerMaxWriteLedgerEntries
	ledgerMaxWriteBytes := contractLedgerCost.LedgerMaxWriteBytes
	txMaxReadLedgerEntries := contractLedgerCost.TxMaxReadLedgerEntries
	txMaxReadBytes := contractLedgerCost.TxMaxReadBytes
	txMaxWriteLedgerEntries := contractLedgerCost.TxMaxWriteLedgerEntries
	txMaxWriteBytes := contractLedgerCost.TxMaxWriteBytes
	feeReadLedgerEntry := contractLedgerCost.FeeReadLedgerEntry
	feeWriteLedgerEntry := contractLedgerCost.FeeWriteLedgerEntry
	feeRead1Kb := contractLedgerCost.FeeRead1Kb
	bucketListTargetSizeBytes := contractLedgerCost.BucketListTargetSizeBytes
	writeFee1KbBucketListLow := contractLedgerCost.WriteFee1KbBucketListLow
	writeFee1KbBucketListHigh := contractLedgerCost.WriteFee1KbBucketListHigh
	bucketListWriteFeeGrowthFactor := contractLedgerCost.BucketListWriteFeeGrowthFactor

	contractHistoricalData, ok := configSetting.GetContractHistoricalData()
	feeHistorical1Kb := contractHistoricalData.FeeHistorical1Kb

	contractMetaData, _ := configSetting.GetContractEvents()
	txMaxContractEventsSizeBytes := contractMetaData.TxMaxContractEventsSizeBytes
	feeContractEvents1Kb := contractMetaData.FeeContractEvents1Kb

	contractBandwidth, _ := configSetting.GetContractBandwidth()
	ledgerMaxTxsSizeBytes := contractBandwidth.LedgerMaxTxsSizeBytes
	txMaxSizeBytes := contractBandwidth.TxMaxSizeBytes
	feeTxSize1Kb := contractBandwidth.FeeTxSize1Kb

	paramsCpuInsns, _ := configSetting.GetContractCostParamsCpuInsns()
	contractCostParamsCpuInsns := serializeParams(paramsCpuInsns)

	paramsMemBytes, _ := configSetting.GetContractCostParamsMemBytes()
	contractCostParamsMemBytes := serializeParams(paramsMemBytes)

	contractDataKeySizeBytes, ok := configSetting.GetContractDataKeySizeBytes()

	contractDataEntrySizeBytes, ok := configSetting.GetContractDataEntrySizeBytes()

	stateExpirationSettings, _ := configSetting.GetStateExpirationSettings()
	maxEntryExpiration := stateExpirationSettings.MaxEntryExpiration
	minTempEntryExpiration := stateExpirationSettings.MinTempEntryExpiration
	minPersistentEntryExpiration := stateExpirationSettings.MinPersistentEntryExpiration
	persistentRentRateDenominator := stateExpirationSettings.PersistentRentRateDenominator
	tempRentRateDenominator := stateExpirationSettings.TempRentRateDenominator
	maxEntriesToExpire := stateExpirationSettings.MaxEntriesToExpire
	bucketListSizeWindowSampleSize := stateExpirationSettings.BucketListSizeWindowSampleSize
	evictionScanSize := stateExpirationSettings.EvictionScanSize
	startingEvictionScanLevel := stateExpirationSettings.StartingEvictionScanLevel

	contractExecutionLanes, _ := configSetting.GetContractExecutionLanes()
	ledgerMaxTxCount := contractExecutionLanes.LedgerMaxTxCount

	bucketList, _ := configSetting.GetBucketListSizeWindow()
	bucketListSizeWindow := make([]uint64, 0, len(bucketList))
	for _, sizeWindow := range bucketList {
		bucketListSizeWindow = append(bucketListSizeWindow, uint64(sizeWindow))
	}

	transformedPool := ConfigSettingOutput{
		ConfigSettingId:                 int32(configSettingId),
		ContractMaxSizeBytes:            uint32(contractMaxSizeBytes),
		LedgerMaxInstructions:           int64(ledgerMaxInstructions),
		TxMaxInstructions:               int64(txMaxInstructions),
		FeeRatePerInstructionsIncrement: int64(feeRatePerInstructionsIncrement),
		TxMemoryLimit:                   uint32(txMemoryLimit),
		LedgerMaxReadLedgerEntries:      uint32(ledgerMaxReadLedgerEntries),
		LedgerMaxReadBytes:              uint32(ledgerMaxReadBytes),
		LedgerMaxWriteLedgerEntries:     uint32(ledgerMaxWriteLedgerEntries),
		LedgerMaxWriteBytes:             uint32(ledgerMaxWriteBytes),
		TxMaxReadLedgerEntries:          uint32(txMaxReadLedgerEntries),
		TxMaxReadBytes:                  uint32(txMaxReadBytes),
		TxMaxWriteLedgerEntries:         uint32(txMaxWriteLedgerEntries),
		TxMaxWriteBytes:                 uint32(txMaxWriteBytes),
		FeeReadLedgerEntry:              int64(feeReadLedgerEntry),
		FeeWriteLedgerEntry:             int64(feeWriteLedgerEntry),
		FeeRead1Kb:                      int64(feeRead1Kb),
		BucketListTargetSizeBytes:       int64(bucketListTargetSizeBytes),
		WriteFee1KbBucketListLow:        int64(writeFee1KbBucketListLow),
		WriteFee1KbBucketListHigh:       int64(writeFee1KbBucketListHigh),
		BucketListWriteFeeGrowthFactor:  uint32(bucketListWriteFeeGrowthFactor),
		FeeHistorical1Kb:                int64(feeHistorical1Kb),
		TxMaxContractEventsSizeBytes:    uint32(txMaxContractEventsSizeBytes),
		FeeContractEvents1Kb:            int64(feeContractEvents1Kb),
		LedgerMaxTxsSizeBytes:           uint32(ledgerMaxTxsSizeBytes),
		TxMaxSizeBytes:                  uint32(txMaxSizeBytes),
		FeeTxSize1Kb:                    int64(feeTxSize1Kb),
		ContractCostParamsCpuInsns:      contractCostParamsCpuInsns,
		ContractCostParamsMemBytes:      contractCostParamsMemBytes,
		ContractDataKeySizeBytes:        uint32(contractDataKeySizeBytes),
		ContractDataEntrySizeBytes:      uint32(contractDataEntrySizeBytes),
		MaxEntryExpiration:              uint32(maxEntryExpiration),
		MinTempEntryExpiration:          uint32(minTempEntryExpiration),
		MinPersistentEntryExpiration:    uint32(minPersistentEntryExpiration),
		PersistentRentRateDenominator:   int64(persistentRentRateDenominator),
		TempRentRateDenominator:         int64(tempRentRateDenominator),
		MaxEntriesToExpire:              uint32(maxEntriesToExpire),
		BucketListSizeWindowSampleSize:  uint32(bucketListSizeWindowSampleSize),
		EvictionScanSize:                uint64(evictionScanSize),
		StartingEvictionScanLevel:       uint32(startingEvictionScanLevel),
		LedgerMaxTxCount:                uint32(ledgerMaxTxCount),
		BucketListSizeWindow:            bucketListSizeWindow,
		LastModifiedLedger:              uint32(ledgerEntry.LastModifiedLedgerSeq),
		LedgerEntryChange:               uint32(changeType),
		Deleted:                         outputDeleted,
	}
	return transformedPool, nil
}

func serializeParams(costParams xdr.ContractCostParams) []map[string]string {
	params := make([]map[string]string, 0, len(costParams))
	for _, contractCostParam := range costParams {
		serializedParam := map[string]string{}
		serializedParam["ExtV"] = strconv.Itoa(int(contractCostParam.Ext.V))
		serializedParam["ConstTerm"] = strconv.Itoa(int(contractCostParam.ConstTerm))
		serializedParam["LinearTerm"] = strconv.Itoa(int(contractCostParam.LinearTerm))
		params = append(params, serializedParam)
	}

	return params
}
