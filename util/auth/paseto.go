package auth

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)


var (
	PrivateKey paseto.V4AsymmetricSecretKey
	PublicKey  paseto.V4AsymmetricPublicKey
	TokenTTL   time.Duration
)

func GenerateKey() {
	privateKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := privateKey.Public()
	fmt.Println("Private Key:", base64.StdEncoding.EncodeToString(privateKey.ExportBytes()))
	fmt.Println("Public Key :", base64.StdEncoding.EncodeToString(publicKey.ExportBytes()))
}

func Init() {
	privKeyStr := os.Getenv("PASETO_PRIVATE_KEY")
	pubKeyStr := os.Getenv("PASETO_PUBLIC_KEY")

	privKeyBytes, err := base64.StdEncoding.DecodeString(privKeyStr)
	if err != nil {
		log.Fatal("Invalid private key format:", err)
	}
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		log.Fatal("Invalid public key format:", err)
	}

	PrivateKey, err = paseto.NewV4AsymmetricSecretKeyFromBytes(privKeyBytes)

	if err != nil {
		log.Fatal("Failed to import private key:", err)
	}

	PublicKey, err = paseto.NewV4AsymmetricPublicKeyFromBytes(pubKeyBytes)

	if err != nil {
		log.Fatal("Failed to import public key:", err)
	}

	TokenTTL, _ = time.ParseDuration(os.Getenv("TOKEN_EXPIRATION"))
	fmt.Println(TokenTTL)
}

func GenerateToken(id string, username string) (string, error) {
	 now := time.Now()
	 exp := now.Add(TokenTTL)
	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(exp)
	token.SetSubject("talkey-auth")

	//set keys and values
	token.Set("id", id)
	token.Set("username", username)

	signed := token.V4Sign(PrivateKey, nil)

	return signed, nil
}