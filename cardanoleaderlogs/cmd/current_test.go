package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestCurrent(t *testing.T) {

poolid := os.Getenv("POOLID")
bfid := os.Getenv("BFID")
callargs := []string{poolid, bfid}
currentEpoch(callargs)

fmt.Printf("testing %v %v \n",poolid, bfid)
t.Log()
}