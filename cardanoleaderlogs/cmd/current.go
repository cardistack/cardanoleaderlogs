/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

        if (poolidstatus && bfidstatus){
		currentEpoch(args)
		} else {
          
			fmt.Println("Please Enter your --poolid poolid --bfid  bfid")

		}




	},
}

func init() {
	rootCmd.AddCommand(currentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// currentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// currentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	currentCmd.Flags().BoolP("poolid", "p", false, "please Enter poolid")
	currentCmd.Flags().BoolP("bfid", "b", false, "please Enter blockfrost(bfid)")
}

func currentEpoch(arg []string){
    ac := accounting.Accounting{Symbol: "", Precision: 0}
	Ada := " \u20B3"
	Lovelaces := 1000000
    fmt.Println(Ada)
	fmt.Printf("Lovelace is %v\n",Lovelaces)

	url := "https://nonce.armada-alliance.io/current"
	poolurl := "https://cardano-mainnet.blockfrost.io/api/v0/pools/"+arg[0]
	poolmetaurl := "https://cardano-mainnet.blockfrost.io/api/v0/pools/"+arg[0]+"/metadata"
	epochlatestparamurl := "https://cardano-mainnet.blockfrost.io/api/v0/epochs/latest/parameters"
	epochlatesturl := "https://cardano-mainnet.blockfrost.io/api/v0/epochs/latest"
	genesisurl := "https://cardano-mainnet.blockfrost.io/api/v0/genesis"
	firstshelleyurl := "https://cardano-mainnet.blockfrost.io/api/v0/blocks/4555184" 
	
	responseBytes := getCurrentArmada(url)
	poolblockresponseBytes := poolBlockfrostQuery(poolurl, arg[1])
    metaresponseBytes := poolBlockfrostQuery(poolmetaurl, arg[1])
	epochparamresponseBytes := poolBlockfrostQuery(epochlatestparamurl, arg[1])
	epochlatestresponseBytes := poolBlockfrostQuery(epochlatesturl, arg[1])
	genesisresponseBytes := poolBlockfrostQuery(genesisurl, arg[1])
	firstshelleyresponseBytes := poolBlockfrostQuery(firstshelleyurl, arg[1])

	armadaCurrent := ArmadaCurrent{}
	poolStakeParam := PoolStakeParam{}
    poolMetaData := PoolMetaData{}
	epochParam := EpochParam{}
	epochLatest :=  EpochLatest{}
	genesisParam := GenesisParam{}
	firstshelleyParam := FirstshelleyParam{}
	var poolVrfSkey string


	if err := json.Unmarshal(responseBytes, &armadaCurrent); err != nil {
        fmt.Printf("Could not unmarshal reponseBytes. %v", err)
    }

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
    
	file, _ := ioutil.ReadFile("vrf.skey")
    data := Vrfkey{}
	
	_ = json.Unmarshal([]byte(file), &data)
	if len(data.CborHex) == 0 {
		fmt.Println("Cant get cbor check your vrf.key")
	} else {

		poolVrfSkey = data.CborHex
		fmt.Println(poolVrfSkey)
	}

	nactiveStake, _ := strconv.Atoi(epochLatest.Active_stake)	
	nStake := nactiveStake / Lovelaces

    pactiveStake, _ := strconv.Atoi(poolStakeParam.Active_stake)	
	pStake := pactiveStake / Lovelaces

	//calculate first slot of target epoch 
	firstSlotOfEpoch := (firstshelleyParam.Slot) + (epochParam.Epoch - 211)*genesisParam.Epoch_length
	

	color.HiMagenta("Checking SlotLeader Schedules for Stakepool: "+poolMetaData.Ticker)
	fmt.Println("Pool Id: "+arg[0])
	fmt.Println("Epoch: "+strconv.Itoa(epochParam.Epoch))
	fmt.Println("Nonce: "+epochParam.Nonce)
	fmt.Println("Network Active Stake in Epoch "+strconv.Itoa(epochLatest.Epoch) + " ["+color.HiCyanString((ac.FormatMoney(nStake)))+"]")
	fmt.Println("Pool Active Stake in Epoch "+strconv.Itoa(epochParam.Epoch) + " [" +color.HiCyanString((ac.FormatMoney(pStake)))+"]")
	fmt.Println("firstSlotOfEpoch: "+strconv.Itoa(firstSlotOfEpoch))
	

	
	

    //vrfEvalCertified(getseed ,tpraosCan)
    sigma := poolStakeParam.Active_size
    epochLength := genesisParam.Epoch_length
	
    var slotcount = 0
	for slot := firstSlotOfEpoch; slot < epochLength+firstSlotOfEpoch; slot++ {
		//for slot := firstSlotOfEpoch; slot < genesisParam.Epoch_length+firstSlotOfEpoch; slot++ {
		//fmt.Println(slot)
		

		slotLeader := isSlotLeader(slot, float64(genesisParam.Active_slots_coefficient), sigma, epochParam.Nonce, poolVrfSkey)

		if slotLeader {
			timestamp := time.Unix(int64(slot+1591566291), 0)
			
			slotcount += 1
			fmt.Println(
				    "Leader At Slot: " + fmt.Sprintf(
					"%v",
					slot-firstSlotOfEpoch,
				) + " - Local Time " + fmt.Sprintf(
					"%v",
					time.Time(timestamp).Format(
						"%Y-%m-%d %H:%M:%S",
					)+" - Scheduled Epoch Blocks: "+fmt.Sprintf(
						"%v",
						slotcount,
					),
				),
			)
		}
	}
	
	if slotcount == 0 {
		fmt.Println(
			"No SlotLeader Schedules Found for Epoch " + fmt.Sprintf("%v", epochParam.Epoch),
		)
		
	}
	
}

func getCurrentArmada(baseAPI string)  []byte {

	request, err := http.NewRequest(
        http.MethodGet, //method
        baseAPI,        //url
        nil,            //body
    )

    if err != nil {
        log.Printf("Could not request. %v", err)
    }

    request.Header.Add("Accept", "application/json")
    request.Header.Add("User-Agent", "cardanoleaderlogs CLI "+baseAPI)

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
