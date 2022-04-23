/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net/http"
    "io/ioutil"
	"log"
    "encoding/json"
	"github.com/spf13/cobra"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("next called")
		nextEpoch()
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
	
	
}

type ArmadaNext struct {
	Epoch  int `json:"epoch"`
	Nonce  string `json:"nonce"`
}

func nextEpoch(){

	Ada := " \u20B3"
	Lovelaces := 1000000
    fmt.Println(Ada)
	fmt.Println(Lovelaces)

	url := "https://nonce.armada-alliance.io/next"
	responseBytes := getNextArmada(url)
    armadaNext := ArmadaNext{}

	if err := json.Unmarshal(responseBytes, &armadaNext); err != nil {
        fmt.Printf("Could not unmarshal reponseBytes. %v", err)
    }

    fmt.Println((armadaNext))


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