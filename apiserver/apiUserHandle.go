package apiserver

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"

	"github.com/lhzd863/autoflow/internal/db"
	"github.com/lhzd863/autoflow/internal/glog"
	"github.com/lhzd863/autoflow/internal/jwt"
	"github.com/lhzd863/autoflow/internal/module"
	"github.com/lhzd863/autoflow/internal/util"

	"github.com/satori/go.uuid"
)

type ResponseResourceUser struct {
}

func NewResponseResourceUser() *ResponseResourceUser {
	return &ResponseResourceUser{}
}

func (rrs *ResponseResourceUser) SystemUserRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemUserRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER)
	defer bt.Close()

	bt.Remove(p.Id)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemUserGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemUserBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if m.UserName != p.UserName || m.Enable != "1" {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserListHandler(request *restful.Request, response *restful.Response) {

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemUserBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemUserAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.UserName) == 0 || len(p.Password) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER)
	defer bt.Close()

	m := new(module.MetaSystemUserBean)
	m.UserName = p.UserName
	m.Password = EnPwdCode(p.Password)
	m.Avatar = p.Avatar
	m.Introduction = p.Introduction
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.CreateTime = timeStr
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	u1 := uuid.Must(uuid.NewV4())
	m.Id = fmt.Sprint(u1)

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemUserAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.UserName) == 0 || len(p.Password) == 0 || len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER)
	defer bt.Close()

	fb0 := bt.Get(p.Id)
	m := new(module.MetaSystemUserBean)
	err = json.Unmarshal([]byte(fb0.(string)), &m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	m.Password = EnPwdCode(p.Password)
	m.Avatar = p.Avatar
	m.Introduction = p.Introduction
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

var PwdKey = "linkbook1qaz*WSX"

//PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//填充的反向操作，删除填充字符串
func PKCS7UnPadding1(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unpadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unpadding)], nil
	}
}

//实现加密
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//实现解密
func AesDeCrypt(cypted []byte, key []byte) (string, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = PKCS7UnPadding1(origData)
	if err != nil {
		return "", err
	}
	return string(origData), err
}

//加密base64
func EnPwdCode(pwdStr string) string {
	pwd := []byte(pwdStr)
	result, err := AesEcrypt(pwd, []byte(PwdKey))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(result)
}

//解密
func DePwdCode(pwd string) string {
	temp, _ := hex.DecodeString(pwd)
	//执行AES解密
	res, _ := AesDeCrypt(temp, []byte(PwdKey))
	return res
}

func (rrs *ResponseResourceUser) SystemRoleRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE)
	defer bt.Close()

	bt.Remove(p.Id)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRoleGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemRoleBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if m.Role != p.Role || m.Enable != "1" {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRoleListHandler(request *restful.Request, response *restful.Response) {

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemRoleBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRoleAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Role) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE)
	defer bt.Close()

	m := new(module.MetaSystemRoleBean)
	m.Role = p.Role
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.CreateTime = timeStr
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	u1 := uuid.Must(uuid.NewV4())
	m.Id = fmt.Sprint(u1)

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRoleUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Role) == 0 || len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE)
	defer bt.Close()

	fb0 := bt.Get(p.Id)
	m := new(module.MetaSystemRoleBean)
	err = json.Unmarshal([]byte(fb0.(string)), &m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRightRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRightRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_RIGHT)
	defer bt.Close()

	bt.Remove(p.Id)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRightGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRightGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_RIGHT)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemRightBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if m.Right != p.Right || m.Enable != "1" {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRightListHandler(request *restful.Request, response *restful.Response) {

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_RIGHT)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemRightBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRightAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRightAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Right) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_RIGHT)
	defer bt.Close()

	m := new(module.MetaSystemRightBean)
	m.Right = p.Right
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.CreateTime = timeStr
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	u1 := uuid.Must(uuid.NewV4())
	m.Id = fmt.Sprint(u1)

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRightUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Role) == 0 || len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_RIGHT)
	defer bt.Close()

	fb0 := bt.Get(p.Id)
	m := new(module.MetaSystemRightBean)
	err = json.Unmarshal([]byte(fb0.(string)), &m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserRoleRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemUserRoleRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER_ROLE)
	defer bt.Close()

	bt.Remove(p.Id)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserRoleGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemUserRoleGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER_ROLE)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemUserRoleBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if m.Role != p.Role || m.Enable != "1" {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserRoleListHandler(request *restful.Request, response *restful.Response) {

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER_ROLE)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemUserRoleBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserRoleAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemUserRoleAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Role) == 0 || len(p.UserName) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER_ROLE)
	defer bt.Close()

	m := new(module.MetaSystemUserRoleBean)
	m.UserName = p.UserName
	m.Role = p.Role
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.CreateTime = timeStr
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	u1 := uuid.Must(uuid.NewV4())
	m.Id = fmt.Sprint(u1)

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserRoleUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Role) == 0 || len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER_ROLE)
	defer bt.Close()

	fb0 := bt.Get(p.Id)
	m := new(module.MetaSystemUserRoleBean)
	err = json.Unmarshal([]byte(fb0.(string)), &m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)

}

func (rrs *ResponseResourceUser) SystemRoleRightRemoveHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleRightRemoveBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE_RIGHT)
	defer bt.Close()

	bt.Remove(p.Id)
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRoleRightGetHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleRightGetBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE_RIGHT)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemRoleRightBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if m.Role != p.Role || m.Enable != "1" {
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRoleRightListHandler(request *restful.Request, response *restful.Response) {

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE_RIGHT)
	defer bt.Close()

	strlist := bt.Scan()
	bt.Close()
	retlst := make([]interface{}, 0)
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemRoleRightBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				continue
			}
			retlst = append(retlst, m)
		}
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRoleRightAddHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleRightAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Role) == 0 || len(p.Right) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE_RIGHT)
	defer bt.Close()

	m := new(module.MetaSystemRoleRightBean)
	m.Role = p.Role
	m.Right = p.Right
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.CreateTime = timeStr
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	u1 := uuid.Must(uuid.NewV4())
	m.Id = fmt.Sprint(u1)

	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemRoleRightUpdateHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemRoleRightAddBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.Role) == 0 || len(p.Id) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_ROLE_RIGHT)
	defer bt.Close()

	fb0 := bt.Get(p.Id)
	m := new(module.MetaSystemRoleRightBean)
	err = json.Unmarshal([]byte(fb0.(string)), &m)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parse json error.%v", err), nil)
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	m.UpdateTime = timeStr
	m.Enable = p.Enable
	jsonstr, _ := json.Marshal(m)
	err = bt.Set(m.Id, string(jsonstr))
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("data in db update error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("data in db update error.%v", err), nil)
		return
	}
	retlst := make([]interface{}, 0)
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) SystemUserTokenHandler(request *restful.Request, response *restful.Response) {
	p := new(module.MetaParaSystemUserTokenBean)
	err := request.ReadEntity(&p)
	if err != nil {
		glog.Glog(LogF, fmt.Sprintf("Parse json error.%v", err))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("Parse json error.%v", err), nil)
		return
	}
	if len(p.UserName) == 0 || len(p.Password) == 0 {
		glog.Glog(LogF, fmt.Sprintf("parameter missed."))
		util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("parameter missed."), nil)
		return
	}

	bt := db.NewBoltDB(conf.BboltDBPath+"/"+util.FILE_AUTO_SYS_DBSTORE, util.TABLE_AUTO_SYS_USER)
	defer bt.Close()

	retlst := make([]interface{}, 0)
	strlist := bt.Scan()
	bt.Close()
	for _, v := range strlist {
		for _, v1 := range v.(map[string]interface{}) {
			m := new(module.MetaSystemUserBean)
			err := json.Unmarshal([]byte(v1.(string)), &m)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("get cmd failed.%v", err), nil)
				return
			}
			if m.UserName != p.UserName || m.Enable != "1" {
				continue
			}
			enpassword := EnPwdCode(p.Password)
			if enpassword != m.Password {
				util.ApiResponse(response.ResponseWriter, 700, "password entre err.", nil)
				return
			}
			n := new(module.MetaSystemUserTokenBean)
			n.Token, err = rrs.createToken(30 * 24)
			if err != nil {
				glog.Glog(LogF, fmt.Sprint(err))
				util.ApiResponse(response.ResponseWriter, 700, fmt.Sprintf("create token err.%v", err), nil)
				return
			} else {
				retlst = append(retlst, n)
				util.ApiResponse(response.ResponseWriter, 200, "", retlst)
				return
			}
		}
	}
	if len(retlst) == 0 {
		util.ApiResponse(response.ResponseWriter, 700, "username entre err.", nil)
		return
	}
	util.ApiResponse(response.ResponseWriter, 200, "", retlst)
}

func (rrs *ResponseResourceUser) createToken(keeptime int64) (string, error) {
	exp1, _ := strconv.ParseFloat(fmt.Sprintf("%v", time.Now().Unix()+3600*keeptime), 64)
	claims := map[string]interface{}{
		"iss": conf.JwtKey,
		"exp": exp1,
	}
	key := []byte(conf.JwtKey)
	tokenbyte, err := jwt.Encode(
		claims,
		key,
		"HS256",
	)

	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to encode: %v", err))
	}
	return string(tokenbyte), nil
}
