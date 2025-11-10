package encrypt

import "github.com/speps/go-hashids/v2"

var (
	hashSalt = "ecomsecret"
	hashLen  = 8
)

type DigitalInterface interface {
	Encrypt(num int) string
	Decrypt(str string) int
}
type Digital struct {
	Length int
	Salt   string
}

//NewDigitalSecret
//@Description: 实例化digital配置
//@param length
//@param salt
//@return *Digital
func NewDigitalSecret(length int, salt string) *Digital {
	return &Digital{
		Length: length,
		Salt:   salt,
	}
}

//Encrypt
//@Description: 数字加密
//@receiver d
//@param num 数字加密为字符串
//@return string
func (d *Digital) Encrypt(num int) string {
	hd := hashids.NewData()
	hd.MinLength = d.Length
	hd.Salt = d.Salt

	hid, _ := hashids.NewWithData(hd)
	secret, _ := hid.Encode([]int{num})

	return secret
}

//Decrypt
//@Description: 字符串解密为数字
//@receiver d
//@param encryptStr
//@return int
func (d *Digital) Decrypt(encryptStr string) int {
	hd := hashids.NewData()
	hd.MinLength = d.Length
	hd.Salt = d.Salt
	hid, _ := hashids.NewWithData(hd)
	decryptVal, _ := hid.DecodeWithError(encryptStr)
	num := 0
	if len(decryptVal) > 0 {
		num = decryptVal[0]
	}
	return num
}

//Encrypt
//@Description:hashId加密
//@receiver h
//@param num
//@return string
func Encrypt(num int) string {
	hd := hashids.NewData()
	hd.MinLength = hashLen
	hd.Salt = hashSalt

	hid, _ := hashids.NewWithData(hd)
	secret, _ := hid.Encode([]int{num})

	return secret
}

//Decrypt
//@Description: hashId解密
//@receiver h
//@param encryptStr
//@return int
func Decrypt(encryptStr string) int {
	hd := hashids.NewData()
	hd.MinLength = hashLen
	hd.Salt = hashSalt
	hid, _ := hashids.NewWithData(hd)
	decryptVal, _ := hid.DecodeWithError(encryptStr)
	num := 0
	if len(decryptVal) > 0 {
		num = decryptVal[0]
	}
	return num
}
