package main 
import (
	"fmt"
	"os"
)

type School struct{
	Name string
	Location string
	Address string
	PriKey string
	PubKey string
	StudentAddress []string
}

func main(){
	keypair, err := rsa_functions.generateKeyPair (1024) 
}