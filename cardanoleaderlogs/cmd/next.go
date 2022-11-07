/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/leekchan/accounting"
	"github.com/spf13/cobra"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Check Next Epoch",
	Long: `Check Next Epoch:

	Check Current Epoch for any Scheduled block for yoor Pool.`,
	Run: func(cmd *cobra.Command, args []string) {

		poolidstatus, _ := cmd.Flags().GetBool("poolid")
		bfidstatus, _ := cmd.Flags().GetBool("bfid")

        if (poolidstatus && bfidstatus){
			
		nextEpoch(args)
		} else {
          
			fmt.Println("Please Enter your --poolid poolid --bfid  bfid")

		}
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nextCmd.Flags().BoolP("poolid", "p", false, "please Enter poolid")
	nextCmd.Flags().BoolP("bfid", "b", false, "please Enter blockfrost(bfid)")
	
}

func getNextArmada(baseAPI string)  []byte {

	request, err := http.NewRequest(
        http.MethodGet, //method
        baseAPI,        //url
        nil,            //body
    )

    if err != nil {
        log.Printf("Could not request. %v", err)
    }

    request.Header.Add("Accept", "application/json")
    request.Header.Add("User-Agent", "cardanoleaderlogs CLI (https://nonce.armada-alliance.io/next)")

	response, err := http.DefaultClient.Do(request)
    if err != nil {
        log.Printf("Could not make a request. %v", err)
    }

    responseBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Printf("Could not read response body. %v", err)
    }
    
	
	
    return responseBytes

}

func nextEpoch(arg []string){
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

	url := "https://nonce.armada-alliance.io/next"
	poolurl := "https://cardano-mainnet.blockfrost.io/api/v0/pools/"+arg[0]
    poolmetaurl := "https://cardano-mainnet.blockfrost.io/api/v0/pools/"+arg[0]+"/metadata"
	genesisurl := "https://cardano-mainnet.blockfrost.io/api/v0/genesis"
	firstshelleyurl := "https://cardano-mainnet.blockfrost.io/api/v0/blocks/4555184" 

	responseBytes := getNextArmada(url)
	poolblockresponseBytes := poolBlockfrostQuery(poolurl, arg[1])
    metaresponseBytes := poolBlockfrostQuery(poolmetaurl, arg[1])
	genesisresponseBytes := poolBlockfrostQuery(genesisurl, arg[1])
	firstshelleyresponseBytes := poolBlockfrostQuery(firstshelleyurl, arg[1])

	armadaNext := ArmadaNext{}
	poolStakeParam := PoolStakeParam{}
	poolMetaData := PoolMetaData{}
	netStakeParam := NetStakeParam{}
	genesisParam := GenesisParam{}
	firstshelleyParam := FirstshelleyParam{}
    var poolVrfSkey string



	if err := json.Unmarshal(responseBytes, &armadaNext); err != nil {
        fmt.Printf("Could not unmarshal reponseBytes. \n %v", err)
		fmt.Printf(" probably something wrong accessing https://nonce.armada-alliance.io/next exiting .. \n")
		os.Exit(1)
    }

	if err := json.Unmarshal(poolblockresponseBytes, &poolStakeParam); err != nil {
        fmt.Printf("Could not unmarshal poolblockreponseBytes. %v", err)
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

	netstakeurl := "https://cardano-mainnet.blockfrost.io/api/v0/epochs/"+ strconv.Itoa(armadaNext.Epoch - 1)
	netstakeresponseBytes := poolBlockfrostQuery(netstakeurl,arg[1])


    if err := json.Unmarshal(netstakeresponseBytes, &netStakeParam); err != nil {
        fmt.Printf("Could not unmarshal poolmetaData. %v", err)
    }
    
	if strings.Contains(armadaNext.Nonce, "errorMessage") {

		fmt.Println(string(colorRed),"(New Nonce Not yet Available) Exiting ...", string(colorReset))
		//fmt.Println(string(colorBlue),"Nonce not yet available : ")

		os.Exit(1)
	} else {

		
        fmt.Println(string(colorGreen),"(New Nonce Available) Proceeding ...", string(colorReset))

	}
	

    file, _ := ioutil.ReadFile("vrf.skey")
    data := Vrfkey{}

	_ = json.Unmarshal([]byte(file), &data)
	if len(data.CborHex) == 0 {
		fmt.Println("Cant get cbor check your vrf.key")
	} else {

		poolVrfSkey = data.CborHex[4:]
		
	}

	nactiveStake, _ := strconv.Atoi(netStakeParam.Active_stake)	
	nStake := nactiveStake / Lovelaces

    pactiveStake, _ := strconv.Atoi(poolStakeParam.Active_stake)	
	pStake := pactiveStake / Lovelaces

	//calculate first slot of target epoch 
	firstSlotOfEpoch := (firstshelleyParam.Slot) + (armadaNext.Epoch - 211)*genesisParam.Epoch_length
	
	color.HiMagenta("Checking SlotLeader Schedules for Stakepool: "+poolMetaData.Ticker + ".......")
	fmt.Println("Pool Id: "+arg[0])
	fmt.Println("Epoch: "+strconv.Itoa(armadaNext.Epoch))
	fmt.Println("Nonce: "+armadaNext.Nonce)
	fmt.Println("Network Active Stake in Epoch "+strconv.Itoa(armadaNext.Epoch - 1) + " ["+color.HiCyanString((ac.FormatMoney(nStake)))+"]")
	fmt.Println("Pool Active Stake in Epoch "+strconv.Itoa(armadaNext.Epoch - 1) + " [" +color.HiCyanString((ac.FormatMoney(pStake)))+"]")
    
	sigma := poolStakeParam.Active_size
    epochLength := genesisParam.Epoch_length
	
    var slotcount = 0
	for slot := firstSlotOfEpoch; slot < epochLength+firstSlotOfEpoch; slot++ {
	
		

	    slotLeader := isSlotLeader(slot, float64(genesisParam.Active_slots_coefficient), sigma, armadaNext.Nonce, poolVrfSkey)

		if slotLeader {
			timestamp := time.Unix(int64(slot+1591566291), 0)
			
			slotcount += 1
			fmt.Printf("Leader At Slot: %v  - Local Time %v - Scheduled Epoch Blocks: %v \n",slot-firstSlotOfEpoch,time.Time(timestamp),slotcount)
		}
	}
	
	if slotcount == 0 {
		fmt.Println(
			"No SlotLeader Schedules Found for Epoch " + fmt.Sprintf("%v", armadaNext.Epoch),
		)
		
	}
}

