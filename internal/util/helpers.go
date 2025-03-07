package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/shashank-sharma/backend/internal/logger"
)

func GenerateRandomId() string {
	return security.RandomStringWithAlphabet(core.DefaultIdLength, core.DefaultIdAlphabet)
}

func GetUserId(tokenString string) (string, error) {
	// Split the token into header, payload, and signature
	parts := strings.Split(tokenString, ".")

	// Decode the payload (no signature verification)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		logger.LogError("Error decoding payload:", err)
		return "", err
	}

	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		logger.LogError("Error unmarshalling payload:", err)
		return "", err
	}

	return claims["id"].(string), nil
}

// AnyToString converts any type to a string representation
func AnyToString(value interface{}) string {
    if value == nil {
        return ""
    }

    // Use reflection to handle special types
    v := reflect.ValueOf(value)
    
    // Handle pointer types by dereferencing
    if v.Kind() == reflect.Ptr {
        if v.IsNil() {
            return ""
        }
        v = v.Elem()
        value = v.Interface()
    }

    switch val := value.(type) {
    case string:
        return val
    case bool:
        return strconv.FormatBool(val)
    case int:
        return strconv.Itoa(val)
    case int64:
        return strconv.FormatInt(val, 10)
    case int32:
        return strconv.FormatInt(int64(val), 10)
    case int16:
        return strconv.FormatInt(int64(val), 10)
    case int8:
        return strconv.FormatInt(int64(val), 10)
    case uint:
        return strconv.FormatUint(uint64(val), 10)
    case uint64:
        return strconv.FormatUint(val, 10)
    case uint32:
        return strconv.FormatUint(uint64(val), 10)
    case uint16:
        return strconv.FormatUint(uint64(val), 10)
    case uint8:
        return strconv.FormatUint(uint64(val), 10)
    case float64:
        return strconv.FormatFloat(val, 'f', -1, 64)
    case float32:
        return strconv.FormatFloat(float64(val), 'f', -1, 32)
    case time.Time:
        return val.Format(time.RFC3339)
    case []byte:
        return string(val)
    default:
        // For complex types, use fmt.Sprintf
        return fmt.Sprintf("%v", val)
    }
}
