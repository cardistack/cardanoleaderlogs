package cmd

/*
Copyright © 2022 cardistack cardistack@protonmail.com
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/leekchan/accounting"
	"github.com/spf13/cobra"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Check Current Epoch",
	Long: `Check Current Epoch:

Check Current Epoch for any Scheduled block for yoor Pool.`,
	Run: func(cmd *cobra.Command, args []string) {

		poolidstatus, _ := cmd.Flags().GetBool("poolid")
		bfidstatus, _ := cmd.Flags().GetBool("bfid")

		if poolidstatus && bfidstatus {

			currentEpoch(args)
		} else {

			fmt.Println("Please Enter your --poolid poolid --bfid  bfid")

		}

	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
	currentCmd.Flags().BoolP("poolid", "p", false, "please Enter poolid")
	currentCmd.Flags().BoolP("bfid", "b", false, "please Enter blockfrost(bfid)")
}

//main function for current parameter
func currentEpoch(arg []string) {
    
	ac := accounting.Accounting{Symbol: "", Precision: 0}
	Ada := " \u20B3"
	Lovelaces := 1000000

	//different colours for print
	const colorReset = "\033[0m"
	const colorRed = "\033[31m"
	const colorGreen = "\033[32m"
	const colorYellow = "\033[33m"
	const colorBlue = "\033[34m"
	const colorPurple = "\033[35m"
	const colorCyan = "\033[36m"
	const colorWhite = "\033[37m"

	fmt.Printf("Lovelace is %v%v\n", Ada, Lovelaces)

	//various URLS to access Blockfrost
	poolurl := "https://cardano-mainnet.blockfrost.io/api/v0/pools/" + arg[0]
	poolmetaurl := "https://cardano-mainnet.blockfrost.io/api/v0/pools/" + arg[0] + "/metadata"
	epochlatestparamurl := "https://cardano-mainnet.blockfrost.io/api/v0/epochs/latest/parameters"
	epochlatesturl := "https://cardano-mainnet.blockfrost.io/api/v0/epochs/latest"
	genesisurl := "https://cardano-mainnet.blockfrost.io/api/v0/genesis"
	firstshelleyurl := "https://cardano-mainnet.blockfrost.io/api/v0/blocks/4555184"

	//Return Bytes for Blockfrost urls
	poolblockresponseBytes := poolBlockfrostQuery(poolurl, arg[1])
	metaresponseBytes := poolBlockfrostQuery(poolmetaurl, arg[1])
	epochparamresponseBytes := poolBlockfrostQuery(epochlatestparamurl, arg[1])
	epochlatestresponseBytes := poolBlockfrostQuery(epochlatesturl, arg[1])
	genesisresponseBytes := poolBlockfrostQuery(genesisurl, arg[1])
	firstshelleyresponseBytes := poolBlockfrostQuery(firstshelleyurl, arg[1])

	//initialize structs
	poolStakeParam := PoolStakeParam{}
	poolMetaData := PoolMetaData{}
	epochParam := EpochParam{}
	epochLatest := EpochLatest{}
	genesisParam := GenesisParam{}
	firstshelleyParam := FirstshelleyParam{}

	var poolVrfSkey string
	var data Vrfkey

	//unmarshal Blockfrost Data
	if err := json.Unmarshal(poolblockresponseBytes, &poolStakeParam); err != nil {
		fmt.Printf("Could not unmarshal poolblockreponseBytes. %v", err)
	}

	if err := json.Unmarshal(metaresponseBytes, &poolMetaData); err != nil {
		fmt.Printf("Could not unmarshal poolmetaData. %v", err)
	}

	if err := json.Unmarshal(epochparamresponseBytes, &epochParam); err != nil {
		fmt.Printf("Could not unmarshal EpochParameters. %v", err)
	}

	if err := json.Unmarshal(epochlatestresponseBytes, &epochLatest); err != nil {
		fmt.Printf("Could not unmarshal EpochLatest. %v", err)
	}

	if err := json.Unmarshal(genesisresponseBytes, &genesisParam); err != nil {
		fmt.Printf("Could not unmarshal genesisdata. %v", err)
	}

	if err := json.Unmarshal(firstshelleyresponseBytes, &firstshelleyParam); err != nil {
		fmt.Printf("Could not unmarshal genesisdata. %v", err)
	}

	//Load vrf key from file
	file, err := ioutil.ReadFile("vrf.skey")
	
    
   
	if err != nil {
		log.Fatalf("check vrf key exist %v", err)

	} else {
		data = Vrfkey{}

	}

	//unmarshal vrfkey
	_ = json.Unmarshal([]byte(file), &data)
	if len(data.CborHex) == 0 {
		fmt.Println("Cant get cbor check your vrf.key")
	} else {

		poolVrfSkey = data.CborHex[4:]
		fmt.Println(poolVrfSkey)
	}

	nactiveStake, _ := strconv.Atoi(epochLatest.Active_stake)
	nStake := nactiveStake / Lovelaces

	pactiveStake, _ := strconv.Atoi(poolStakeParam.Active_stake)
	pStake := pactiveStake / Lovelaces

	//calculate first slot of target epoch
	firstSlotOfEpoch := (firstshelleyParam.Slot) + (epochParam.Epoch-211)*genesisParam.Epoch_length

	color.HiMagenta("Checking SlotLeader Schedules for Stakepool: " + poolMetaData.Ticker + ".......")
	fmt.Println("Pool Id: " + arg[0])
	fmt.Println("Epoch: " + strconv.Itoa(epochParam.Epoch))
	fmt.Println("Nonce: " + epochParam.Nonce)
	fmt.Println("Network Active Stake in Epoch " + strconv.Itoa(epochLatest.Epoch) + " [" + color.HiCyanString((ac.FormatMoney(nStake))) + "]")
	fmt.Println("Pool Active Stake in Epoch " + strconv.Itoa(epochParam.Epoch) + " [" + color.HiCyanString((ac.FormatMoney(pStake))) + "]")
	//fmt.Println("firstSlotOfEpoch: "+strconv.Itoa(firstSlotOfEpoch))

	sigma := poolStakeParam.Active_size
	epochLength := genesisParam.Epoch_length

	// check for  Stake Pool Leader
	var slotcount = 0
	for slot := firstSlotOfEpoch; slot < epochLength+firstSlotOfEpoch; slot++ {

		slotLeader := isSlotLeader(slot, float64(genesisParam.Active_slots_coefficient), sigma, epochParam.Nonce, poolVrfSkey)

		if slotLeader {
			timestamp := time.Unix(int64(slot+1591566291), 0)

			slotcount ++
			fmt.Printf("Leader At Slot: %v  - Local Time %v - Scheduled Epoch Blocks: %v \n", slot-firstSlotOfEpoch, time.Time(timestamp), slotcount)
		}
	}

	if slotcount == 0 {
		fmt.Println(
			"No SlotLeader Schedules Found for Epoch " + fmt.Sprintf("%v", epochParam.Epoch),
		)

	}

}
