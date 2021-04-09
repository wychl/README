# aes 加密

## 基础知识

- 对称加密算法，加密和解密使用一样的密钥的加解密算法。
- 分组密码（block cipher），是每次只能处理特定长度的一块（block）数据的一类加解密算法。
- 目前常见的对称加密算法DES、3DES、AES都是属于分组密码。

AES高级加密标准，这个标准用来替代原先的DES，是现行的对称加密标准。

分组：128bit

秘钥：128bit、192bit、256bit

## ECB模式

- 全称Electronic CodeBook mode，电子密码本模式。
- 分组方式：将明文分组加密之后的结果直接称为密文分组。
- 优点：
    1. 一个分组的损坏不影响其他分组
    2. 可以并行加解密
- 缺点：
    1. 相同的明文分组会转换为相同的密文分组。
    2. 无需破译密码就能操纵明文（每个分组独立且前后文无关，直接增加或删除一个分组不影响其它分组解密过程的正确性）

## CBC模式

- 全称Cipher Block Chaining mode，密码分组链接模式。
- 分组方式：将明文分组与前一个密文分组进行XOR运算，然后再进行加密。每个分组的加解密都依赖于前一个分组。而第一个分组没有前一个分组，因此需要一个初始化向量（initialization vector）。
- 优点：
    1. 加密结果与前文相关，有利于提高加密结果的随机性。
    2. 可并行解密。
- 缺点
    1. 无法并行加密。
    2. 一个分组损坏，如果密文长度不变，则两个分组受影响。
    3. 一个分组损坏，如果密文长度改变，则后面所有分组受影响。

## CFB模式

- 全称Cipher FeedBack mode，密文反馈模式。
- 分组方式：前一个密文分组会被送回到密码算法的输入端。
    在CBC和EBC模式中，明文分组都是通过密码算法进行加密的。而在CFB模式中，明文分组并没有通过加密算法直接进行加密，明文分组和密文分组之间只有一个XOR。
- CFB模式是通过将“明文分组”与“密码算法的输出”进行XOR运行生成“密文分组”。
- CFB模式中由密码算法生成的比特序列称为密钥流（key stream）。密码算法相当于密钥流的伪随机数生成器，而初始化向量相当于伪随机数生成器的种子。（CFB模式有点类似一次性密码本。）
- 优点：
    1. 支持并行解密。
    2. 不需要填充（padding）。
- 缺点：
    1. 不能抵御重放攻击（replay attack）。
    2. 不支持并行加密。

## OFB模式

- Output FeedBack mode 输出反馈模式
- 密码算法的输出会反馈到密码算法的输入中（具体见下图）。
- OFB模式中，XOR所需的比特序列（密钥流）可以事先通过密码算法生成，和明文分组无关。只需要提前准备好所需的密钥流，然后进行XOR运算就可以了。

## 分组模式小结

推荐使用CBC模式。为什么要填充？ECB和CBC模式要求明文数据必须填充至长度为分组长度的整数倍。那么需要填充多少个字节？

paddingSize = blockSize - textLength % blockSize

填充什么内容？

- ANSI X.923：填充序列的最后一个字节填paddingSize，其它填0。
- ISO 10126：填充序列的最后一个字节填paddingSize， 其它填随机数。
- PKCS7：填充序列的每个字节都填paddingSize。

## 代码实例

```go
package main

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "fmt"
    "encoding/base64"
)

func main() {
    TestCBCAES()
}

var commonIV = []byte{0x1A, 0x3B, 0x36, 0x22, 0xD6, 0xE2, 0x2E, 0xD0, 0x22, 0xFB, 0xB8, 0x75, 0xDD, 0x38, 0x22, 0x11}
var secretKey = []byte{0x12, 0x4D, 0x4A, 0x3E, 0xC2, 0x08, 0x4A, 0x21, 0x41, 0xC1, 0xD5, 0xC5, 0xA8, 0x6A, 0xEE, 0xA1}
var planText = []byte(`ciika test go aes`)

func TestCBCAES() {
    fmt.Printf("Plain Text : %s\n", planText)

    // 加密
    result, err := AESEncrypt(planText, commonIV, secretKey)
    if err != nil {
        panic(err)
    }
    //nkmnjWoKu89yXXe+tTDoZxJIe7q/RGAD2JqXDHRgPoU=
    fmt.Println(base64.StdEncoding.EncodeToString(result))
    fmt.Printf("Encrypted Hex Data : %x\n", result)

    // 解密
    origData, err := AESDecrypt(result, commonIV, secretKey)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decrypted String : %s\n", origData)
}

// 加密函数
func AESEncrypt(origData, iv []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    origData = PKCS5Padding(origData, blockSize)
    blockMode := cipher.NewCBCEncrypter(block, iv)
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}

// 解密函数
func AESDecrypt(crypted, iv []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockMode := cipher.NewCBCDecrypter(block, iv)
    origData := make([]byte, len(crypted))
    // origData := crypted
    blockMode.CryptBlocks(origData, crypted)
    origData = PKCS5UnPadding(origData)
    return origData, nil
}

// padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

// unpadding
func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    // remove the last byte
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
```