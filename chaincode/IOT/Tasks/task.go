package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "strconv"
	//"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Task struct{
	ObjectType string `json:"docType"`
	TaskId string `json:"taskId"`
	DeviceId string `json:"deviceId"`
	URI string `json:"uri"`
	IpAddress string `json:"ipaddress"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response { 
	function, args := APIstub.GetFunctionAndParameters()
	if function == "addTask" {
		return s.addTask(APIstub, args)
	} else if function == "getTask"{
		return s.getTask(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "verifyRecord"{
	        return s.verifyRecord(APIstub, args)
        }

	return shim.Error("Invalid Smart Contract function name.")
}


func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) addTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var err error
	TaskId := args[0] 
	DeviceId := args[1]
	URI := args[2]
	IpAddress := args[3] 

	ObjectType := "Task"
	var Task = &Task{ObjectType,TaskId,DeviceId,URI,IpAddress}
	TaskJSON, err := json.Marshal(Task)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(TaskId, TaskJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) verifyRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	DeviceData, _ := APIstub.GetState(args[0])
	return shim.Success(DeviceData)
}



func (s *SmartContract) getTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(buffer.Bytes())
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

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
	buffer.WriteString("]")

	return &buffer, nil
}


func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
