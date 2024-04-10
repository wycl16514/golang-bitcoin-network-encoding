For signature we can't use SEC for encoding, because there is not close relation between r and s, you know one of them but you can't derive the other. There is a scheme to encode signature
called DER (Distinguished Encoding Rules), following are steps for encoding signature:

1, set the first byte to 0x30
2, the second byte is the length of signature (usually 0x44 or 0x45)
3, the third byte set to 0x02, this is an indcator that following bytes are for r.
4, transfer r into bytes array, if the first byte of r is >= 0x80, than we append a byte 0x00 in the beginning of the bytes, compute the length of the bytes array, append the value of length
after the marker byte 0x02 in step 3, and append the bytes array following the value of length
5, add a marker byte 0x02 at the end of bytes array from step 4.
6. append s as how we append r in step 4.
