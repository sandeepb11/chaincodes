package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"bytes"
	//"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the structure
type Record struct {
	ObjectType    string `json:"docType"`
	Orderid string `json:"uid"`
	ProductId         string `json:"productid"`
	Price         string `json:"price"`
	ManufacturerId  string `json:"manufacturerid"`
	DistributorId      string `json:"distributorid"`
	RetailerId          string `json:"retailerid"`
	ModifiedBy string `json:"modifiedby"`
	Orderstatus        string `json:"ordererstatus"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

//Invoke Method
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "publishContract" {
		return s.publishContract(APIstub, args)
	} else if function == "update" {
		return s.update(APIstub, args)
	} else if function == "verifyRecord" {
		return s.verifyRecord(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "gettingFullRecord" {
		return s.gettingFullRecord(APIstub, args)
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
func (s *SmartContract) publishContract(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	var err error

	Orderid := args[0]
	ProductId := args[1]
	Price := args[2]
	ManufacturerId:= args[3]
	DistributorId:= args[4]
	RetailerId:= args[5]
	ModifiedBy:= args[6]
	Orderstatus := args[7]

	objectType := "Record"
	Record := &Record{objectType, Orderid, ProductId, Price, ManufacturerId, DistributorId, RetailerId, ModifiedBy, Orderstatus}
	RecordJSON, err := json.Marshal(Record)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(Orderid, RecordJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
func (s *SmartContract) update(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	Orderid := args[0]
	ordererstatus := strings.ToLower(args[1])

	RecordJSON, err := APIstub.GetState(Orderid)
	if err != nil {
		return shim.Error("Failed to get Record :" + err.Error())
	} else if RecordJSON == nil {
		return shim.Error("Record does not exist")
	}

	update := Record{}
	err = json.Unmarshal(RecordJSON, &update) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	update.Orderstatus = ordererstatus //change the owner
	RecordJSON1, _ := json.Marshal(update)
	err = APIstub.PutState(Orderid, RecordJSON1)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
func (s *SmartContract) gettingFullRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

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

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

