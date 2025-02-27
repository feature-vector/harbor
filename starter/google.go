package starter

import "os"

func InitGoogle(credentialsPath string) {
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credentialsPath)
	if err != nil {
		panic(err)
	}
}
