package main 
import (
	"fmt"
	"./rsa_functions"
	"crypto/rsa"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
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
	initializedSchool := School {Name: name, Address:GenerateRandomAddress(), PriKey: keypair.
						PriKey, PubKey: keypair.PubKey, StudentAddress: stuAddress}
	return initializedSchool, err
}

func SchoolInformation(school School){
	fmt.Println("School Name:", school.Name)
	fmt.Println("Address:", school.Address)
	fmt.Println("PubKey - N (modulus):", school.PubKey.N)
	fmt.Println("PubKey - E (exponent):", school.PubKey.E)
}

func GenerateRandomAddress() (string){
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))

	address := hex.EncodeToString(h.Sum(nil))
	return address
}

func main(){
	SJTU_school, err := SchoolInitializer("SJTU")
	if err != nil {
		fmt.Println("School initialization failed.")
		return
	}
	SchoolInformation(SJTU_school)

}