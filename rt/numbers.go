package rt

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Integer formats an integer to a string based on a user-specified format.
func Integer(n int, format string) string {
	return renderFloat(float64(n), format)
}

// Float formats a float64 to a string based on a user-specified format.
func Float(n float64, format string) string {
	return renderFloat(n, format)
}

// Decimal formats a big.Rat to a string.
func Decimal(n big.Rat) string {
	return n.FloatString(2)
}

// IntegerTrimZero formats an integer to a string and removes the trailing ".00", if any.
func IntegerTrimZero(n int, format string) string {
	return TrimZero(Integer(n, format))
}

// FloatTrimZero formats a float64 to a string and removes the trailing ".00", if any.
func FloatTrimZero(n float64, format string) string {
	return TrimZero(Float(n, format))
}

// Decimal formats a big.Rat to a string and removes trailing "0"s, if any.
func DecimalTrimZero(n big.Rat) string {
	return strings.TrimRight(strings.TrimRight(Decimal(n), "0"), ".")
}

// Ordinal formats an integer to ordinal format string (1st, 2nd, 3rd, etc).
func Ordinal(n int) string {
	return ordinal(n)
}

// ------------------------------------------------------------------
//
//
// internal - ordinal
//
//
// ------------------------------------------------------------------

// ordinal formats an integer to ordinal format string (1st, 2nd, 3rd, etc).
//
// All credit goes to the original author.
// Ref: https://github.com/dustin/go-humanize/blob/master/ordinals.go
func ordinal(n int) string {
	suffix := "th"
	switch n % 10 {
	case 1:
		if n%100 != 11 {
			suffix = "st"
		}
	case 2:
		if n%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if n%100 != 13 {
			suffix = "rd"
		}
	}
	return strconv.Itoa(n) + suffix
}

// ------------------------------------------------------------------
//
//
// internal - renderFloat
//
//
// ------------------------------------------------------------------

// renderFloat renders a float64 value to a string based on the following user-specific criteria:
//   - thousands separator
//   - decimal separator
//   - decimal precision
//
// Examples of format strings, given n = 12345.6789:
//   - "#,###.##" => "12,345.67"
//   - "#,###." => "12,345"
//   - "#,###" => "12345,678"
//   - "#\u202F###,##" => "12 345,67"
//   - "#.###,###### => 12.345,678900
//   - "" (aka default format) => 12,345.67
//
// The highest precision allowed is 9 digits after the decimal symbol.
//
// All credit goes to the original author.
// Ref: https://gist.github.com/gorhill/5285193
func renderFloat(n float64, format string) string {
	// Special cases:
	//   NaN = "NaN"
	//   +Inf = "+Infinity"
	//   -Inf = "-Infinity"
	if math.IsNaN(n) {
		return "NaN"
	}
	if n > math.MaxFloat64 {
		return "Infinity"
	}
	if n < -math.MaxFloat64 {
		return "-Infinity"
	}

	// default format
	precision := 2
	decimalStr := "."
	thousandStr := ","
	positiveStr := ""
	negativeStr := "-"

	if len(format) > 0 {
		// If there is an explicit format directive,
		// then default values are these:
		precision = 9
		thousandStr = ""

		// collect indices of meaningful formatting directives
		formatDirectiveChars := []rune(format)
		formatDirectiveIndices := make([]int, 0)
		for i, char := range formatDirectiveChars {
			if char != '#' && char != '0' {
				formatDirectiveIndices = append(formatDirectiveIndices, i)
			}
		}

		if len(formatDirectiveIndices) > 0 {
			// Directive at index 0:
			//   Must be a '+'
			//   Raise an error if not the case
			// index: 0123456789
			//        +0.000,000
			//        +000,000.0
			//        +0000.00
			//        +0000
			if formatDirectiveIndices[0] == 0 {
				if formatDirectiveChars[formatDirectiveIndices[0]] != '+' {
					panic("RenderFloat(): invalid positive sign directive")
				}
				positiveStr = "+"
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			// Two directives:
			//   First is thousands separator
			//   Raise an error if not followed by 3-digit
			// 0123456789
			// 0.000,000
			// 000,000.00
			if len(formatDirectiveIndices) == 2 {
				if (formatDirectiveIndices[1] - formatDirectiveIndices[0]) != 4 {
					panic("RenderFloat(): thousands separator directive must be followed by 3 digit-specifiers")
				}
				thousandStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				formatDirectiveIndices = formatDirectiveIndices[1:]
			}

			// One directive:
			//   Directive is decimal separator
			//   The number of digit-specifier following the separator indicates wanted precision
			// 0123456789
			// 0.00
			// 000,0000
			if len(formatDirectiveIndices) == 1 {
				decimalStr = string(formatDirectiveChars[formatDirectiveIndices[0]])
				precision = len(formatDirectiveChars) - formatDirectiveIndices[0] - 1
			}
		}
	}

	// generate sign part
	var signStr string
	if n >= 0.000000001 {
		signStr = positiveStr
	} else if n <= -0.000000001 {
		signStr = negativeStr
		n = -n
	} else {
		signStr = ""
		n = 0.0
	}

	// split number into integer and fractional parts
	intf, fracf := math.Modf(n + renderFloatPrecisionRounders[precision])

	// generate integer part string
	intStr := strconv.Itoa(int(intf))

	// add thousand separator if required
	if len(thousandStr) > 0 {
		for i := len(intStr); i > 3; {
			i -= 3
			intStr = intStr[:i] + thousandStr + intStr[i:]
		}
	}

	// no fractional part, we can leave now
	if precision == 0 {
		return signStr + intStr
	}

	// generate fractional part
	fracStr := strconv.Itoa(int(fracf * renderFloatPrecisionMultipliers[precision]))
	// may need padding
	if len(fracStr) < precision {
		fracStr = "000000000000000"[:precision-len(fracStr)] + fracStr
	}

	return signStr + intStr + decimalStr + fracStr
}

var renderFloatPrecisionMultipliers = [10]float64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
}

var renderFloatPrecisionRounders = [10]float64{
	0.5,
	0.05,
	0.005,
	0.0005,
	0.00005,
	0.000005,
	0.0000005,
	0.00000005,
	0.000000005,
	0.0000000005,
}
