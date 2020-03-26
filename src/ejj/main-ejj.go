package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/chaincode/ejj/view"
	"github.com/chaincode/ejj/model"
	"github.com/chaincode/ejj/common"
	"github.com/chaincode/ejj/common/doc"
)

//EnhancementJudiciaryJustice 链码基本结构
type EJJ struct {
}

//Init chaincode基本接口
func (ejj *EJJ) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

//Invoke chaincode基本接口
func (ejj *EJJ) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "takeNotarization" {
		// 公证书证据上链
		return ejj.takeNotarization(stub, args)
	} else if function == "takeReminder" {
		// 催告记录和执行证书的哈希上链
		return ejj.takeReminder(stub, args)
	} else if function == "queryByKey" {
		// 按唯一key查询单条记录
		return ejj.queryByKey(stub, args)
	} else if function == "queryAll" {
		// 按条件查询多条记录
		return ejj.queryAll(stub, args)
	} else if function == "queryWithPagination" {
		// 按条件分页查询
		return ejj.queryWithPagination(stub, args)
	} else if function == "queryTXChainForKey" {
		// 查询key的交易历史
		return ejj.queryTXChainForKey(stub, args)
	}

	res := common.GetRetString(1, "Chaincode Unkown method!")

	return shim.Error(res)
}

//takeNotarization 公证书证据上链
// args: 0 - {Notarization Object}
func (ejj *EJJ) takeNotarization(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		res := common.GetRetString(1, "Chaincode invoke[takeNotarization] args != 1")
		return shim.Error(res)
	}

	var nn model.Notarization
	return doc.PersistDocAsUnique(stub, args[0], &nn)
}

//takeReminder 催告记录和执行证书的哈希上链
// args: 0 - {Reminder Object}
func (ejj *EJJ) takeReminder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		res := common.GetRetString(1, "Chaincode invoke[takeReminder] args != 1")
		return shim.Error(res)
	}

	rr := model.Reminder{ReminderHashList: make([]string, 0)}
	return doc.PersistReminderDocWithAppend(stub, args[0], &rr)
}

//queryWithPagination 分页查询数据记录
//  0 - {CouchDB selector query} ; 1 - count of page ; 2 - pagination bookmark
func (ejj *EJJ) queryWithPagination(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	queryResults, err := view.QueryWithPaginationShim(stub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//queryAll 返回所有符合条件的记录
//  0 - {CouchDB selector query}
func (ejj *EJJ) queryAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return view.QueryAllShim(stub, args)
}

//queryTXChainForKey 根据Key查询交易链
//  0 - Table Name; 1 - ID ;
func (t *EJJ) queryTXChainForKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return view.QueryTXChainForKeyShim(stub, args)
}

// 根据ID查询记录
// args: 0 - Table Name; 1 - id
func (ejj *EJJ) queryByKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return view.QueryByKeyShim(stub, args)
}

func main() {
	if err := shim.Start(new(EJJ)); err != nil {
		fmt.Printf("Error starting EJJ: %s", err)
	}
}
