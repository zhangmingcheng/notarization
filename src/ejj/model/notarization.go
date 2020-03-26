package model

import(
	"github.com/chaincode/ejj/common"
)

//Notarization 公证书结构
type Notarization struct {
	Type				string				`json:"Type"`					//类别，类似表名，区分大小写
	NotaryMatter		string				`json:"NotaryMatter"`			//公证事项
	NotaryInstitution	string				`json:"NotaryInstitution"`		//公证单位
	NotaryName			string				`json:"NotaryName"`				//公证员名称
	NotarizationID		string				`json:"NotarizationID"`			//公证书编号
	NotarizationHash	string				`json:"NotarizationHash"`		//公证书哈希
	NotarizationNO		string				`json:"NotarizationNO"`			//公证书号，如（2019）京中信电证字00665号
	LitigantList		[]common.Subject	`json:"LitigantList,omitempty"`	//当事人列表
	ContentDocHashList	[]string			`json:"ContentDocHashList"`		//申请公正的相关资料文件哈希列表，包含贷款合同、公证申请表、赋强公证效力告知书、电子询问笔录、送达地址确认书
	CustomerID			string				`json:"CustomerID,omitempty"`	//客户编号
	CustomerName		string				`json:"CustomerName,omitempty"`	//客户名称
	MerchantID			string				`json:"MerchantID,omitempty"`	//商户编号
	MerchantName		string				`json:"MerchantName,omitempty"`	//商户名称
	DebtID				string				`json:"DebtID"`					//债项编号
	AcceptDate			int64				`json:"AcceptDate"`				//受理时间
	NotarizationDate	int64				`json:"NotarizationDate"`		//出证时间
	CreateDate			int64				`json:"CreateDate"`				//创建时间
	Debtor				common.Subject		`json:"Debtor,omitempty"`		//借款人信息
}

func (nn Notarization) GetType() string {
	return nn.Type
}

func (nn Notarization) GetID() string {
	return nn.NotarizationID
}

func (nn *Notarization) AppendString(strs ...string) error {
	return nil
}