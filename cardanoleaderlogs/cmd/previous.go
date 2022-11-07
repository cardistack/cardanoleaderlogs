/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/leekchan/accounting"

	//"github.com/leekchan/accounting"
	"github.com/spf13/cobra"
)

// previousCmd represents the previous command
var previousCmd = &cobra.Command{
	Use:   "previous",
	Short: "Check Previous Epoch",
	Long: `Check Previous Epoch:

	Check Previous Epoch for any Scheduled block for yoor Pool.`,
	Run: func(cmd *cobra.Command, args []string) {
		
        epochnostatus, _ := cmd.Flags().GetBool("epochno")
		poolidstatus, _ := cmd.Flags().GetBool("poolid")
		bfidstatus, _ := cmd.Flags().GetBool("bfid")

        if (epochnostatus && poolidstatus && bfidstatus){
		previousEpoch(args)
		} else {
          
			fmt.Println("Please Enter your --epochno epoch --poolid poolid --bfid  bfid")

		}
		
	},
}

func init() {
	rootCmd.AddCommand(previousCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// previousCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// previousCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    previousCmd.Flags().BoolP("epochno", "e", false, "please Enter epochno")
	previousCmd.Flags().BoolP("poolid", "p", false, "please Enter poolid")
	previousCmd.Flags().BoolP("bfid", "b", false, "please Enter blockfrost(bfid)")
}


func previousEpoch(arg []string){

	ac := accounting.Accounting{Symbol: "", Precision: 0}
	
	Ada := " \u20B3"
	Lovelaces := 1000000

	const colorReset = "\033[0m"
    const colorRed = "\033[31m"
    const colorGreen = "\033[32m"
    const colorYellow = "\033[33m"
    const colorBlue = "\033[34m"
    const colorPurple = "\033[35m"
    const colorCyan = "\033[36m"
    const colorWhite = "\033[37m"
   
	fmt.Printf("Lovelace is %v%v\n",Ada,Lovelaces)

	epochlatestparamurl := "https://cardano-mainnet.blockfrost.io/api/v0/epochs/"+arg[0]+"/parameters"
	netstakeparamurl := "https://cardano-mainnet.blockfrost.io/api/v0/epochs/"+arg[0]
	poolhiststakeurl := "https://cardano-mainnet.blockfrost.io/api/v0/pools/"+arg[1]+"/history?page=2"
	poolmetaurl := "https://cardano-mainnet.blockfrost.io/api/v0/pools/"+arg[1]+"/metadata"
	genesisurl := "https://cardano-mainnet.blockfrost.io/api/v0/genesis"
	firstshelleyurl := "https://cardano-mainnet.blockfrost.io/api/v0/blocks/4555184"

    epochprevrespBytes := poolBlockfrostQuery(epochlatestparamurl, arg[2])
	netstakeprevrespBytes := poolBlockfrostQuery(netstakeparamurl, arg[2])
	poolhistrespBytes     := poolBlockfrostQuery(poolhiststakeurl, arg[2])
	metaresponseBytes := poolBlockfrostQuery(poolmetaurl, arg[2])
	genesisresponseBytes := poolBlockfrostQuery(genesisurl, arg[2])
	firstshelleyresponseBytes := poolBlockfrostQuery(firstshelleyurl, arg[2])

	epochParam := EpochParam{}
	netstakeParam := NetStakeParam{}
	poolhistoryParam := PoolHistoryParam{}
	poolMetaData := PoolMetaData{}
	genesisParam := GenesisParam{}
    firstshelleyParam := FirstshelleyParam{}
	var sigma float64
	var pstake int
	var nstake float64
	var poolVrfSkey string

	
	//nStake := netstakeParam.Active_stake

	if err := json.Unmarshal(epochprevrespBytes, &epochParam); err != nil {
        fmt.Printf("Could not unmarshal epochParam. %v", err)
    }

	if err := json.Unmarshal(netstakeprevrespBytes, &netstakeParam); err != nil {
        fmt.Printf("Could not unmarshal netstakeParam. %v", err)
    }

	if err := json.Unmarshal(poolhistrespBytes, &poolhistoryParam); err != nil {
        fmt.Printf("Could not unmarshal poolhistoryParam %v", err)
    }
	if err := json.Unmarshal(metaresponseBytes, &poolMetaData); err != nil {
        fmt.Printf("Could not unmarshal poolmetaData. %v", err)
    }

	if err := json.Unmarshal(genesisresponseBytes, &genesisParam); err != nil {
        fmt.Printf("Could not unmarshal genesisdata. %v", err)
    }

	if err := json.Unmarshal(firstshelleyresponseBytes, &firstshelleyParam); err != nil {
        fmt.Printf("Could not unmarshal genesisdata. %v", err)
    }
    //fmt.Println(poolhistoryParam[1].Epoch)
	//fmt.Println(netstakeParam.Epoch)
	//fmt.Println(PoolHistoryParam.poolshistory)
    epochno,_ := strconv.Atoi(arg[0])
	for i := range poolhistoryParam {
       
        
		//fmt.Println(poolhistoryParam[i].Epoch)
		if poolhistoryParam[i].Epoch == epochno {

            
			sigma = poolhistoryParam[i].Active_size
			pstakeconver,_ := strconv.Atoi(poolhistoryParam[i].Active_stake)
			pstake = pstakeconver / Lovelaces
			nstakeconver , _ := strconv.ParseFloat(netstakeParam.Active_stake, 32)
			nstake = nstakeconver / float64(Lovelaces)
			 //fmt.Println(sigma)
			 //fmt.Println(pstake)
		} 
		

	}
    file, _ := ioutil.ReadFile("vrf.skey")
    data := Vrfkey{}
	
	_ = json.Unmarshal([]byte(file), &data)
	if len(data.CborHex) == 0 {
		fmt.Println("Cant get cbor check your vrf.key")
	} else {

		poolVrfSkey = data.CborHex[4:]
		//fmt.Println(poolVrfSkey)
	}

	color.HiMagenta("Checking SlotLeader Schedules for Stakepool: "+poolMetaData.Ticker)
	fmt.Println("Pool Id: "+arg[1])
	fmt.Println("Epoch: "+strconv.Itoa(epochno))
	fmt.Println("Nonce: "+epochParam.Nonce)
	fmt.Println("Network Active Stake in Epoch "+strconv.Itoa(netstakeParam.Epoch) + " ["+color.HiCyanString((ac.FormatMoney(nstake)))+"]")
    fmt.Printf("Pool Active Stake in Epoch %v [ %v ] \n",strconv.Itoa(epochno),pstake)

	epochLength := genesisParam.Epoch_length
    firstSlotOfEpoch := (firstshelleyParam.Slot) + (epochParam.Epoch - 211)*genesisParam.Epoch_length


    //fmt.Println(epochLength)
    fmt.Println("Checking Leaderlogs...")
	var slotcount = 0
	
	for slot := firstSlotOfEpoch; slot < epochLength+firstSlotOfEpoch; slot++ {
	
	
 
		slotLeader := isSlotLeader(slot, float64(genesisParam.Active_slots_coefficient), sigma, epochParam.Nonce, poolVrfSkey)

		if slotLeader {
			timestamp := time.Unix(int64(slot+1591566291), 0)
			
			slotcount += 1
			

          fmt.Printf("Leader At Slot: %v  - Local Time %v - Scheduled Epoch Blocks: %v \n",slot-firstSlotOfEpoch,time.Time(timestamp),slotcount)

		}
	}
	
    


	
	if slotcount == 0 {
		fmt.Println(
			"No SlotLeader Schedules Found for Epoch " + fmt.Sprintf("%v", epochParam.Epoch),
		)
		
	}



}