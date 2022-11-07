/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

/*
#cgo CFLAGS: -Wall -I${SRCDIR}/../libsodium
#cgo LDFLAGS: -L${SRCDIR}/../libsodium -l:libsodium.a
#include "../libsodium/include/sodium/core.h"
#include "../libsodium/include/sodium/crypto_vrf.h"
#include "../libsodium/include/sodium/crypto_vrf_ietfdraft03.h"
#include "../libsodium/include/sodium.h"
#include "string.h"
#include "stdio.h"
#include "stddef.h"
*/
import "C"
import (
	"bytes"
	"encoding/binary"

	//"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	"golang.org/x/crypto/blake2b"
)



//make see using go bindings for libsodium
func mkSeed(slot int, eta0 string) []byte{
    //initialize nonce
	noncebyte := []byte{0, 0, 0, 0, 0, 0, 0, 1}
	hashnonce := blake2b.Sum256(noncebyte)

	eta0byte := unhexlify(eta0)
	
    slottobytes := IntToBytes(int64(slot))
	slotplusnoncebyte := append(slottobytes, eta0byte...)
    hashslot := blake2b.Sum256(slotplusnoncebyte)
	seed := []byte{}
	for i, x := range hashnonce{
	seed = append(seed, x ^ hashslot[i])
   }
   return seed
}

func unhexlify(str string) []byte {
    res := make([]byte, 0)
    for i := 0; i < len(str); i+=2 {
        x, _ := strconv.ParseInt(str[i:i+2], 16, 32)
        res = append(res, byte(x))
    }
    return res
}



func IntToBytes(num int64) []byte {
    buff := new(bytes.Buffer)
    bigOrLittleEndian := binary.BigEndian
    err := binary.Write(buff, bigOrLittleEndian, num)
	if err != nil {
        log.Panic(err)
    }

    return buff.Bytes()
}


func IntFromBytes(x *big.Int, buf []byte) *big.Int {
    if len(buf) == 0 {
        return x
    }

    if (0x80 & buf[0]) == 0 { // positive number
        return x.SetBytes(buf)
		
    }

    for i := range buf {
        buf[i] = ^buf[i]
    }

    return x.SetBytes(buf).Add(x, big.NewInt(1)).Neg(x)
}


func vrfEvalCertified(seeding []byte, tpraosCanBeLeaderSignKeyVRF []byte ) []byte{
  
   
	if C.sodium_init() == -1 {
		panic("sodium_init() failed")
	}


	
	proofbytes := C.crypto_vrf_ietfdraft03_proofbytes()

	
	proofarray := make([]byte,proofbytes)

	//seed_length := uint64(len(seeding))
	
	
	C.crypto_vrf_prove((*C.uchar)((&proofarray[0])), (*C.uchar)(&tpraosCanBeLeaderSignKeyVRF[0]),(*C.uchar)((&seeding[0])),(C.ulonglong)(len(seeding)))
	
    
	outbytes := C.crypto_vrf_outputbytes()
   
	outarray := make([]byte,outbytes)

   
	C.crypto_vrf_proof_to_hash((*C.uchar)((&outarray[0])),(*C.uchar)((&proofarray[0])))
	


	return outarray

	


  
}

func isSlotLeader(slot int, activeSlotCoeff float64, sigma float64, eta0 string, poolVrfSkey string) bool {
  
	seed := mkSeed(slot, eta0)

	tpraosCanBeLeaderSignKeyVRFb := unhexlify(poolVrfSkey)

	cert := vrfEvalCertified(seed,tpraosCanBeLeaderSignKeyVRFb)

	
	
	z := new(big.Int)
    certNat := z.SetBytes(cert)

	x := new(big.Int).SetInt64(2)
	y := new(big.Int).SetInt64(512)
	certNatMax  := x.Exp(x, y, nil)
	bignumber := big.Int{}
	denominator := bignumber.Sub(certNatMax, certNat)
	denominatorfloat := big.NewFloat(0).SetInt(denominator)
    certNatMaxfloat  := big.NewFloat(0).SetInt(certNatMax )
	q := new(big.Float).Quo(certNatMaxfloat, denominatorfloat)
	c := math.Log(1.0 - activeSlotCoeff)

	
	sigmaOfF := math.Exp(-sigma * c)

	sigmafloat := big.NewFloat(0).SetFloat64(sigmaOfF)
   
	
    compare := q.Cmp(sigmafloat)

    result := false
	if compare == -1 {
		result = true

	} else if

	compare == 1 {
		result = false
	} else {
     result = true
     
	}
	return result
}