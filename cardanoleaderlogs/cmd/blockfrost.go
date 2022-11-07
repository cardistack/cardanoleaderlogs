/*
Copyright © 2022 cardistack cardistack@protonmail.com

*/
package cmd

import (	
	"net/http"
    "io"
	"log" 
)

type EpochParam struct {
	Epoch    int `json:"epoch"`
	Min_fee_a int `json:"min_fee_a"`
	Min_fee_b int  `json:"min_fee_b"`
	Max_block_size int `json:"max_block_size"`
	Max_tx_size int `json:"max_tx_size"`
	Nonce string `json:"nonce"`
}

type EpochLatest struct {
	Epoch    int `json:"epoch"`
	Block_count int `json:"block_count"`
	Fees string `json:"fees"`
	Active_stake string `json:"active_stake"`
}

type PoolMetaData struct {

	Pool_id string `json:"pool_id"`
	Hex string `json:"hex"`
	Url   string   `json:"url"`
	Hash string `json:"Hash"`
	Ticker string `json:"ticker"`
	Name string `json:"name"`
	Description string `json:"description"`
	Homepage string `json:"homepage"`
}

type ArmadaCurrent struct {
	Epoch  int `json:"epoch"`
	Nonce  string `json:"nonce"`
}

type PoolStakeParam struct {
	Pool_id    string `json:"pool_id"`
	Blocks_minted int `json:"blocks_minted"`
	Blocks_epoch int  `json:"blocks_epoch"`
	Live_stake string `json:"live_stake"`
	Live_delegators int `json:"live_delegators"`
	Active_stake string `json:"active_stake"`
	Active_size float64 `json:"active_size"`
}

type GenesisParam struct {
	Active_slots_coefficient float32 `json:"active_slots_coefficient"`
	Epoch_length int `json:"epoch_length"`
	Slot_length int  `json:"slot_length"`
}

type FirstshelleyParam struct {
	Epoch    int `json:"epoch"`
	Slot int `json:"slot"`
}

type NetStakeParam struct {
	Epoch   int `json:"epoch"`
	Fees string `json:"fees"`
	Active_stake string `json:"active_stake"`
}

type PoolHistoryParam []PoolHistory

type PoolHistory struct {
	Epoch    int `json:"epoch"`
	Blocks int `json:"blocks"`
	Active_stake string `json:"active_stake"`
	Active_size float64 `json:"active_size"`
	Delegators_count int  `json:"delegators_count"`
	Rewards string `json:"rewards"`
	Fees string `json:"fees"`
}

type ArmadaNext struct {
	Epoch  int `json:"epoch"`
	Nonce  string `json:"nonce"`
}

func poolBlockfrostQuery(baseAPI string, blockfrostID string)  []byte {

	request, err := http.NewRequest(
        http.MethodGet, //method
        baseAPI,        //url
        nil,            //body
    )

    if err != nil {
        log.Printf("Could not request. %v", err)
    }

    request.Header.Add("content-type", "application/json")
	request.Header.Add("project_id", blockfrostID)
    request.Header.Add("User-Agent", "cardanoleaderlogs CLI "+baseAPI)

	response, err := http.DefaultClient.Do(request)
    if err != nil {
        log.Printf("Could not make a request. %v", err)
    }

    poolblockfrostresponseBytes, err := io.ReadAll(response.Body)
    if err != nil {
        log.Printf("Could not read response body. %v", err)
    }

    return poolblockfrostresponseBytes
}