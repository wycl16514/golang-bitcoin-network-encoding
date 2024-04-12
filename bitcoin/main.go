package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {

	//0xdeadbeef54321
	// p := new(big.Int)
	// p.SetString("deadbeef54321", 16)
	// privateKey := ecc.NewPrivateKey(p)
	// pubKey := privateKey.GetPublicKey()
	// fmt.Printf("pub key is %s\n", pubKey)
	// secBinUnCompressed := new(big.Int)
	// secBinUnCompressed.SetString(pubKey.Sec(false), 16)
	// unUnCompressedDecode := ecc.ParseSEC(secBinUnCompressed.Bytes())
	// fmt.Printf("decode sec uncompressed format:%s\n", unUnCompressedDecode)

	// secBinCompressed := new(big.Int)
	// secBinCompressed.SetString(pubKey.Sec(true), 16)
	// compressedDecode := ecc.ParseSEC(secBinCompressed.Bytes())
	// fmt.Printf("decode sec compressed format:%s\n", compressedDecode)

	// r := new(big.Int)
	// r.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	// rField := ecc.S256Field(r)
	// s := new(big.Int)
	// s.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	// sField := ecc.S256Field(s)
	// sig := ecc.NewSignature(rField, sField)
	// derEncode := sig.Der()
	// fmt.Printf("der encoding for signatue is %x\n", derEncode)

	// val := new(big.Int)
	// val.SetString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 16)
	// fmt.Printf("base58 encoding is %s\n", ecc.EncodeBase58(val.Bytes()))

	// val.SetString("eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c", 16)
	// fmt.Printf("base58 encoding is %s\n", ecc.EncodeBase58(val.Bytes()))

	// val.SetString("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", 16)
	// fmt.Printf("base58 encoding is %s\n", ecc.EncodeBase58(val.Bytes()))

	// //5003
	// privateKey := ecc.NewPrivateKey(big.NewInt(int64(5003)))
	// fmt.Printf("wif for private key 5003 is %s\n", privateKey.Wif(true, true))

	// //2021^5
	// var expOp big.Int
	// privateKey = ecc.NewPrivateKey(expOp.Exp(big.NewInt(int64(2021)), big.NewInt(int64(5)), nil))
	// fmt.Printf("wif for private key 2021^5 is %s\n", privateKey.Wif(false, true))

	// //0xdeadbeef54321
	// p := new(big.Int)
	// p.SetString("deadbeef54321", 16)
	// privateKey = ecc.NewPrivateKey(p)
	// fmt.Printf("wif for private key 0xdeadbeef54321 is %s\n", privateKey.Wif(true, false))

	p := new(big.Int)
	p.SetString("12345678", 16)
	bytes := p.Bytes()
	fmt.Printf("bytes for 0x12345678 is %x\n", bytes)

	littleEndianByte := ecc.BigIntToLittleEndian(p, ecc.LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("little endian for 0x12345678 is %x\n", littleEndianByte)

	littleEndianByteToInt64 := ecc.LittleEndianToBigInt(littleEndianByte, ecc.LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("little endian bytes into int is %x\n", littleEndianByteToInt64)
}
