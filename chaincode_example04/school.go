package main 
import (
	"fmt"
	"./rsa_functions"
	"crypto/rsa"
	"./interactions"
	"bytes"
    "io/ioutil"
    "net/http"
)

type School struct{
	Name string
	Address string
	PriKey *rsa.PrivateKey
	PubKey *rsa.PublicKey
	StudentAddress []string
}


func SchoolInitializer(name string) (School, error){
	keypair, err := rsa_functions.GenerateKeyPair (1024) 
	if err != nil {
		fmt.Println("Key pair generation failed.")
		// cannot handle the error here, may change it later
	}

	var stuAddress []string
	address := rsa_functions.GenerateRandomAddress()
	initializedSchool := School {Name: name, Address: address, PriKey: keypair.
						PriKey, PubKey: keypair.PubKey, StudentAddress: stuAddress}
	return initializedSchool, err
}



func SchoolInformation(school School){
	fmt.Println("=======================================")
	fmt.Println("School Name:", school.Name)
	fmt.Println("Address:", school.Address)
	fmt.Println("PubKey - N (modulus):", school.PubKey.N)
	fmt.Println("PubKey - E (exponent):", school.PubKey.E)
	fmt.Println("=======================================")
}

func DeployChaincodeCreateSchool (enrollId string, school School, chaincodePath string) (string error){
	
}

func main(){
	SJTU_school, err := SchoolInitializer("SJTU")
	if err != nil {
		fmt.Println("School initialization failed.")
		return
	}
	SchoolInformation(SJTU_school)
	return_body, _ := interactions.Login("alice", "CMS10pEQlB16", "47.90.123.204:7050")
	fmt.Println("login return:", return_body)
}