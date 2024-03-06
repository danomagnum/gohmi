package logix

import (
	"fmt"
	"strconv"
)

func convertTypeFromString(value string, target any) (any, error) {
	switch target.(type) {
	case int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return intValue, nil
	case int8:
		intValue, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return nil, err
		}
		return int8(intValue), nil
	case int16:
		intValue, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return nil, err
		}
		return int16(intValue), nil
	case int32:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(intValue), nil
	case int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return int64(intValue), nil
	case uint:
		uintValue, err := strconv.ParseUint(value, 10, strconv.IntSize)
		if err != nil {
			return nil, err
		}
		return uint(uintValue), nil
	case uint8:
		uintValue, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return nil, err
		}
		return uint8(uintValue), nil
	case uint16:
		uintValue, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return nil, err
		}
		return uint16(uintValue), nil
	case uint32:
		uintValue, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint32(uintValue), nil
	case uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return uint64(uintValue), nil
	case float32:
		floatValue, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		return float32(floatValue), nil
	case float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return floatValue, nil
	case string:
		return value, nil
	case bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		return boolValue, nil
	default:
		return nil, fmt.Errorf("unsupported target type")
	}
}
