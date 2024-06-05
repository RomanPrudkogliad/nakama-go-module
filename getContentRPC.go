package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/heroiclabs/nakama-common/runtime"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type ZeptoRequest struct {
	Type    string `json:"type,omitempty" validate:"oneof=core ios android"`
	Version string `json:"version,omitempty" validate:"oneof=1.0.0 0.1.1 0.0.1"`
	Hash    string `json:"hash,omitempty"`
}

type ZeptoResponse struct {
	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
}

type ContentStorageObject struct {
	Request   ZeptoRequest `json:"request"`
	FilePath  string       `json:"filePath"`
	CreatedAt int64        `json:"created_at"`
}

func GetContent(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("GetContent request: %s", payload)
	var request = &ZeptoRequest{Type: "core", Version: "1.0.0", Hash: "d6e4677dc8987b7b140ad75384bb7a49adea29c7bbb5c1191e932420fe8a067e"}
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", fmt.Errorf("unable to unmarshal request: %w, payload %s", err, payload)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validationErr := validate.Struct(request)
	if validationErr != nil {
		fmt.Println(validationErr)
		return "", fmt.Errorf("unable to process request: %w, payload %s", validationErr, payload)
	}

	filePath := "data/" + request.Type + "/" + request.Version + ".json"
	file, fileErr := os.ReadFile(filePath)
	if fileErr != nil {
		return "", fmt.Errorf("unable to read file: %w", fileErr)
	}

	now := time.Now().Unix()
	storageObject := &ContentStorageObject{
		Request:   *request,
		FilePath:  filePath,
		CreatedAt: now,
	}

	var storageRaw, errStorageErr = json.Marshal(storageObject)
	if errStorageErr != nil {
		logger.Error("Failed to marshal response.", errStorageErr)
		return "", fmt.Errorf("unable to Marshal storageObject: %w", errStorageErr)
	}

	writes := []*runtime.StorageWrite{{
		Collection:      "content_collection",
		Key:             fmt.Sprintf("request_%d", now),
		Value:           string(storageRaw),
		UserID:          "",
		PermissionRead:  2,
		PermissionWrite: 1,
	}}
	_, storageErr := nk.StorageWrite(ctx, writes)
	if storageErr != nil {
		logger.WithField("err", storageErr).Error("Storage write error.")
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, strings.NewReader(string(file))); err != nil {
		log.Fatal(err)
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))
	responseContent := ""
	if strings.Compare(fileHash, request.Hash) == 0 {
		responseContent = string(file)
	} else {
		fileHash = ""
	}

	var zeptoResponse = &ZeptoResponse{
		Type:    request.Type,
		Version: request.Version,
		Hash:    fileHash,
		Content: responseContent,
	}

	var responseRaw, err = json.Marshal(zeptoResponse)
	if err != nil {
		logger.Error("Failed to marshal response.", err)
		return "", fmt.Errorf("unable to Marshal response: %w", err)
	}

	/*
		userId, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
		recordList, _, readErr := nk.StorageList(ctx, userId, "", "content_collection", 10, "")

		if readErr != nil {
			logger.WithField("err", readErr).Error("Storage read error.")
		} else {
			logger.Info("Storage read records %s", len(recordList))
			for _, record := range recordList {
				logger.Info("value: %s", record.Value)
			}
		}*/

	return string(responseRaw), nil
}
