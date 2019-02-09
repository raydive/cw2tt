package typetalk

import "github.com/nulab/go-typetalk/typetalk/v1"

func MakeTypetalkBot(token string) *v1.Client {
	client := v1.NewClient(nil)
	client.SetTypetalkToken(token)
	return client
}
