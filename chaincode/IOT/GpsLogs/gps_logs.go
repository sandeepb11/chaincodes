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

type GpsLog struct{
	ObjectType string `json:"docType"`
	GpsID string `json:"sensorId"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
	TimeStamp string `json:"timestamp"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response { 
	function, args := APIstub.GetFunctionAndParameters()
	if function == "addGpsData" {
		return s.addGpsData(APIstub, args)
	} else if function == "queryGpsLog" {
		return s.queryGpsLog(APIstub, args)
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


func (s *SmartContract) verifyRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	RecordData, _ := APIstub.GetState(args[0])
	return shim.Success(RecordData)
}



func (s *SmartContract) queryGpsLog(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
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

//add Gps data
// GpsID string `json:"gpsId"`
// Latitude string `json:"latitude"`
// Longitude string `json:"longitude"`
// TimeStamp string `json:"timestamp"`


func (s *SmartContract) addGpsData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var err error
	GpsID := args[0] 
	Latitude := args[1]
	Longitude := args[2]
	TimeStamp := args[3] 

	ObjectType := "GpsLog"
	var GpsLog = &GpsLog{ObjectType,GpsID,Latitude,Longitude,TimeStamp}
	GpsJSON, err := json.Marshal(GpsLog)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(GpsID, GpsJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
