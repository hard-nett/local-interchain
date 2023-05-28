package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

type LogOutput struct {
	ChainID     string `json:"chain-id"`
	ChainName   string `json:"chain-name"`
	RPCAddress  string `json:"rpc-address"`
	GRPCAddress string `json:"grpc-address"`
	IBCPath     string `json:"ibc-path"`
}

const filename = "../configs/logs.json"

func WriteRunningChains(bz []byte) {
	_ = ioutil.WriteFile(filename, bz, 0644)
}

func DumpChainsInfoToLogs(t *testing.T, config *MainConfig, chains []ibc.Chain) (*cosmos.CosmosChain, int) {
	// This may be un-needed.
	var longestTTLChain *cosmos.CosmosChain
	ttlWait := 0

	var outputLogs []LogOutput

	// Iterate chain config & get the ibc chain's to save data to logs.
	for idx, chain := range config.Chains {
		chainObj := chains[idx].(*cosmos.CosmosChain)
		t.Logf("\n\n\n\nWaiting for %d blocks on chain %s", chain.BlocksTTL, chainObj.Config().ChainID)

		v := LogOutput{
			// TODO: Rest API Address?
			ChainID:     chainObj.Config().ChainID,
			ChainName:   chainObj.Config().Name,
			RPCAddress:  chainObj.GetHostRPCAddress(),
			GRPCAddress: chainObj.GetHostGRPCAddress(),
			IBCPath:     chain.IBCPath,
		}

		if chain.BlocksTTL > ttlWait {
			ttlWait = chain.BlocksTTL
			longestTTLChain = chainObj
		}

		outputLogs = append(outputLogs, v)
	}

	// dump output logs to file
	bz, _ := json.MarshalIndent(outputLogs, "", "  ")
	WriteRunningChains([]byte(bz))

	return longestTTLChain, ttlWait
}
