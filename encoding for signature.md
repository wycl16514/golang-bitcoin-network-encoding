For signature we can't use SEC for encoding, because there is not close relation between r and s, you know one of them but you can't derive the other. There is a scheme to encode signature
called DER (Distinguished Encoding Rules), following are steps for encoding signature:

1, set the first byte to 0x30

2, the second byte is the length of signature (usually 0x44 or 0x45)

3, the third byte set to 0x02, this is an indcator that following bytes are for r.

4, transfer r into bytes array, if the first byte of r is >= 0x80, than we append a byte 0x00 in the beginning of the bytes, compute the length of the bytes array, append the value of length
after the marker byte 0x02 in step 3, and append the bytes array following the value of length

5, add a marker byte 0x02 at the end of bytes array from step 4.

6. append s as how we append r in step 4.

We need to encode the length of r and s, because r and s at most has 32 bytes, but some times their length may shorter than this. Let's see an example for the encoding of r and s:

30 45 02 21 00 ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f 02 20 7a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed

the first byte is 0x30 as mentioned above, the second byte is 0x45 that is the total length of r and s. The third byte is marker 0x02 indicating from next byte is the biginning of bytes array for r.
According to 3, the byte following from marker 0x02 is the length of bytes array for r, the value is 0x21. Byte following 0x21 is 0x00, by step 4, it indicates the first byte of r is >= 0x80, we can
see the byte follow 0x00 is 0xed which is indeed more than 0x80, the length of r is 0x21 - 1 = 0x20, which means the following 32 bytes are bytes array for r, we extract it out as following:

r: 81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f

following the last byte of r is 0x20 , it is indicator for the beginning of s, the byte following 0x02 is 0x20 which indicates the length of s is 0x20, the following byte is not 0x00, which means the
first btye of s is not more than 0x80, and the byte following 0x20 is 0x7a , this is the beginning byte for s and is smaller than 0x80, the byte array for s is :

s: 7a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed

Let's put the encoding scheme into code, in signature.go we add the following code:
```g
func (s *Signature) Der() []byte {
	rBin := s.r.num.Bytes()
	//if the first byte >= 0x80, append 0x00 at beginning
	if rBin[0] >= 0x80 {
		rBin = append([]byte{0x00}, rBin...)
	}
	//insert indicator 0x02 and the length of rBin
	rBin = append([]byte{0x02, byte(len(rBin))}, rBin...)
	//do the same for s
	sBin := s.s.num.Bytes()
	//if the first byte of s >= 0x80, append 0x00 at beginning of s
	if sBin[0] >= 0x80 {
		sBin = append([]byte{0x00}, sBin...)
	}
	//insert indicator 0x02 and the length of sBin
	sBin = append([]byte{0x02, byte(len(sBin))}, sBin...)
	//combine rBin , sBin and insert 0x30 and the total length of sBin, rBin at the beginning
	derBin := append([]byte{0x30, byte(len(rBin) + len(sBin))}, rBin...)
	derBin = append(derBin, sBin...)

	return derBin
}

```
Let's test the code above, we will construct a signature and call the Der method to get the encoding:
```g
package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
  r := new(big.Int)
	r.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	rField := ecc.S256Field(r)
	s := new(big.Int)
	s.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	sField := ecc.S256Field(s)
	sig := ecc.NewSignature(rField, sField)
	derEncode := sig.Der()
	fmt.Printf("der encoding for signature is: %x\n", derEncode)
}
```
The result of running the above code is :
```g
der encoding for signature is: 3045022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec
```

