package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	slotWhitelistDeployerMapping = map[string]uint64{
		"whitelisted":  1,
		"whitelistAll": 2,
	}
	slotValidatorMapping = map[string]uint64{
		"validators": 1,
	}
)

// IsWhitelistedDeployer reads the contract storage to check if an address is allow to deploy
func IsWhitelistedDeployer(statedb *StateDB, address common.Address) bool {
	contract := common.HexToAddress(common.WhitelistDeployerSC)
	whitelistAllSlot := slotWhitelistDeployerMapping["whitelistAll"]
	whitelistAll := statedb.GetState(contract, GetLocSimpleVariable(whitelistAllSlot))
	if whitelistAll.Big().Cmp(big.NewInt(1)) == 0 {
		return true
	}

	whitelistedSlot := slotWhitelistDeployerMapping["whitelisted"]
	valueLoc := GetLocMappingAtKey(address.Hash(), whitelistedSlot)
	whitelisted := statedb.GetState(contract, valueLoc)

	return whitelisted.Big().Cmp(big.NewInt(1)) == 0
}

func GetValidators(statedb *StateDB) []common.Address {
	slot := slotValidatorMapping["validators"]
	slotHash := common.BigToHash(new(big.Int).SetUint64(slot))
	arrLength := statedb.GetState(common.HexToAddress(common.ValidatorSC), slotHash)
	keys := []common.Hash{}
	for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
		key := GetLocDynamicArrAtElement(slotHash, i, 1)
		keys = append(keys, key)
	}
	rets := []common.Address{}
	for _, key := range keys {
		ret := statedb.GetState(common.HexToAddress(common.ValidatorSC), key)
		rets = append(rets, common.HexToAddress(ret.Hex()))
	}
	return rets
}
