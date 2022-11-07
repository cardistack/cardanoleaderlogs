/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/cardistack/cardanoleaderlogs/cmd"
)

func main() {
	//handle := C.dlopen(C.CString("/usr/local/lib/libsodium.so"), C.RTLD_LAZY)
	//libsodium := C.dlsym(handle, C.CString("str_length"))

	if runtime.GOOS == "linux" {

		fmt.Println("you run linux")
	}
	var path = "vrf.skey"
	var file, err = os.OpenFile(path, os.O_RDONLY, 0644)

	if isError(err) {
		return
	}
	defer file.Close()

	cmd.Execute()

}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Please ensure you have vrf.skey file in the root directory")
	}
	return (err != nil)

}
