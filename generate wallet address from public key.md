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

      return prefix + result
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

The SEC format we mentioned before has flaws:
1. it still too long to be memorize or recognize by human
2. has pitfalls on security

In order to remove flaws above, we will use new encoding scheme when we want to generate bitcoin wallet address:

1, if the address is for mainnet, set first byte to 0x00, for testnet set first byte to 0x6f
2, encode the public key in SEC format(compressed or uncompressed), do sha256 and then follow ripemd160 hash, we can combine these two hash
into an operation called hash160
3, combine the first byte from step 1 and bytes from step 2
4, do a hash256 on the result of step 3 and get the first 4 bytes from the result, this is called base58 checksum.
5, combine bytes array from step 3 and step 4 together and encode it by using base58

Let's use code to implement steps above, first we do the base58 checksum and hash160 first in util.go:
```g
import (
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

...

func Base58Checksum(s []byte) []byte {
	//do hash256 on s and append the first 4 bytes of the result to s
	hash256 := Hash256(string(s))
	return append(s, hash256[:4]...)
}

func Hash160(s []byte) []byte {
	//do hash256 and do ripemd160 following
	hash256 := Hash256(string(s))
	hasher := ripemd160.New()
	hasher.Write(hash256)
	hashBytes := hasher.Sum(nil)
	return hashBytes
}
```

Now let's generate the wallet address by using above methods, in point.go add the following code:
```g
func (p *Point) Sec(compressed bool) (string, []byte) {
	secBytes := []byte{}
	if !compressed {
		/*
			uncompressed sec:
			1. first byte 04
			2. x in big endian hex string
			3. y in big endian hex string
			padding x,y with leading 0
		*/
		secBytes = append(secBytes, 0x04)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		secBytes = append(secBytes, p.y.num.Bytes()...)
		return fmt.Sprintf("04%064x%064x", p.x.num, p.y.num), secBytes
	}

	var opMod big.Int
	if opMod.Mod(p.y.num, big.NewInt(int64(2))).Cmp(big.NewInt(int64(0))) == 0 {
		//y is even, set first byte t0 0x02
		secBytes = append(secBytes, 0x02)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		return fmt.Sprintf("02%064x", p.x.num), secBytes
	} else {
		secBytes = append(secBytes, 0x03)
		secBytes = append(secBytes, p.x.num.Bytes()...)
		return fmt.Sprintf("03%064x", p.x.num), secBytes
	}
}


func (p *Point) hash160(compressed bool) []byte {
	_, secBytes := p.Sec(compressed)
	return Hash160(secBytes)
}

func (p *Point) Address(compressed bool, testnet bool) string {
	hash160 := p.hash160(compressed)
	//if mainnet set first byte to 0x00 , 0x6f for testnet
	prefix := []byte{}
	if testnet {
		prefix = append(prefix, 0x6f)
	} else {
		prefix = append(prefix, 0x00)
	}
	//do base58 checksum
	return Base58Checksum(append(prefix, hash160...))
}
```

Let's test the above code:
```g
package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
        privateKey := ecc.NewPrivateKey(big.NewInt(int64(5002)))
	pubKey := privateKey.GetPublicKey()
	fmt.Printf("wallet address for 5002*G is %s\n", pubKey.Address(false, true))

	//2020 ^ 5*G
	var expOp big.Int
	privateKey = ecc.NewPrivateKey(expOp.Exp(big.NewInt(int64(2020)), big.NewInt(int64(5)), nil))
	pubKey = privateKey.GetPublicKey()
	fmt.Printf("wallet address for 2020^5 * G is %s\n", pubKey.Address(true, true))

	//0x12345deadbeef * G
	p := new(big.Int)
	p.SetString("12345deadbeef", 16)
	privateKey = ecc.NewPrivateKey(p)
	pubKey = privateKey.GetPublicKey()
	fmt.Printf("wallet address for 0x12345deadbeef*G is %s\n", pubKey.Address(true, false))
}

```

The running result of the above code is :
```g
wallet address for 5002*G is mmTPbXQFxboEtNRkwfh6K51jvdtHLxGeMA
wallet address for 2020^5 * G is mopVkxp8UhXqRYbCYJsbeE1h1fiF64jcoH
wallet address for 0x12345deadbeef*G is 1F1Pn2y6pDb68E5nYJJeba4TLg2U7B6KF1
```

All those are wallet address generated by given public key for mainnet or testnet. Beside the encoding for public key, there is also encoding for private 
key, of course the private key is not a thing should transmit on the net, if you loss your private key, you loss all your assets in your account.

But sometimes you may need to write down your private key on paper or transmit the private key from one wallet to other, then an encoding scheme for 
private key is also needed. The Encoding scheme for private key called WIF, its steps are following:

1, set first byte to 0x80 for mainnet, 0xef for testnet

2, append the private key bytes array after the first byte, if the length of bytes array smaller than 32, then we
need to padd it with 0

3, if the public key address was compressed , add a suffix byte with value 0x01 after step 2.

4, do hash256 on the result of step 3 and get its first 4 bytes

5, combine result from step 3 and 4, and encode it by Base58

Let's see how to use code for it, in private_key.go, add the following:
```g
func (p *PrivateKey) Wif(compressed bool, testnet bool) string {
	//add first byte with value 0x80 for mainnet, 0xef for testnet
	bytes := []byte{}
	if testnet {
		bytes = append(bytes, 0xef)
	} else {
		bytes = append(bytes, 0x80)
	}
       	//append the secret bytes array
	secretBytes := p.secret.Bytes()
	if len(secretBytes) < 32 {
		//two chars into one byte
		s := fmt.Sprintf("%064x", p.secret.Bytes())
		paddingBytes, err := hex.DecodeString(s)
		if err != nil {
			panic(fmt.Sprintf("padding secret bytes err: %v\n", err))
		}
		secretBytes = paddingBytes
	}

	bytes = append(bytes, secretBytes...)
	//if the public is SEC compressed, add suffix byte with value 0x01
	if compressed {
		bytes = append(bytes, 0x01)
	}

	return Base58Checksum(bytes)
}
```
Let's test the code above:
```g
package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
privateKey := ecc.NewPrivateKey(big.NewInt(int64(5003)))
	fmt.Printf("WIF for private key 5003 is %s\n", privateKey.Wif(true, true))

	//2021 ^ 5*G
	var expOp big.Int
	privateKey = ecc.NewPrivateKey(expOp.Exp(big.NewInt(int64(2021)), big.NewInt(int64(5)), nil))
	fmt.Printf("WIF for private key 2021^5  is %s\n", privateKey.Wif(false, true))

	//0xdeadbeef54321 * G
	p := new(big.Int)
	p.SetString("deadbeef54321", 16)
	privateKey = ecc.NewPrivateKey(p)
	fmt.Printf("WIF for private key deadbeef54321 is %s\n", privateKey.Wif(true, false))
}
```
The running result for the above code is:
```g
WIF for private key 5003 is cMahea7zqjxrtgAbB7LSGbcQUr1uX1ojuat9jZodMN8rFTv2sfUK
WIF for private key 2021^5  is 91avARGdfge8E4tZfYLoxeJ5sGBdNJQH4kvjpWAxgzczjbCwxic
WIF for private key deadbeef54321 is KwDiBf89QgGbjEhKnhXJuH7LrciVrZi3qYjgtNr6kJz3AAYY7Thi
```
Those are encodings for given private keys



