package common

// 赋强司法公正合约所有状态
const (
	Applied	= "applied"	// 已申请
	ExigentCautioned	= "cautioned"	// 已催告
)

//Tables
var TABLES = map[string]string{
	"notarization": "NTRN_",
	"reminder": "RMDR_",
}

type Subject struct {
	Name					string	`json:"Name"`	//姓名
	LegalityCertificateType	string	`json:"LegalityCertificateType"`	//有效证件类型
	LegalityCertificateID	string	`json:"LegalityCertificateID"`		//有效证件号码
	FixedTelephone			string	`json:"FixedTelephone,omitempty"`	//固定电话
	MobilePhone				string	`json:"MobilePhone,omitempty"`	//手机
	EMail					string	`json:"EMail,omitempty"`	//电子邮箱
	ContactAddress			string	`json:"ContactAddress,omitempty"`	//常用或联系地址
	Gender					string	`json:"Gender,omitempty"`	//性别
}

type Evidence interface {
	GetType() string
	GetID()	string
	AppendString(strs ...string) error
}