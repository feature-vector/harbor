package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	firebaseApp *firebase.App
	authClient  *auth.Client
)

func Init(ctx context.Context, credentialFile string) {
	var err error

	opt := option.WithCredentialsFile(credentialFile)
	config := &firebase.Config{}
	firebaseApp, err = firebase.NewApp(ctx, config, opt)
	if err != nil {
		panic(err)
	}
	authClient, err = firebaseApp.Auth(ctx)
	if err != nil {
		panic(err)
	}
}

func VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return authClient.VerifyIDToken(ctx, idToken)
}
