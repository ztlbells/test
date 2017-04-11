/* testing functions from rsa package */
package rsa_functions
import (
	"os"
	"fmt"
	"crypto"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
)


type KeyPair struct{
	PriKey *rsa.PrivateKey
	PubKey *rsa.PublicKey
}

//encryption
func GenerateCiphertext (bits int, secretMessage []byte, label []byte) ([]byte, error){

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate key pairs.")
		return nil, err
	}

	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &privateKey.PublicKey, secretMessage, label)
	if err != nil {
	        fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
	        return nil, err
	}
	return ciphertext, err
}

//generate pri/pub key pair
func GenerateKeyPair (bits int) (KeyPair, error){
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println("Failed to generate key pairs.")
	}
	keypair := KeyPair {PriKey: privateKey, PubKey: &privateKey.PublicKey}
	return keypair, err

}

// generate address: school/student identifier
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

func Test(){
	// CIPHERTEXT
	secretMessage := []byte ("The patrol car is in pursuit.")
	/*label := []byte ("order")

	ciphertext, err := generateCiphertext (2048, secretMessage, label)
	if err != nil {
		fmt.Println("Failed to generate ciphertext.")
		return
	}
	fmt.Println("Ciphertext:", ciphertext)*/

	// KEYPAIR GENERATION
	keypair, err := GenerateKeyPair(1024)
	if err != nil {
		fmt.Println("Failed to generate keypair.")
		return
	}
	//fmt.Println("privateKey", keypair.priKey)
	//fmt.Println("publicKey", keypair.pubKey)

	//SIGN MESSAGE
	signature_hashed := sha256.Sum256(secretMessage)
	signature, err := rsa.SignPKCS1v15(rand.Reader, keypair.PriKey, crypto.SHA256, signature_hashed[:])
	if err != nil {
	        fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
	        return
	}
	fmt.Printf("Signature: %x\n", signature)

	verification_hashed := sha256.Sum256(secretMessage)

	err = rsa.VerifyPKCS1v15(keypair.PubKey, crypto.SHA256, verification_hashed[:], signature)
	if err != nil {
	        fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
	        return
}
}

