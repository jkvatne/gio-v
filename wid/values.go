package wid

import (
	"fmt"
	"math"
	"strconv"
)

// ValueToSTring converts a value to string.
// Accepts both pointers to values and values.
// If the value is numeric, it is converted to string.
func ValueToString(v interface{}, dp int) string {
	if v == nil {
		return "nil"
	} else if x, ok := v.(*int); ok {
		if *x == math.MinInt {
			return "---"
		} else {
			return fmt.Sprintf("%d", *x)
		}
	} else if x, ok := v.(int); ok {
		if x == math.MinInt {
			return "---"
		} else {
			return fmt.Sprintf("%d", x)
		}
	} else if x, ok := v.(*int32); ok {
		if *x == math.MinInt32 {
			return "---"
		} else {
			return fmt.Sprintf("%d", *x)
		}
	} else if x, ok := v.(int32); ok {
		if x == math.MinInt32 {
			return "---"
		} else {
			return fmt.Sprintf("%d", x)
		}
	} else if x, ok := v.(*int64); ok {
		if *x == math.MinInt64 {
			return "---"
		} else {
			return fmt.Sprintf("%d", *x)
		}
	} else if x, ok := v.(int64); ok {
		if x == math.MinInt {
			return "---"
		} else {
			return fmt.Sprintf("%d", x)
		}
	} else if x, ok := v.(*float32); ok {
		if *x == math.MaxFloat32 {
			return "---"
		} else {
			return fmt.Sprintf("%.*f", dp, *x)
		}
	} else if x, ok := v.(*float64); ok {
		if *x == math.MaxFloat64 {
			return "---"
		} else {
			return fmt.Sprintf("%.*f", dp, *x)
		}
	} else if x, ok := v.(*string); ok {
		return *x
	} else if x, ok := v.(string); ok {
		return x
	}
	return ""
}

// StringToValue will convert a string to a numeric value
// and store it at the address supplied. Will also accept string pointers
func StringToValue(p interface{}, current string) {
	if _, ok := p.(*int); ok {
		x, err := strconv.Atoi(current)
		if err == nil {
			*p.(*int) = x
		}
	} else if _, ok := p.(*int32); ok {
		x, err := strconv.Atoi(current)
		if err == nil {
			*p.(*int) = x
		}
	} else if _, ok := p.(*int64); ok {
		x, err := strconv.ParseInt(current, 10, 64)
		if err == nil {
			*p.(*int64) = x
		}
	} else if _, ok := p.(*float32); ok {
		f, err := strconv.ParseFloat(current, 32)
		if err == nil {
			*p.(*float32) = float32(f)
		}
	} else if _, ok := p.(*float64); ok {
		f, err := strconv.ParseFloat(current, 64)
		if err == nil {
			*p.(*float64) = f
		}
	} else if _, ok := p.(*string); ok {
		*p.(*string) = current
	} else {
		panic("Edit value should be pointer to value")
	}
}
