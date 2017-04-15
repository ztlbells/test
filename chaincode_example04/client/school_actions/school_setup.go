package main 
import (
	"fmt"
	"github.com/test/chaincode_example04/interactions" 
	"github.com/test/chaincode_example04/client/school"

)

func main(){
	SJTU_school, err := school.SchoolInitializer("SJTU")
	if err != nil {
		fmt.Println("School initialization failed.")
		return
	}
	school.SchoolInformation(SJTU_school)
	login_return_body, _ := interactions.Login("alice", "CMS10pEQlB16", "47.90.123.204:7050")
	fmt.Println("login return:", login_return_body)

	deploy_return_body, _ := school.DeployChaincode_CreateSchool("alice", SJTU_school, "47.90.123.204:7050", "https://github.com/ztlbells/test/chaincode_example04")
	fmt.Println("create return:", deploy_return_body)

	CC_ID, _ := school.GetCCID(deploy_return_body)
	fmt.Println("CC_ID:", CC_ID)
}