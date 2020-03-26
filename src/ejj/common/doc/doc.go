package doc

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	
	"github.com/chaincode/ejj/model"
	"github.com/chaincode/ejj/common"
)

// pObj: pointer to individual object
func PersistDocAsUnique(stub shim.ChaincodeStubInterface, jsonObjStr string, pObj common.Evidence) pb.Response {
	err := common.NewObjectFromJsonString(jsonObjStr, pObj)
	if err != nil {
		res := fmt.Sprintf("Chaincode invoke[PersistDoc] failed: %s", err.Error())
		res = common.GetRetString(1, res)
		return shim.Error(res)
	}

	dt := model.DocTable{pObj.GetType(), stub}

	err = dt.SaveObjectAsUnique(pObj.GetID(), pObj)
	if err != nil {
		return shim.Error(err.Error())
	}

	res := common.GetRetByte(0, "Chaincode invoke Success")
	return shim.Success(res)
}

// pObj: pointer to individual object
func PersistReminderDocWithAppend(stub shim.ChaincodeStubInterface, jsonObjStr string, pObj common.Evidence) pb.Response {
	err := common.NewObjectFromJsonString(jsonObjStr, pObj)
	if err != nil {
		res := fmt.Sprintf("Chaincode invoke[PersistReminderDocWithAppend] failed: %s", err.Error())
		res = common.GetRetString(1, res)
		return shim.Error(res)
	}

	dt := model.DocTable{pObj.GetType(), stub}

	rr := model.Reminder{ReminderHashList: make([]string, 0)}
	
	err = dt.GetObject(pObj.GetID(), &rr)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	err = pObj.AppendString(rr.ReminderHashList...)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	err = dt.SaveObject(pObj.GetID(), pObj)
	if err != nil {
		return shim.Error(err.Error())
	}

	res := common.GetRetByte(0, "Chaincode invoke Success")
	return shim.Success(res)
}
