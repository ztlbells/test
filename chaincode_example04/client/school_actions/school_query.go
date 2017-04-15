package main 
import (
	"fmt"
	//"github.com/test/chaincode_example04/interactions" 
	"github.com/test/chaincode_example04/client/school"

)
func main(){
	CC_ID := "e0a47d839b8cea0765a9736bbb751ec9294e505bb0127deaf67440ba6b1c3bb64d2f4da70369dd3364eac8b18149413b28fcc10fcc6d8c82928139f2564cc4ca"
	query_return_body, _ := school.QueryChaincode_GetSchoolByAddress("alice", "ca076f5421db9454ae3ea9568c2cabe2" ,"47.90.123.204:7050", CC_ID)
	fmt.Println("query_return:", query_return_body)
}