package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

/*func TestMain(m *testing.M) {
//ctx := context.Background()
/*ctx = nktest.WithAlwaysPullFromEnv(ctx, "PULL")
ctx = nktest.WithUnderCIFromEnv(ctx, "CI")
ctx = nktest.WithHostPortMap(ctx)
var opts []nktest.BuildConfigOption
if os.Getenv("CI") == "" {
	opts = append(opts, nktest.WithDefaultGoEnv(), nktest.WithDefaultGoVolumes())
}*/
/*nktest.Main(ctx, m,
		nktest.WithDir("."),
		nktest.WithBuildConfig(".", opts...),
	)
}
*/

type GetContentRequest struct {
	UserId string `json:"userId"`
	Body   string `json:"body"`
}

type ContentResponse struct {
	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
	Hash    string `json:"hash,omitempty"`
	Content string `json:"content,omitempty"`
}

type GetContentResponse struct {
	Body         string `json:"body"`
	ErrorMessage string `json:"error_message"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func processNakamaRequest(t *testing.T, bodyRaw string) *http.Response {
	ctx := context.Background()
	requestBody := GetContentRequest{UserId: "", Body: bodyRaw}
	requestBodyRaw, err := json.Marshal(requestBody)
	fmt.Println(string(requestBodyRaw))
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"http://127.0.0.1:7351/v2/console/api/endpoints/rpc/getcontent",
		bytes.NewReader(requestBodyRaw))

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TOKEN")))
	req.Header.Set("Content-Type", "application/json")

	cl := &http.Client{
		Transport: &http.Transport{},
	}
	res, err := cl.Do(req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	return res
}

func readGetContentResponse(t *testing.T, res *http.Response) GetContentResponse {
	bodyResponseRaw, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	fmt.Println(string(bodyResponseRaw))
	response := GetContentResponse{}
	errBody := json.Unmarshal(bodyResponseRaw, &response)
	if errBody != nil {
		log.Fatalln(errBody)
		t.Fatalf("expected no error, got: %v", err)
	}
	return response
}

func readContentResponse(t *testing.T, responseRaw string) ContentResponse {
	contentResponse := ContentResponse{}
	errContent := json.Unmarshal([]byte(responseRaw), &contentResponse)
	if errContent != nil {
		log.Fatalln(errContent)
		t.Fatalf("expected no error, got: %v", errContent)
	}
	return contentResponse
}

func TestGetContentValidRequest(t *testing.T) {
	bodyRaw := "{\"version\": \"1.0.0\", \"hash\": \"d6e4677dc8987b7b140ad75384bb7a49adea29c7bbb5c1191e932420fe8a067e\", \"type\": \"core\"}"
	res := processNakamaRequest(t, bodyRaw)
	// check response
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got: %d", http.StatusOK, res.StatusCode)
	}

	response := readGetContentResponse(t, res)
	contentResponse := readContentResponse(t, response.Body)

	assert.Equal(t, contentResponse.Type, "core")
	assert.Equal(t, contentResponse.Version, "1.0.0")
	assert.Equal(t, contentResponse.Hash, "d6e4677dc8987b7b140ad75384bb7a49adea29c7bbb5c1191e932420fe8a067e")
	assert.Equal(t, contentResponse.Content, "{\"jsonFileWith\": \"Content\"}")
	assert.Empty(t, response.ErrorMessage)
}

func TestGetContentDefaultRequest(t *testing.T) {
	bodyRaw := "{}"
	res := processNakamaRequest(t, bodyRaw)
	// check response
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got: %d", http.StatusOK, res.StatusCode)
	}

	response := readGetContentResponse(t, res)
	contentResponse := readContentResponse(t, response.Body)

	assert.Equal(t, "core", contentResponse.Type)
	assert.Equal(t, "1.0.0", contentResponse.Version)
	assert.Equal(t, "d6e4677dc8987b7b140ad75384bb7a49adea29c7bbb5c1191e932420fe8a067e", contentResponse.Hash)
	assert.Equal(t, "{\"jsonFileWith\": \"Content\"}", contentResponse.Content)
	assert.Empty(t, response.ErrorMessage)
}

func TestGetContentWrongHash(t *testing.T) {
	bodyRaw := "{\"version\": \"1.0.0\", \"hash\": \"123\", \"type\": \"core\"}"
	res := processNakamaRequest(t, bodyRaw)
	// check response
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got: %d", http.StatusOK, res.StatusCode)
	}

	response := readGetContentResponse(t, res)
	contentResponse := readContentResponse(t, response.Body)

	assert.Equal(t, "core", contentResponse.Type)
	assert.Equal(t, "1.0.0", contentResponse.Version)
	assert.Equal(t, "", contentResponse.Hash)
	assert.Equal(t, "", contentResponse.Content)
	assert.Empty(t, response.ErrorMessage)
}

func TestGetContentWrongType(t *testing.T) {
	bodyRaw := "{\"version\": \"1.0.0\", \"type\": \"rore\"}"
	res := processNakamaRequest(t, bodyRaw)
	// check response
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got: %d", http.StatusOK, res.StatusCode)
	}

	response := readGetContentResponse(t, res)

	assert.NotEmpty(t, response.ErrorMessage)
}
