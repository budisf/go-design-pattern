package firebase

import (
	"context"
	"encoding/json"
	"ethical-be/app/config"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

var Firebase_Credential map[string]string

func FirebaseCredentialInit(conf *config.Conf) {
	switch conf.App.Mode {
	case "staging":
		Firebase_Credential = map[string]string{
			"type":                        conf.Firebase_staging.Type,
			"project_id":                  conf.Firebase_staging.Project_id,
			"private_key_id":              conf.Firebase_staging.Private_key_id,
			"private_key":                 strings.Replace(string(conf.Firebase_staging.Private_key), "\\n", "\n", -1),
			"client_email":                conf.Firebase_staging.Client_email,
			"client_id":                   conf.Firebase_staging.Client_id,
			"auth_uri":                    conf.Firebase_staging.Auth_uri,
			"token_uri":                   conf.Firebase_staging.Token_uri,
			"auth_provider_x509_cert_url": conf.Firebase_staging.Auth_provider_x509_cert_url,
			"client_x509_cert_url":        conf.Firebase_staging.Client_x509_cert_url,
		}
	case "production":
		Firebase_Credential = map[string]string{
			"type":                        conf.Firebase_prod.Type,
			"project_id":                  conf.Firebase_prod.Project_id,
			"private_key_id":              conf.Firebase_prod.Private_key_id,
			"private_key":                 strings.Replace(string(conf.Firebase_prod.Private_key), "\\n", "\n", -1),
			"client_email":                conf.Firebase_prod.Client_email,
			"client_id":                   conf.Firebase_prod.Client_id,
			"auth_uri":                    conf.Firebase_prod.Auth_uri,
			"token_uri":                   conf.Firebase_prod.Token_uri,
			"auth_provider_x509_cert_url": conf.Firebase_prod.Auth_provider_x509_cert_url,
			"client_x509_cert_url":        conf.Firebase_prod.Client_x509_cert_url,
		}
	default:
		Firebase_Credential = map[string]string{
			"type":                        conf.Firebase_staging.Type,
			"project_id":                  conf.Firebase_staging.Project_id,
			"private_key_id":              conf.Firebase_staging.Private_key_id,
			"private_key":                 strings.Replace(string(conf.Firebase_staging.Private_key), "\\n", "\n", -1),
			"client_email":                conf.Firebase_staging.Client_email,
			"client_id":                   conf.Firebase_staging.Client_id,
			"auth_uri":                    conf.Firebase_staging.Auth_uri,
			"token_uri":                   conf.Firebase_staging.Token_uri,
			"auth_provider_x509_cert_url": conf.Firebase_staging.Auth_provider_x509_cert_url,
			"client_x509_cert_url":        conf.Firebase_staging.Client_x509_cert_url,
		}
	}

}

func FirebaseInit() *firebase.App {
	var cd []byte

	cd, _ = json.Marshal(Firebase_Credential)

	// opt := option.WithCredentialsFile("D:/Go/src/Golang-Fiber/serviceAccountKey.json")
	opt := option.WithCredentialsJSON(cd)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.SetFlags(16)
		log.Printf("error initializing app: %v\n", err)
	}
	return app
}

func CloudFirestore() *firestore.Client {
	app := FirebaseInit()
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.SetFlags(16)
		log.Printf("error initializing client: %v\n", err)
	}
	return client
}

func AuthClient() *auth.Client {
	app := FirebaseInit()
	client, err := app.Auth(context.Background())
	if err != nil {
		log.SetFlags(16)
		log.Printf("error getting Auth client: %v\n", err)
	}

	return client
}

func StorageClient() *storage.Client {
	var cd []byte
	var config *firebase.Config
	cd, _ = json.Marshal(Firebase_Credential)
	config = &firebase.Config{
		StorageBucket: os.Getenv("FIREBASE_STORAGE_BUCKET_URL"),
	}

	opt := option.WithCredentialsJSON(cd)
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.SetFlags(16)
		log.Println(err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.SetFlags(16)
		log.Println(err)
	}
	return client
}
