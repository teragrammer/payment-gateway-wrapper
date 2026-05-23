package utils

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StringToInt converts a string to an integer and returns an error if the conversion fails.
func StringToInt(s string, def int) int {
	// Attempt to convert the string to an integer
	num, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return num
}

// StringToInt64 converts a string to an integer64 and returns an error if the conversion fails.
func StringToInt64(s string, def int64) int64 {
	// Attempt to convert the string to an integer
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return num
}

func IntToString(num int) string {
	return strconv.Itoa(num) // Use strconv.Itoa to convert int to string
}

func StringToObjectID(s string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.NilObjectID
	}
	return objectID
}
