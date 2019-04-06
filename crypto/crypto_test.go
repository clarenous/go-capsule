package crypto

import (
	"fmt"
	"golang.org/x/crypto/sha3"
	"testing"
)

func TestRipemd160(t *testing.T) {
	fmt.Println(Ripemd160(nil))
	fmt.Println(sha3.Sum256(nil))
}
