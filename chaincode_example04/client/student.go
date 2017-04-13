package main 
import (
	"fmt"
	"../rsa_functions"
)

type Student struct{
	Name string
	Address string
	TXid string //transaction ID
	recordCode string // for verification
}

func StudentInitializer(name string) (Student){
	address := rsa_functions.GenerateRandomAddress()
	initializedStudent := Student {Name: name, Address: address, TXid: "NULL", recordCode: "NULL"}
	return initializedStudent
}

func StudentInformation(student Student){
	fmt.Println("Student Name:", student.Name)
	fmt.Println("Address:", student.Address)
	fmt.Println("Transaction ID:", student.TXid)
	fmt.Println("RecordCode:", student.recordCode)
}



func main(){
	ZTLUO_student := StudentInitializer("ztluo")
	StudentInformation(ZTLUO_student)

}