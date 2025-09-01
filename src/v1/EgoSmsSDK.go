package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Pahappa-LTD/EgoSmsGoSDK/src/v1/models"
	"github.com/Pahappa-LTD/EgoSmsGoSDK/src/v1/utils"
)

var API_URL = "https://api.egosms.co/api/v1/json/"

type EgoSmsSDK struct {
	apiKey          string
	username        string
	password        string
	senderId        string
	isAuthenticated bool
}

func (sdk *EgoSmsSDK) GetApiKey() string {
	return sdk.apiKey
}

func (sdk *EgoSmsSDK) GetUsername() string {
	return sdk.username
}

func (sdk *EgoSmsSDK) GetPassword() string {
	return sdk.password
}

func (sdk *EgoSmsSDK) SetAuthenticated(isAuthenticated bool) {
	sdk.isAuthenticated = isAuthenticated
}

func (sdk *EgoSmsSDK) GetApiURL() string {
	return API_URL
}

func (sdk *EgoSmsSDK) SendSMS(numbers interface{}, message string) (bool, error) {
	return sdk.SendSMSFull(numbers, message, sdk.senderId, models.HIGHEST)
}

func (sdk *EgoSmsSDK) SendSMSWithSenderId(numbers interface{}, message string, senderId string) (bool, error) {
	return sdk.SendSMSFull(numbers, message, senderId, models.HIGHEST)
}

func (sdk *EgoSmsSDK) SendSMSWithPriority(numbers interface{}, message string, priority models.MessagePriority) (bool, error) {
	return sdk.SendSMSFull(numbers, message, sdk.senderId, priority)
}

func UseSandBox() {
	API_URL = "http://sandbox.egosms.co/api/v1/json/"
}

func UseLiveServer() {
	API_URL = "https://www.egosms.co/api/v1/json/"
}

func AuthenticateWithApiKey(apiKey string) (*EgoSmsSDK, error) {
	return nil, errors.New("API Key authentication is not supported in this version. Please use username and password authentication")
}

func Authenticate(username string, password string) (*EgoSmsSDK, error) {
	sdk := &EgoSmsSDK{
		username: username,
		password: password,
		senderId: "EgoSms",
	}

	isValid, err := utils.ValidateCredentials(sdk)
	if err != nil {
		return nil, err
	}

	sdk.isAuthenticated = isValid
	return sdk, nil
}

func (sdk *EgoSmsSDK) WithSenderId(senderId string) *EgoSmsSDK {
	sdk.senderId = senderId
	return sdk
}

func (sdk *EgoSmsSDK) SendSMSFull(numbers interface{}, message string, senderId string, priority models.MessagePriority) (bool, error) {
	if !sdk.isAuthenticated {
		return false, errors.New("SDK is not authenticated. Please authenticate before performing actions")
	}

	var numberSlice []string
	switch v := numbers.(type) {
	case string:
		numberSlice = []string{v}
	case []string:
		numberSlice = v
	default:
		return false, errors.New("numbers must be a string or a slice of strings")
	}

	if len(numberSlice) == 0 {
		return false, errors.New("numbers list cannot be null or empty")
	}

	if message == "" {
		return false, errors.New("message cannot be null or empty")
	}

	if len(message) == 1 {
		return false, errors.New("message cannot be a single character")
	}

	if strings.TrimSpace(senderId) == "" {
		senderId = sdk.senderId
	}

	if len(senderId) > 11 {
		fmt.Println("Warning: Sender ID length exceeds 11 characters. Some networks may truncate or reject messages.")
	}

	validatedNumbers := utils.ValidateNumbers(numberSlice)

	if len(validatedNumbers) == 0 {
		return false, errors.New("no valid phone numbers provided. Please check inputs")
	}

	var messageModels []models.MessageModel
	for _, number := range validatedNumbers {
		messageModels = append(messageModels, models.MessageModel{
			Number:   number,
			Message:  message,
			SenderId: senderId,
			Priority: priority,
		})
	}

	apiRequest := models.ApiRequest{
		Method:      "SendSms",
		Userdata:    models.UserData{Username: sdk.username, Password: sdk.password},
		MessageData: messageModels,
	}

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(API_URL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var apiResponse models.ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return false, err
	}

	if apiResponse.Status == models.OK {
		fmt.Println("SMS sent successfully.")
		return true, nil
	} else {
		return false, errors.New(apiResponse.Message)
	}
}

func (sdk *EgoSmsSDK) GetBalance() (string, error) {
	if !sdk.isAuthenticated {
		return "", errors.New("SDK is not authenticated. Please authenticate before performing actions")
	}

	apiRequest := models.ApiRequest{
		Method:   "Balance",
		Userdata: models.UserData{Username: sdk.username, Password: sdk.password},
	}

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(API_URL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResponse models.ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", err
	}

	return apiResponse.Balance, nil
}
