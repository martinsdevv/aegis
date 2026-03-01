package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	sum := sha256.Sum256([]byte("cliente2"))
	fmt.Println(hex.EncodeToString(sum[:]))
}
