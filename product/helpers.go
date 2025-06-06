package product

import "strconv"

func StringToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
