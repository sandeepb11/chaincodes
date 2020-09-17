package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"strconv"
	//"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)



// Define the Smart Contract structure
type SmartContract struct {
}


// Define the structure
type Device struct{
	ObjectType string `json:"docType"`
	DeviceId string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
	DeviceOwnerId string `json:"deviceOwnerId"`
	DeviceAddress string `json:"deviceAddress"`
	GpsId string `json:"gpsId"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


//Invoke Method
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response { 
	function, args := APIstub.GetFunctionAndParameters()
	if function == "getDeviceRecords" {
		return s.getDeviceRecords(APIstub, args)
	} else if function == "createDevice" {
		return s.createDevice(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "verifyDevice"{
	        return s.verifyDevice(APIstub, args)
        }
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}
// Device parameters ==>
// 	DeviceId string `json:"deviceId"`
// 	DeviceName string `json:"deviceName"`
// 	DeviceOwnerId string `json:"deviceOwnerId"`
// 	DeviceAddress string `json:"deviceAddress"`
// 	GpsId string `json:"gpsId"`

func (s *SmartContract) createDevice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var err error
	DeviceId := args[0] 
	DeviceName := args[1]
	DeviceOwnerId := args[2]
	DeviceAddress := args[3] 
	GpsId := args[4]

	ObjectType := "Device"
	var Device = &Device{ObjectType,DeviceId,DeviceName,DeviceOwnerId,DeviceAddress,GpsId}
	DeviceJSON, err := json.Marshal(Device)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(DeviceId, DeviceJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) verifyDevice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	DeviceData, _ := APIstub.GetState(args[0])
	return shim.Success(DeviceData)
}



//Device data

func (s *SmartContract) getDeviceRecords(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

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
