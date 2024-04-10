bitcoin as an application of blockchain, its build up on distributed system with "nodes" split around the world. In order to make all nodes coorperate together to function as a whole 
system, it need to keep all nodes to sync in the same state, that means nodes will communicate with each other frequently and exchange lots of data messages with each other. This will
require message or data to sending on the network to be encoded with some format to keep the content of data in safty and the encoding scheme should make sure the encoded packet as small
as possible to ensure the efficiency of communication between lots of nodes.

![image](https://github.com/wycl16514/golang-bitcoin-network-encoding/assets/7506958/5dd78856-d4a6-437a-9cab-ce3742f4a6b6)


When a user of bitcoin create a wallet, he/she need to publish the address of his or her wallet and he can receive or sending funds to others. Wallet address is actually public key 
we created at last section and encode in some kind of format called SEC(standards for efficient cryptography). SEC format has two forms they are uncompressed and compressed,let's check
the uncompressed format first. For a public key P=(x,y), the coordinate of x and y are 32 bytes integer, we use the following steps to encode the key into uncompressed SEC format:

1, the beginning byte set to 0x04,
2, turn the x value into big-endian and append it after byte 0x04
3, turn the y value into big-endian and append it after the end of x

Let's see an example of SEC uncompressed format:

047211a824f55b505228e4c3d5194c1fcfaa15a456abdf37f9b9d97a4040afc073dee6c89064984f03385237d92167c13e236446b417ab79a0fcae412ae3316b77

let's split the chunk of data into three parts:

1. beginning byte 0x04
2. x value in big-endian: 7211a824f55b505228e4c3d5194c1fcfaa15a456abdf37f9b9d97a4040afc073
3. y value in big-endian: dee6c89064984f03385237d92167c13e236446b417ab79a0fcae412ae3316b77

let's check whether the point with the given x,y is on the bitcoin curve:
```go
import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
	x := new(big.Int)
	x.SetString("7211a824f55b505228e4c3d5194c1fcfaa15a456abdf37f9b9d97a4040afc073", 16)
	y := new(big.Int)
	y.SetString("dee6c89064984f03385237d92167c13e236446b417ab79a0fcae412ae3316b77", 16)
	//check the point on curve
	ecc.NewEllipticPoint(ecc.S256Field(x), ecc.S256Field(y),
		ecc.S256Field(big.NewInt(int64(0))), ecc.S256Field(big.NewInt(int64(7))))
	fmt.Println("point is on the curve")
}
```
If the point is not on the curve, S256Point will throw a panic otherwise the last line of code will be executed and print out the string "point is on the curve", following is the result of running 
the code:
<img width="561" alt="截屏2024-04-10 12 28 53" src="https://github.com/wycl16514/golang-bitcoin-network-encoding/assets/7506958/0d44edae-9296-42d9-aa45-22fb9cc4a81a">

Now let's put the uncompressed SEC format into code:
```go
func (p *Point) Sec() string {
	/*
		uncompressed sec:
		1. first byte 04
		2. x in big endian hex string
		3. y in big endian hex string
	*/
	return fmt.Sprintf("04%x%x", p.x.num, p.y.num)
}
```
Let's create some point and check the code above, in the Point struct we add the following method:
```g
func (p *Point) Sec() string {
	/*
		uncompressed sec:
		1. first byte 04
		2. x in big endian hex string
		3. y in big endian hex string
    padding x,y with leading 0
	*/
	return fmt.Sprintf("04%064x%064x", p.x.num, p.y.num)
}
```
Let's test above code:
```g
func main() {

	privateKey := ecc.NewPrivateKey(big.NewInt(int64(5000)))
	pubKey := privateKey.GetPublicKey()
	fmt.Printf("sec format for 5000*G is %s\n", pubKey.Sec())

	// var op big.Int
	//2018^5 * G
	var expOp big.Int
	privateKey = ecc.NewPrivateKey(expOp.Exp(big.NewInt(int64(2018)), big.NewInt(int64(5)), nil))
	pubKey = privateKey.GetPublicKey()
	fmt.Printf("sec format for 2018^5 * G is %s\n", pubKey.Sec())

	//0xdeadbeef12345 * G
	p := new(big.Int)
	p.SetString("deadbeef12345", 16)
	privateKey = ecc.NewPrivateKey(p)
	pubKey = privateKey.GetPublicKey()
	fmt.Printf("sec format for 0xdeadbeef12345*G is %s\n", pubKey.Sec())
}
```
Running the code we can get the following result:

![截屏2024-04-10 13 06 35](https://github.com/wycl16514/golang-bitcoin-network-encoding/assets/7506958/dbab1c05-d510-45c8-b169-14a9c71c3f8a)

We have uncompressed format which means we will have compressed format, and it is much troublesome than the uncompressed one. For the curve equation y^2 = x^3 + ax +b, if point (x,y) is on the curve,
which means (x,-y) also on the curve because y^2 in the equation. And remember x, y are finit field element, if the order of the field is p, then we have -y is equivalence to p-y this lead to the 
conlusion that if (x,y) on the curve then (x, p-y) is also on the curve.

Because p is prime then p is also odd. If y is even, then p-y is odd, if y is odd, then p - y is even, we will use these info to compress the public key by only using the x coordinate and an 
indicator for y is even or odd, by doing this we can avoid appending y in the data, that is to say we "compress" the y coordinate, following is te step used to generate the SEC compressed format:

1, if y is even, set the starting byte to 0x02, otherwise set the starting byte to 0x03

2, append x coordinate in 32 bytes as big endian hex string

Let's check an example:

0349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a

we can see the first byte is 0x03, which means the value of y is odd, and the x value is "49fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a", having the value of x and getting y is 
very tedious:

1. Find w^2 = v, when v is given, if the order p can satisfy p % 4 == 3, then finding w will be easy.
2. if p % 4 == 3 => (p+1) % 4 == 0 => (p+1) / 4 is an integer,
3. by Fermat's little theory, W^(p-1) % p == 1, => w^2 == W^2 * 1 => W^2 * W^(p-1) = w^(p+1)
4. because p is prime , then p is odd, then (p+1)/2 is an integer,
5. w * w = w ^ 2 = w ^(p+1) => w = w ^(1) = w^(2/2) = w^((p+1)/2) => w == w^((p+1)/2)
6. (p+1)/4 is integer , (p+1)/2 == 2*((p+1)/4) => W^((p+1)/2) == w^(2*(p+1)/4) => (w^2)^((p+1)/4) == (v)^((p+1)/4)
7.  w == w^((p+1)/2) => w == v^((p+1)/4)

Luckly the order p  for bitcoin field element can satisfy p % 3 == 1, which means if we know the value of x, we plus x into w = x^3+7, then we can y by y = w^((p+1)/4), let's put those steps into 
code and compute the y value given by the data above:
```g
func (f *FieldElement) Sqrt() *FieldElement {
	//make sure (p+1) % 4 == 0
	var opAdd big.Int
	orderAddOne := opAdd.Add(f.order, big.NewInt(int64(1)))
	var opMod big.Int
	modRes := opMod.Mod(orderAddOne, big.NewInt(int64(4)))
	if modRes.Cmp(big.NewInt(int64(0))) != 0 {
		panic("order plus one mod 4 is not 0")
	}

	var opDiv big.Int
	return f.Power(opDiv.Div(orderAddOne, big.NewInt(int64(4))))

}
```
Let's compute y^2 by putting x into the equation for bitcoin curve and using the Sqrt method to get the value of y and check (x,y) is on the bitcoin curve:
```g
import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
x := new(big.Int)
	x.SetString("49fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a", 16)
	//y^2 = x^3 + 7
	y2 := (ecc.S256Field(x).Power(big.NewInt(int64(3)))).Add(ecc.S256Field(big.NewInt(int64(7))))
	y := y2.Sqrt()
	fmt.Printf("y value for gien x is %s\n", y)
	//check x,y on the curve
	ecc.NewEllipticPoint(ecc.S256Field(x), y,
		ecc.S256Field(big.NewInt(int64(0))), ecc.S256Field(big.NewInt(int64(7))))
	fmt.Println("point from the SEC compressed format is on the curve")
}
```
The result of running the above code is :

![截屏2024-04-10 14 34 33](https://github.com/wycl16514/golang-bitcoin-network-encoding/assets/7506958/a8c929a5-a4d3-46b5-be0d-c14815230be4)

Let's add code for the SEC compressed format:
```g
func (p *Point) Sec(compressed bool) string {
	if !compressed {
		/*
			uncompressed sec:
			1. first byte 04
			2. x in big endian hex string
			3. y in big endian hex string
			padding x,y with leading 0
		*/
		return fmt.Sprintf("04%064x%064x", p.x.num, p.y.num)
	}

	var opMod big.Int
	if opMod.Mod(p.y.num, big.NewInt(int64(2))).Cmp(big.NewInt(int64(0))) == 0 {
		//y is even, set first byte t0 0x02
		return fmt.Sprintf("02%064x", p.x.num)
	} else {
		return fmt.Sprintf("03%064x", p.x.num)
	}
}

```
Let's check the SEC compressed format for points we have given below:
```g
import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
privateKey := ecc.NewPrivateKey(big.NewInt(int64(5001)))
	pubKey := privateKey.GetPublicKey()
	fmt.Printf("sec compressed format for 5000*G is %s\n", pubKey.Sec(true))

	// var op big.Int
	//2018^5 * G
	var expOp big.Int
	privateKey = ecc.NewPrivateKey(expOp.Exp(big.NewInt(int64(2019)), big.NewInt(int64(5)), nil))
	pubKey = privateKey.GetPublicKey()
	fmt.Printf("sec compressed format for 2019^5 * G is %s\n", pubKey.Sec(true))

	//0xdeadbeef54321 * G
	p := new(big.Int)
	p.SetString("deadbeef54321", 16)
	privateKey = ecc.NewPrivateKey(p)
	pubKey = privateKey.GetPublicKey()
	fmt.Printf("sec compressed format for 0xdeadbeef54321*G is %s\n", pubKey.Sec(true))
}
```
The running result for the above code is:
![截屏2024-04-10 14 50 56](https://github.com/wycl16514/golang-bitcoin-network-encoding/assets/7506958/ff5c0038-da85-470f-b495-4693fb40c3f6)


Given SEC format string, let's see how we can decode it into the point on bitcoin curve, in util.go we add the following:
```g
func ParseSEC(secBin []byte) *Point {
	//check first byte to determine its sec uncompress or compressed format
	if secBin[0] == 4 {
		//uncompress
		x := new(big.Int)
		x.SetBytes(secBin[1:33])
		y := new(big.Int)
		y.SetBytes(secBin[33:65])
		return S256Point(x, y)
	}

	isEven := (secBin[0] == 2)
	x := new(big.Int)
	x.SetBytes(secBin[1:])
	y2 := S256Field(x).Power(big.NewInt(int64(3))).Add(S256Field(big.NewInt(int64(7))))
	y := y2.Sqrt()
	var modOp big.Int
	var yEven *FieldElement
	var yOdd *FieldElement
	if modOp.Mod(y.num, big.NewInt(int64(2))).Cmp(big.NewInt(int64(0))) == 0 {
		yEven = y
		yOdd = y.Negate()
	} else {
		yOdd = y
		yEven = y.Negate()
	}

	if isEven {
		return S256Point(x, yEven.num)
	} else {
		return S256Point(x, yOdd.num)
	}
}

```
We will use the function above to decode a compressed and uncompressed sec format:
```g
package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
//0xdeadbeef54321 * G
	p := new(big.Int)
	p.SetString("deadbeef54321", 16)
	privateKey := ecc.NewPrivateKey(p)
	pubKey := privateKey.GetPublicKey()
	fmt.Printf("pub key is %s\n", pubKey)
	secBinUnCompressed := new(big.Int)
	secBinUnCompressed.SetString(pubKey.Sec(false), 16)
	unCompressDecode := ecc.ParseSEC(secBinUnCompressed.Bytes())
	fmt.Printf("decode sec uncompressed format: %s\n", unCompressDecode)

	secBinCompressed := new(big.Int)
	secBinCompressed.SetString(pubKey.Sec(true), 16)
	compressedDecode := ecc.ParseSEC(secBinCompressed.Bytes())
	fmt.Printf("decode sec compressed format: %s\n", compressedDecode)
}
```
The running result for the above code is :
```g
pub key is (x: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 68183256789471564274367801875006375475112961631504146173182369992753849792144}, y: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 22766467020259989524460707859208087846960031342427706382509198566352122699446}, a: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 0}, b: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 7})

decode sec uncompressed format: (x: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 68183256789471564274367801875006375475112961631504146173182369992753849792144}, y: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 22766467020259989524460707859208087846960031342427706382509198566352122699446}, a: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 0}, b: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 7})

decode sec compressed format: (x: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 68183256789471564274367801875006375475112961631504146173182369992753849792144}, y: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 22766467020259989524460707859208087846960031342427706382509198566352122699446}, a: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 0}, b: FieldElement{order: 115792089237316195423570985008687907853269984665640564039457584007908834671663, num: 7})
```

we can check that the decode point for compressed and uncompressed SEC data is the same as the public key we generate before.
