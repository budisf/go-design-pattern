package config

type Conf struct {
	App struct {
		Name       string `env:"APP_NAME"`
		Name_api   string `env:"APP_NAME_API"`
		Port       string `env:"APP_PORT"`
		Mode       string `env:"APP_MODE"`
		Url        string `env:"APP_URL"`
		Secret_key string `env:"APP_SECRET"`
	}
	Db struct {
		Host string `env:"DB_HOST_LOCAL"`
		Name string `env:"DB_NAME_LOCAL"`
		User string `env:"DB_USER_LOCAL"`
		Pass string `env:"DB_PASSWORD_LOCAL"`
		Port string `env:"DB_PORT_LOCAL"`
	}
	Db_staging struct {
		Host string `env:"DB_HOST_STAGING"`
		Name string `env:"DB_NAME_STAGING"`
		User string `env:"DB_USER_STAGING"`
		Pass string `env:"DB_PASSWORD_STAGING"`
		Port string `env:"DB_PORT_STAGING"`
	}
	Db_prod struct {
		Host string `env:"DB_HOST_PROD"`
		Name string `env:"DB_NAME_PROD"`
		User string `env:"DB_USER_PROD"`
		Pass string `env:"DB_PASSWORD_PROD"`
		Port string `env:"DB_PORT_PROD"`
	}
	Firebase_staging struct {
		Type                        string `env:"FIREBASE_TYPE"`
		Project_id                  string `env:"FIREBASE_PROJECT_ID"`
		Private_key_id              string `env:"FIREBASE_PRIVATE_KEY_ID"`
		Private_key                 string `env:"FIREBASE_PRIVATE_KEY"`
		Client_email                string `env:"FIREBASE_CLIENT_EMAIL"`
		Client_id                   string `env:"FIREBASE_CLIENT_ID"`
		Auth_uri                    string `env:"FIREBASE_AUTH_URL"`
		Token_uri                   string `env:"FIREBASE_TOKEN_URL"`
		Auth_provider_x509_cert_url string `env:"FIREBASE_AUTH_PROVIDER_X509_CERT_URL"`
		Client_x509_cert_url        string `env:"FIREBASE_CLIENT_X509_CERT_URL"`
	}

	Firebase_prod struct {
		Type                        string `env:"FIREBASE_TYPE"`
		Project_id                  string `env:"FIREBASE_PROJECT_ID"`
		Private_key_id              string `env:"FIREBASE_PRIVATE_KEY_ID"`
		Private_key                 string `env:"FIREBASE_PRIVATE_KEY"`
		Client_email                string `env:"FIREBASE_CLIENT_EMAIL"`
		Client_id                   string `env:"FIREBASE_CLIENT_ID"`
		Auth_uri                    string `env:"FIREBASE_AUTH_URL"`
		Token_uri                   string `env:"FIREBASE_TOKEN_URL"`
		Auth_provider_x509_cert_url string `env:"FIREBASE_AUTH_PROVIDER_X509_CERT_URL"`
		Client_x509_cert_url        string `env:"FIREBASE_CLIENT_X509_CERT_URL"`
	}
}
