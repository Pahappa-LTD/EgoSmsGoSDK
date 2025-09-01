package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Pahappa-LTD/EgoSmsGoSDK/src/v1/models"
)

type EgoSmsSDKInterface interface {
	GetApiKey() string
	GetUsername() string
	GetPassword() string
	SetAuthenticated(bool)
	GetApiURL() string
}

func ValidateCredentials(sdk EgoSmsSDKInterface) (bool, error) {
	if sdk == nil {
		return false, errors.New("EgoSmsSDK instance cannot be null")
	}

	isApiKey := true
	if sdk.GetApiKey() == "" {
		if sdk.GetPassword() == "" || sdk.GetUsername() == "" {
			return false, errors.New("either API Key or Username and Password must be provided")
		} else {
			isApiKey = false
		}
	}

	if !isValidCredential(sdk, isApiKey) {
		fmt.Println("                                                      _                    ")
		fmt.Println("  /\\     _|_ |_   _  ._ _|_ o  _  _. _|_ o  _  ._    |_ _. o |  _   _| | |")
		fmt.Println(" /--\\ |_| |_ | | (/_ | | |_ | (_ (_|  |_ | (_) | |   | (_| | | (/_ (_| o o ")
		fmt.Println("                                                                           ")
		fmt.Println()
		return false, errors.New("credentials validation failed")
	}

	fmt.Println("Validated using an api key")
	fmt.Println()
	sdk.SetAuthenticated(true)
	return true, nil
}

func isValidCredential(sdk EgoSmsSDKInterface, isApiKey bool) bool {
	apiRequest := models.ApiRequest{}
	apiRequest.Method = "Balance"
	apiRequest.Userdata = models.UserData{Username: sdk.GetUsername(), Password: sdk.GetPassword()}

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
		return false
	}

	resp, err := http.Post(sdk.GetApiURL(), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return false
	}

	var apiResponse models.ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return false
	}

	if apiResponse.Status == models.OK {
		fmt.Println("Credentials validated successfully.")
		return true
	} else {
		fmt.Printf("Error validating credentials: %s\n", apiResponse.Message)
		return false
	}
}
