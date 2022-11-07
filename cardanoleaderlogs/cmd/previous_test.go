package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestPrevious(t *testing.T) {

poolid := os.Getenv("POOLID")
bfid := os.Getenv("BFID")
epochno := "360"

callargs := []string{epochno," ", poolid, " ", bfid}
previousEpoch(callargs)

fmt.Printf("testing %v %v \n",poolid, bfid)
t.Log()

}