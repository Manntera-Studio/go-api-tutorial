package main

import (
	"bytes"
    "context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/api/idtoken"
)

func main()  {
	URL := "http://example.com" // アプリのURLに書き換える
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
	}
	audience := "IAP_CLIENT_ID.apps.googleusercontent.com" //設定したサービスアカウントのクライアントIDに書き換える
	respBody := new(bytes.Buffer)
	if err := makeIAPRequest(respBody, request, audience); err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
	}
	fmt.Println(respBody.String())
}

func makeIAPRequest(w io.Writer, request *http.Request, audience string) error {
	ctx := context.Background()

	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return fmt.Errorf("idtoken.NewClient: %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do: %v", err)
	}
	defer response.Body.Close()
	if _, err := io.Copy(w, response.Body); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	return nil
}
