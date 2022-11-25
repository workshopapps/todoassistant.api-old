package firebaseinit

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func SetupFirebase() (*firebase.App, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile("./taskman-firebase-adminsdk.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Println("Unable to Connect To Firebase", err)
		return nil, err
	}

	_, err = app.Messaging(ctx)
	if err != nil {
		fmt.Println("Unable to Connect To Firebase", err)
		return nil, err
	}

	return app, nil
}
