package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"go_mxshop_srvs/user_srv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"strings"
)

func genMd5(s string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, s)
	return hex.EncodeToString(hash.Sum(nil))
}

var passwordOptions = &password.Options{16, 100, 32, sha512.New}

func genNewMd5(s string) string {
	// Using custom options
	salt, encodedPwd := password.Encode(s, passwordOptions)
	// 通过$符号组合加密算法、盐值、加密后数据
	newEncodePwd := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	// 验证数据
	// TODO 验证数据长度
	pwdOption := strings.Split(newEncodePwd, "$")
	check := password.Verify(s, pwdOption[2], pwdOption[3], passwordOptions)
	fmt.Println(check)             // true
	fmt.Println(len(newEncodePwd)) // 96
	return newEncodePwd
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:chengfei@20170917.com@tcp(116.62.167.224:3306)/go_mxshop_srvs_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	fmt.Println(genNewMd5("123"))
}
