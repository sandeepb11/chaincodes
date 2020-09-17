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

type Sensor struct{
	ObjectType string `json:"docType"`
	SensorId string `json:"sensorId"`
	Name string `json:"name"`
	Unit string `json:"unit"`
	Thresold string `json:"thresold"`
	DeviceId string `json:"deviceId"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response { 
	function, args := APIstub.GetFunctionAndParameters()
	if function == "addSensor" {
		return s.addSensor(APIstub, args)
	} else if function == "getDevice" {
		return s.getDevice(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "verifyRecord" {
        	return s.VerifyRecord(APIstub, args)
        }

	return shim.Error("Invalid Smart Contract function name.")
}


func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


//Sensor data (Register)
// SensorId string `json:"sensorId"`
// Name string `json:"name"`
// Unit string `json:"unit"`
// Thresold string `json:"thresold"`
// DeviceId string `json:"deviceId"`

func (s *SmartContract) addSensor(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var err error
	SensorId := args[0] 
	Name := args[1]
	Unit := args[2]
	Thresold := args[3] 
	DeviceId := args[4]

	ObjectType := "Sensor"
	var Sensor = &Sensor{ObjectType,SensorId,Name,Unit,Thresold,DeviceId}
	SensorLogJSON, err := json.Marshal(Sensor)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(SensorId, SensorLogJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) VerifyRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	DeviceData, _ := APIstub.GetState(args[0])
	return shim.Success(DeviceData)
}

func (s *SmartContract) getDevice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
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
