/*
	author:swb
	time:16/06/30
	MIT License
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var bankNo int = 0
var cpNo int = 0
var transactionId int = 0

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type CenterBank struct {
	Name        string
	TotalNumber int
	RestNubmer  int
}

type Bank struct {
	Name        string
	TotalNumber int
	RestNubmer  int
	ID          int
}

type Company struct {
	Name   string
	Number int
	ID     int
}

type Transaction struct {
	FromType int //Bank 0  Company 1
	FromID   int
	ToType   int //Bank 0 Company 1
	ToID     int
	Time     string
	Number   int
	ID       int
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
	var centerBank CenterBank
	totalNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	centerBank = CenterBank{Name: args[0], TotalNumber: totalNumber, RestNubmer: 0}
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
	if function == "createBank" {
		return t.createBank(stub, args)
	} else if function == "createCompany" {
		return t.createCompany(stub, args)
	} else if function == "issueCoin" {
		return t.issueCoin(stub, args)
	} else if function == "issueCoinToBank" {
		return t.issueCoinToBank(stub, args)
	} else if function == "issueCoinToCp" {
		return t.issueCoinToCp(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) createBank(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) createCompany(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) issueCoin(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) issueCoinToBank(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) issueCoinToCp(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if len(args) == 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. name of the key and value to set")
	}

	if function == "getCenterBank" {
		cbBytes, err := getCenterBank(stub, args)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return cbBytes, nil
	} else if function == "getBankById" {
		bankBytes, err := getBankById(stub, args)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return bankBytes, nil
	} else if function == "getCompanyById" {
		cpBytes, err := getCompanyById(stub, args)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return cpBytes, nil
	} else if function == "getTransactionById" {
		tsBytes, err := getTransactionById(stub, args)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return tsBytes, nil
	} else if function == "getBanks" {
		bankBytes, err := getBanks(stub)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return bankBytes, nil
	} else if function == "getCompanys" {
		cpBytes, err := getCompanys(stub)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return cpBytes, nil
	} else if function == "getTransactions" {
		tsBytes, err := getTransactions(stub)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return tsBytes, nil
	}

}

func getCenterBank(stub *shim.ChaincodeStub) (CenterBank, error) {
	var centerBank CenterBank
	cbBytes, err := stub.GetState("centerBank")
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(cbBytes, &centerBank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return cbBytes, nil
}

func getCompanyById(stub *shim.ChaincodeStub, id string) (CenterBank, error) {
	return nil, nil
}

func getBankById(stub *shim.ChaincodeStub, id string) (Bank, error) {
	return nil, nil
}

func getTransactionById(stub *shim.ChaincodeStub, id string) (CenterBank, error) {
	return nil, nil
}

func getBanks(stub *shim.ChaincodeStub) ([]Bank, error) {
	return nil, nil
}

func getCompanys(stub *shim.ChaincodeStub) ([]Company, error) {
	return nil, nil
}

func getTransactions(stub *shim.ChaincodeStub) ([]Transaction, error) {
	return nil, nil
}
