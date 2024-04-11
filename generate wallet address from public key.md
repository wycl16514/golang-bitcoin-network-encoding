As bitcoin user, we always need to send or receive bitcoins from other, this will require to let others know your wallet address. Because wallet address need to be read by human, all the encoding schem
we have before are produce result in binary, therefore we need another scheme to create wallet address in human friendly way.

Wallet address is actually generated from public key and it need to satisfy the following requirement:

1, easy to read and write, user may want to memorize it or write it down on paper

2, Not too long for sending over the internet

3, It should be secure, and harder to make mistake, you don't want you fund transfer to people unknow to you!

The base58 encoding scheme can help us to achive three goals. Compare with the commonly use of base64, it removes characters like l , I, 0, O, -, _ because they are easy to confuse with each other.
Because the encoding schem uses all numbers, and uppercase and lowercase letters and remove 0 O, l , I, which means it will use 58 characters in the encoding process, you will find its algorithm on
internet easily and we will give the code as below, in util.go:
```g
func EncodeBase58(s []byte) string {
	BASE58_ALPHABET := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	count := 0
	for idx := range s {
		if s[idx] == 0 {
			count += 1
		} else {
			break
		}
	}

	prefix := ""
	for i := 0; i < count; i++ {
		prefix += "1"
	}
	result := ""
	num := new(big.Int)
	num.SetBytes(s)
	for num.Cmp(big.NewInt(int64(0))) > 0 {
		var divOp big.Int
		var modOp big.Int
		mod := modOp.Mod(num, big.NewInt(int64(58)))
		num = divOp.Div(num, big.NewInt(int64(58)))
		result = string(BASE58_ALPHABET[mod.Int64()]) + result
	}
```
Let's test the code above:
```go
package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
val := new(big.Int)
	val.SetString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 16)
	fmt.Printf("base58 encoding is %s\n", ecc.EncodeBase58(val.Bytes()))

	val.SetString("eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c", 16)
	fmt.Printf("base58 encoding is %s\n", ecc.EncodeBase58(val.Bytes()))

	val.SetString("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", 16)
	fmt.Printf("base58 encoding is %s\n", ecc.EncodeBase58(val.Bytes()))
}
```
The result of running above code is :
```g
base58 encoding is 9MA8fRQrT4u8Zj8ZRd6MAiiyaxb2Y1CMpvVkHQu5hVM6
base58 encoding is 4fE3H2E6XMp4SsxtwinF7w9a34ooUrwWe4WsW1458Pd
base58 encoding is EQJsjkd6JaGwxrjEhfeqPenqHwrBmPQZjJGNSCHBkcF7
```
