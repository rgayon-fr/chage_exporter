package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Basic chage structure
type ChageStruct struct {
	User               string
	LastPasswordChange time.Time
	PasswordExpireDate time.Time
}

// Get duration until password expires
func (c *ChageStruct) UntilExpiryDate() time.Duration {
	// Default time
	defaultTime := time.Time{}

	// Return a duration of 0 if the password never expires.
	// NOTE: Changed from -1 to remove useless complexity.
	if c.PasswordExpireDate == defaultTime {
		return time.Duration(0)
	}

	return time.Until(c.PasswordExpireDate)
}

// Get age of password change
func (c *ChageStruct) PasswordAge() time.Duration {
	return time.Since(c.LastPasswordChange)
}

// Get the date from a chage field, return default time.Time struct if field is "never"
func GetDateFromChageField(i string) (time.Time, error) {
	// Don't bother with time parsing if the string is "never"
	if i == "never" {
		return time.Time{}, nil
	}

	// Parse the input.
	timeParse, timeError := time.Parse("2006-01-02", i)
	if timeError != nil {
		return time.Time{}, fmt.Errorf("error parsing time from chage field: %s", timeError)
	}

	return timeParse, nil
}

// Gets result of the chage command for a particular user.
func GetChage(u string) (ChageStruct, error) {
	// Init the return struct.
	var chageOutputStruct ChageStruct

	// Set user right out of the bat
	chageOutputStruct.User = u

	chageCommand := exec.Command("chage", "-l", "-i", u)
	chageOutput, chageError := chageCommand.CombinedOutput()
	if chageError != nil {
		return ChageStruct{}, fmt.Errorf("error running the chage command: %s", chageError)
	}

	// Convert the output into a string
	chageOutputString := string(chageOutput[:])

	// Split that output on new lines
	chageOutputLines := strings.Split(chageOutputString, "\n")

	// Compile regex strings
	// These are really simple and probably should be made into something a little more optimized.
	// For now however, they work.

	// Get last password change line
	lastPwChangeRegex, lastPwChangeError := regexp.Compile(`^(Last password change)`)
	if lastPwChangeError != nil {
		return ChageStruct{}, fmt.Errorf("error compiling lastPwChangeRegex: %s", lastPwChangeError)
	}

	// Get password expiry date line
	pwExpiresRegex, pwExpiresError := regexp.Compile(`^(Password expires)`)
	if pwExpiresError != nil {
		return ChageStruct{}, fmt.Errorf("error compiling pwExpiresRegex: %s", pwExpiresError)
	}

	// Simple date recognition regex.
	// TODO: Find better way to recognize dates in lines.
	simpleDateRegex, simpleDateError := regexp.Compile(`([0-9]{4}-[0-9]{2}-[0-9]{2}|never)`)
	if simpleDateError != nil {
		return ChageStruct{}, fmt.Errorf("error compiling simpleDateRegex: %s", simpleDateError)
	}

	// Scan through the output of chage and using regex get the information for the following lines:
	for _, line := range chageOutputLines {
		lineBytes := []byte(line)

		if lastPwChangeRegex.Match(lineBytes) {
			// We've found the last time the password was changed.
			// Get the date in the line.
			lineDate := simpleDateRegex.Find(lineBytes)
			lineDateStruct, lineDateError := GetDateFromChageField(string(lineDate[:]))
			if lineDateError != nil {
				return ChageStruct{}, fmt.Errorf("error parsing chage command: %s", lineDateError)
			}
			chageOutputStruct.LastPasswordChange = lineDateStruct
		}
		if pwExpiresRegex.Match(lineBytes) {
			// We've found the password expiry date
			// Get the date in the line
			lineDate := simpleDateRegex.Find(lineBytes)
			lineDateStruct, lineDateError := GetDateFromChageField(string(lineDate[:]))
			if lineDateError != nil {
				return ChageStruct{}, fmt.Errorf("error parsing chage command: %s", lineDateError)
			}
			chageOutputStruct.PasswordExpireDate = lineDateStruct
		}
	}

	// Last password change (Age calculation)
	// Password expires (Makes things *real* easy.)
	return chageOutputStruct, nil
}
