# Ego SMS sdk for Go

Example:
```go
sdk, err := v1.Authenticate("username", "password")

success, err := sdk.SendSMS("+256772123456", "Test message to single number")

numbers := []string{"+256772123456", "0772123457"}
success, err := sdk.SendSMSWithSenderId(numbers, "Test message to many numbers", "MySenderID")

balance, err := sdk.GetBalance()
```