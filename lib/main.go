package main

import "C"

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

//export DoVRF
func DoVRF(cinput *C.char) *C.char {
	input := C.GoString(cinput)
	return C.CString(doVRF(bufio.NewReader(strings.NewReader(input))))
}

//export VerifyVRF
func VerifyVRF(cinput *C.char) *C.char {
	input := C.GoString(cinput)
	return C.CString(verifyVRF(bufio.NewReader(strings.NewReader(input))))
}

func readString(reader *bufio.Reader) string {
	ii, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	if ii[len(ii)-1] == '\n' {
		ii = ii[:len(ii)-1]
	}
	return ii
}

func getRawPK(reader *bufio.Reader) ecdsa.PublicKey {
	x := readString(reader)
	y := readString(reader)
	n := big.NewInt(0)
	xpk, _ := n.SetString(x, 10)
	m := big.NewInt(0)
	ypk, _ := m.SetString(y, 10)

	return ecdsa.PublicKey{
		Curve: curve,
		X:     xpk,
		Y:     ypk,
	}
}

func getPK(reader *bufio.Reader) PublicKey {
	res := getRawPK(reader)
	return PublicKey{&res}
}

func getRawSK(reader *bufio.Reader) ecdsa.PrivateKey {
	epk := getRawPK(reader)

	d := readString(reader)
	o := big.NewInt(0)
	dsk, _ := o.SetString(d, 10)

	return ecdsa.PrivateKey{
		PublicKey: epk,
		D:         dsk,
	}
}

func getSK(reader *bufio.Reader) PrivateKey {
	res := getRawSK(reader)
	return PrivateKey{&res}
}

func doVRF(reader *bufio.Reader) string {
	sk := getSK(reader)
	smsg := readString(reader)
	msg, err := base64.StdEncoding.DecodeString(smsg)
	if err != nil {
		return "error: invalid message"
	}

	array_idx, proof := sk.Evaluate(msg)
	idx := array_idx[:]

	return (base64.StdEncoding.EncodeToString(idx) + "\n" +
		base64.StdEncoding.EncodeToString(proof))
}

func verifyVRF(reader *bufio.Reader) string {
	pk := getPK(reader)
	smsg := readString(reader)
	msg, err := base64.StdEncoding.DecodeString(smsg)
	if err != nil {
		return "error: invalid message"
	}

	sidx := readString(reader)
	given_idx, err := base64.StdEncoding.DecodeString(sidx)
	if err != nil {
		return "error: invalid index"
	}

	sproof := readString(reader)
	proof, err := base64.StdEncoding.DecodeString(sproof)
	if err != nil {
		return "error: invalid proof"
	}

	array_idx, err := pk.ProofToHash(msg, proof)
	if err != nil {
		return fmt.Sprintf("error: ProofToHash threw error: %v", err)
	}
	idx := array_idx[:]
	if bytes.Compare(idx, given_idx) != 0 {
		return "error: index did not match."
	}

	return ""
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	if len(os.Args) > 1 && os.Args[1] == "do" {
		doVRF(reader)
	} else {
		verifyVRF(reader)
	}
}
