package school
import (
	"fmt"
	"errors"
	"github.com/test/chaincode_example04/rsa_functions"
	//"github.com/test/chaincode_example04/interactions" 
	"github.com/go-simplejson"
	"crypto/rsa"
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

func GetSchoolAddress_Client(school School) (string){
	return school.Address
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

	//fmt.Println(url, " [POST]\n", post)
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

func GetCCID(jsonBody string) (string, error){
	js, err := simplejson. NewJson([]byte(jsonBody))
	if err != nil {
		fmt.Println("Failed to read input bytes.")
		return "", errors.New("Failed to create NewJson.")
	}
	// {"jsonrpc":"2.0",
	//	"result":{
	//				"status":"OK",
	//				"message":"b9844a9c732801222f3e07e2b832904363823930db4b4e44cd8969a41bc80d2553514002c298157d5f4c47d325ceb9bc0dd6d8b098d0482baebf75ee826168c1"
	//			},
	//	"id":1}
	CC_ID := js.Get("result").Get("message").MustString()
	return CC_ID, nil
}


func QueryChaincode_GetSchoolByAddress(enrollId string, address string, serverAddress string, CC_ID string) (string, error) {
	url := "http://" + serverAddress + "/chaincode"
	// paramater list: Address: args[0],
    post := "{\"jsonrpc\": \"2.0\"," +
 				"\"method\": \"query\"," +
  				"\"params\": {" +
    						"\"type\": 1," +
    						"\"chaincodeID\":{" +
      						"\"name\": \"" + CC_ID +
    						"\"}," +
    			"\"ctorMsg\": {" + 
       						"\"args\":[\"getSchoolByAddress\", \"" + address + "\"]}," + 
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
/*func main(){
	SJTU_school, err := SchoolInitializer("SJTU")
	if err != nil {
		fmt.Println("School initialization failed.")
		return
	}
	SchoolInformation(SJTU_school)
	login_return_body, _ := interactions.Login("alice", "CMS10pEQlB16", "192.168.183.175:7050")
	fmt.Println("login return:", login_return_body)

	deploy_return_body, _ := DeployChaincode_CreateSchool("alice", SJTU_school, "192.168.183.175:7050", "https://github.com/ztlbells/test/chaincode_example04")
	fmt.Println("create return:", deploy_return_body)

	CC_ID, _ := GetCCID(deploy_return_body)
	fmt.Println("CC_ID:", CC_ID)

	//query_return_body, _ := QueryChaincode_GetSchoolByAddress("alice", GetSchoolAddress_Client(SJTU_school) ,"192.168.183.197:7050", CC_ID)
	//fmt.Println("query_return:", query_return_body)

}*/