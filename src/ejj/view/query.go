package view

import (
	"fmt"
	"time"
	"strconv"
	"bytes"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/chaincode/ejj/common"
	"github.com/chaincode/ejj/model"
)

//QueryWithPaginationShim 分页查询数据记录
// 0 - Issuer|Drawee|Owner ; 1 - count of page ; 2 - pagination bookmark
func QueryWithPaginationShim(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 3 {
		return nil, errors.New("Chaincode query[queryBillsWithPaginationShim] failed: argument expecting 3")
	}
	
	queryString := args[0]
	//return type of ParseInt is int64
	pageSize, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return nil, err
	}
	
	bookmark := args[2]

	return getQueryResultForQueryStringWithPagination(stub, queryString, int32(pageSize), bookmark)

}

func getQueryResultForQueryStringWithPagination(stub shim.ChaincodeStubInterface, queryString string, pageSize int32, bookmark string) ([]byte, error) {
	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	bf_data, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	bf_meta := constructPaginationMetadataToQueryResults(responseMetadata)

	bf := constructJsonArray(bf_meta, bf_data)

	return bf.Bytes(), nil
}

func constructPaginationMetadataToQueryResults(responseMetadata *pb.QueryResponseMetadata) *bytes.Buffer {
	var buffer bytes.Buffer

	buffer.WriteString("{\"ResponseMetadata\":{\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}")

	return &buffer
}

//QueryAllShim 返回所有符合条件的记录
//  0 - {CouchDB selector query}
func QueryAllShim(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Chaincode query[queryAllShim] failed: argument expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	bf := constructJsonArray(buffer)

	return bf.Bytes(), nil
}

func constructJsonArray(bufs... *bytes.Buffer) *bytes.Buffer {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	for i, b := range bufs {
		if i != 0 && b.Len() > 0 {
			buffer.WriteString(",")
		}
		buffer.Write(b.Bytes())
	}
	buffer.WriteString("]")

	return &buffer
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	return &buffer, nil
}

//QueryTXChainForKeyShim 根据Key查询交易链
//  0 - Table Name; 1 - ID ;
func QueryTXChainForKeyShim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Chaincode query[queryTXChainForKeyShim] failed: argument expecting 2")
	}

	table_name := args[0]
	id := args[1]
	key := common.TABLES[table_name] + id

	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}


// QueryByKeyShim 根据key查询记录
// args: 0 - Table Name; 1 - id
func QueryByKeyShim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		res := common.GetRetString(1, "Chaincode query[queryByKeyShim] args != 2")
		return shim.Error(res)
	}
	
	tableName := args[0]
	id := args[1]
	dt := model.DocTable{tableName, stub}

	objBytes, err := dt.GetObjectBytes(id)
	if  err != nil {
		res := fmt.Sprintf("Chaincode queryByID failed: %s", err.Error())
		res = common.GetRetString(1, res)
		return shim.Error(res)
	}

	if objBytes == nil {
		return shim.Success([]byte("{}"))
	}

	return shim.Success(objBytes)
}
