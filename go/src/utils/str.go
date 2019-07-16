package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"
)

type Str struct {
}



//md5
func (s *Str) Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 字符串转十进制
func (s *Str) StrToBigInt(str string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(str, 10)
	return n
}

// 十六进制转十进制
func (s *Str) HexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	//n, _ = n.SetString(hex[2:], 16)
	n, _ = n.SetString(hex, 16)
	return n
}

// 加密
func (s *Str) Encode(params string, secret string) string {
	//# 获取Md5密钥
	secretHex := strings.ToUpper(s.Md5(secret))
	//# 获取10进制密钥
	secretDec := s.HexToBigInt(secretHex)

	//# 将10进制密钥base64后和内容为待加密内容
	content := base64.StdEncoding.EncodeToString([]byte(secretDec.String())) + "@" + params

	//字符串转16
	infoHex := strings.ToUpper(hex.EncodeToString([]byte(content)))
	//# 转10进制
	infoDec := s.HexToBigInt(infoHex)

	//# 密文与密钥相加
	infoDec.Add(infoDec, secretDec)

	//数字变成字符串
	return infoDec.String()
}

// 解密
func (s *Str) Decode(content string, secret string) string {

	//加密密钥
	secretHex := strings.ToUpper(s.Md5(secret))
	//# 获取10进制密钥
	secretDec := s.HexToBigInt(secretHex)

	// 字符串转BigInt
	contentInt := s.StrToBigInt(content)

	//# 加密内容减密钥得到十进制密文
	private := contentInt.Sub(contentInt, secretDec)

	//十进制转字符串
	result := strings.Replace(string(private.Bytes()), "\n", "", 1)
	return result[strings.Index(result, "@")+1:]
}
