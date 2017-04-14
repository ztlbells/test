package main 
import (
	"fmt"
	"github.com/test/chaincode_example04/rsa_functions"
	"crypto/rsa"
	"github.com/test/chaincode_example04/interactions" //for login
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
	keypair, err := rsa_functions.GenerateKeyPair (256) 
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


func DeployChaincode_CreateSchool(enrollId string, school School, serverAddress string, chaincodePath string) (string, error){
	url := "http://" + serverAddress + "/chaincode"
	// paramater list: Name:args[0], Address: args[1], PriKey:args[2], PubKey:args[3]
    post := "{\"jsonrpc\": \"2.0\"," +
 				"\"method\": \"deploy\"," +
  				"\"params\": {" +
    						"\"type\": 1," +
    						"\"chaincodeID\":{" +
      						"\"path\": \"" + chaincodePath +
    						"\"}," +
    			"\"ctorMsg\": {" + 
       						"\"args\":[\"init\", \"createSchool\", \"" + school.Name + "\"," +
       						"\"" + string(school.Address) + "\", \"" + rsa_functions.GetMarshalledPriKeyString(school.PriKey) + "\"," + 
    						"\"" + rsa_functions.GetMarshalledPubKeyString(school.PubKey) + "\"]}," + 
    			"\"secureContext\": \"" + enrollId + "\"" +
  				"}," + 
  				"\"id\": 1" +
			"}"

	fmt.Println(url, " [POST]\n", post)
    var jsonStr = []byte(post)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    return string(body), err
    return ".", nil
}

func main(){
	SJTU_school, err := SchoolInitializer("SJTU")
	if err != nil {
		fmt.Println("School initialization failed.")
		return
	}
	SchoolInformation(SJTU_school)
	login_return_body, _ := interactions.Login("alice", "CMS10pEQlB16", "47.90.123.204:7050")
	fmt.Println("login return:", login_return_body)

	/*marshalled_pri := GetMarshalledPriKey(SJTU_school.PriKey)
	marshalled_pub := GetMarshalledPubKey(SJTU_school.PubKey)
	
	fmt.Println("pub:", marshalled_pub)

	string_marshalled_pub := string(marshalled_pub)
	fmt.Println("string_pub:", string_marshalled_pub)

	
	//fmt.Println("pub:", marshalled_pub)*/



	deploy_return_body, _ := DeployChaincode_CreateSchool("alice", SJTU_school, "47.90.123.204:7050", "https://github.com/ztlbells/test/chaincode_example04")
	fmt.Println("login return:", deploy_return_body)
}