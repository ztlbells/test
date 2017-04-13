package interactions
import (
    "bytes"
    //"fmt"
    "io/ioutil"
    "net/http"
)

//     enrollId := "alice" 
//    enrollSecret := "CMS10pEQlB16"

func Login(enrollId string, enrollSecret string, serverAddress string) (string, error){
	url := "http://" + serverAddress + "/registrar"
    post := "{\"enrollId\":\"" +  enrollId + "\",\"enrollSecret\":\"" + enrollSecret + "\"}"
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
}

