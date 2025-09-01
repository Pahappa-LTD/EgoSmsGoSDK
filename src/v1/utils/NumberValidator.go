package utils

import (
    "regexp"
    "strings"
)

var regex = regexp.MustCompile(`^\+?(0|\d{3})\d{9}$`)

func ValidateNumbers(numbers []string) []string {
    if len(numbers) == 0 {
        return []string{}
    }

    cleansed := make(map[string]bool)
    for _, number := range numbers {
        if strings.TrimSpace(number) == "" {
            continue
        }

        cleanedNumber := strings.ReplaceAll(strings.ReplaceAll(number, "-", ""), " ", "")
        if regex.MatchString(cleanedNumber) {
            if strings.HasPrefix(cleanedNumber, "0") {
                cleanedNumber = "256" + cleanedNumber[1:]
            } else if strings.HasPrefix(cleanedNumber, "+") {
                cleanedNumber = cleanedNumber[1:]
            }
            cleansed[cleanedNumber] = true
        }
    }

    result := make([]string, 0, len(cleansed))
    for k := range cleansed {
        result = append(result, k)
    }
    return result
}
