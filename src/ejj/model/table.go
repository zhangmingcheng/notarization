package model

import (
	"encoding/json"
	"fmt"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/chaincode/ejj/common"
)


type DocTable struct{
	Name string
	Stub shim.ChaincodeStubInterface
}

func (dt DocTable) composeKey(id string) string {
	return common.TABLES[dt.Name] + id
}

func (dt DocTable) GetObjectBytes(id string) ([]byte, error) {
	if !dt.isTableExist() {
		return nil, errors.New(fmt.Sprintf("The table[%s] not exist", dt.Name))
	}
	
	obj_bytes, err := dt.Stub.GetState(dt.composeKey(id))

	if err != nil {
		return nil, err
	}

	return obj_bytes, nil

}

// pObj: pointer to individual object
func (dt DocTable) GetObject(id string, pObj interface{}) error {
	obj_bytes, err := dt.GetObjectBytes(id)
	
	if err != nil {
		return err
	} else if obj_bytes != nil {
		err = json.Unmarshal(obj_bytes, pObj)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (dt DocTable) SaveObject(id string, obj interface{}) error {
	if !dt.isTableExist() {
		return errors.New(fmt.Sprintf("The table[%s] not exist", dt.Name))
	}
	
	obj_bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = dt.Stub.PutState(dt.composeKey(id), obj_bytes)
	if err != nil {
		return err
	}
	
	return nil
}

func (dt DocTable) SaveObjectAsUnique(id string, obj interface{}) error {
	exist, err := dt.isObjectExist(id)
	if err != nil {
		return err
	} else if exist == true {
		return errors.New(fmt.Sprintf("The record[%s] is existing", id))
	}
	
	return dt.SaveObject(id, obj)
}


func (dt DocTable) isObjectExist(id string) (bool, error) {
	obj_bytes, err := dt.GetObjectBytes(id)
	
	if err != nil {
		return false, err
	} else if obj_bytes == nil {
		return false, nil
	}
	
	return true, nil
}

func (dt DocTable) isTableExist() bool {
	_, exist := common.TABLES[dt.Name]
	return exist
}