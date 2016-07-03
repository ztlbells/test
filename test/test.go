/*
	author:swb
	time:16/06/30
	MIT License
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


var bankNo int = 0
var cpNo int= 0


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type CenterBank struct{
	name string
	TotalNumber int
	RestNubmer int
}

type Bank struct{
	name string
	TotalNumber int
	RestNubmer int
	ID int
}

type Company struct{
	name string
	Number int
	ID int
}

type Transaction struct{
	FromType int  //Bank 0  Company 1
	FromID int   
	ToType int   //Bank 0 Company 1 
	ToID int
	Time string
	Number int
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var totalNumber int
	totalNumber,err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	centerBank = CenterBank{name:args[0],TotalNumber:totalNumber,RestNubmer:0}
	centerBankBytes, err := json.Marshal(&centerBank)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("centerBank", centerBankBytes)
	if err != nil {
		return nil, errors.New("PutState Error" + err.Error())
	}
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, nil
}



// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	var err error
	
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. name of the key and value to set")
	}
	function = args[0]
  	
  	if args[0] == "GetCenterBank" {
  		centerBank,err := GetCenterBank(stub)
  		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		}else{
			cbBytes,err1 := json.Marshal(&centerBank)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}	
			fmt.Println("All success, returning the cp")
			return cbBytes, nil	
		}
  	}else{
  		return nil, nil
  	}
	
}


func GetCenterBank(stub *shim.ChaincodeStub) (CenterBank, error){
	var centerBank CenterBank
	cbBytes, err := stub.GetState("centerBank")	
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(cbBytes, &centerBank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
		
	return centerBank, nil
}
