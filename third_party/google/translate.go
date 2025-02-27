package google

import (
	"cloud.google.com/go/translate"
	"context"
	"errors"
	"golang.org/x/text/language"
)

func Translate(ctx context.Context, input string, targetLanguage string) (string, error) {
	ret, err := TranslateBatch(ctx, []string{input}, targetLanguage)
	if err != nil {
		return "", err
	}
	return ret[0], nil
}

func TranslateBatch(ctx context.Context, inputs []string, targetLanguage string) ([]string, error) {
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return nil, err
	}

	resp, err := client.Translate(ctx, inputs, lang, nil)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, errors.New("translate response is empty")
	}
	var ret []string
	for _, r := range resp {
		ret = append(ret, r.Text)
	}
	return ret, nil
}
