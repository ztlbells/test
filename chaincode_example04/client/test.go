package main
import ("fmt")

func main(){
	 post := "{\"jsonrpc\": \"2.0\"," + "\n" +
 				"\"method\": \"deploy\"," + "\n" +
  				"\"params\": {" + "\n" +
    						"	\"type\": 1," + "\n" +
    						"	\"chaincodeID\":{" + "\n" +
      						"	\"path\": \"" + "[1]\"" + "\n"  +
    						"}," + "\n"  +
    			"\"ctorMsg\": {" + "\n" +
       						"	\"args\":[\"init\", \"createSchool\", \"" + "[2]" + "\"," + "\n"  +
       						"	\"" + "[3]" + "\", \"" + "[4]" + "\"," + "\n" +
    						"	\"" + "[5]" + "\"]}," + "\n" +
    			"\"secureContext\": \"" + "[6]" + "\"" + "\n" +
  				"}," + "\n" +
  				"\"id\": 1" + "\n" +
			"}" 
		fmt.Println(post)

}