package biz

import (
	"math"
	"regexp"

	"github.com/burstsms/mtmo-tp/backend/lib/number"
)

func ParseMobileCountry(mobile string, country string) (string, string, error) {
	return number.ParseMobileCountry(mobile, country)
}

func GetCountryFromPhone(mobile string) (string, error) {
	return number.GetCountryFromPhone(mobile)
}

func IsValidSMS(message string) (int, error) {
	// TODO unsub url length calculation
	// this would be knowing what char count to use for [unsubscribe] in the message

	// message length
	smsCount := CountSMSParts(message)
	if smsCount > 4 {
		return smsCount, ErrInvalidSMSTooManyParts
	}

	return smsCount, nil
}

var gsmRegex = regexp.MustCompile(`[|^€{}[\]~]`)

// isGSMRune tests input for the GSM 3.38 character set using default alphabet and extension table. This set was manually generated from the table:
// https://en.wikipedia.org/wiki/GSM_03.38#GSM_7-bit_default_alphabet_and_extension_table_of_3GPP_TS_23.038_.2F_GSM_03.38
// The cases are ordered to generate a hit early by detecting frequently used characters to not-so-frequently used.
func isGSMRune(r rune) bool {
	switch {
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	case '0' <= r && r <= '9':
		return true
	case r == ' ' || r == '\n' || r == '\r':
		return true
	case r == '!' || r == '"' || r == '#' || r == '@' || r == '_' || r == '%' || r == '&' || r == '\'' || r == '\\' || r == '(' || r == ')' || r == '*' || r == '+' || r == ',' || r == '-' || r == '.' || r == '/' || r == ':' || r == ';' || r == '<' || r == '=' || r == '>' || r == '?' || r == '$' || r == '^' || r == '{' || r == '}' || r == '[' || r == '~' || r == ']' || r == '|':
		return true
	case r == '£' || r == '¥' || r == 'è' || r == 'é' || r == 'ù' || r == 'ì' || r == 'ò' || r == 'Ç' || r == 'Ø' || r == 'ø' || r == 'Å' || r == 'å' || r == 'Δ' || r == 'Φ' || r == 'Γ' || r == 'Λ' || r == 'Ω' || r == 'Π' || r == 'Ψ' || r == 'Σ' || r == 'Θ' || r == 'Ξ' || r == 'Æ' || r == 'æ' || r == 'ß' || r == 'É' || r == '¤' || r == '¡' || r == 'Ä' || r == 'Ö' || r == 'Ñ' || r == 'Ü' || r == '§' || r == '¿' || r == 'ä' || r == 'ö' || r == 'ñ' || r == 'ü' || r == 'à' || r == '€':
		return true
	}

	return false
}

// IsGSMString is a helper wrapper for testing a whole string contains GSM Chars
func IsGSMString(s string) bool {
	for _, r := range s {
		if !isGSMRune(r) {
			return false
		}
	}

	return true
}

// CountSMSParts returns the sms count of a message body
func CountSMSParts(content string) int {
	var start, cutOff int = 161, 153 // GSM by default

	if !IsGSMString(content) { // If not GSM string, use unicode
		start = 71
		cutOff = 67
	}

	contentLength := contentLength(content)
	if contentLength >= start {
		return int(math.Ceil(float64(contentLength) / float64(cutOff)))
	}

	return 1
}

// contentLength returns length of a message body considering gsm escape chars
func contentLength(content string) int {
	var length = len([]rune(content))

	if IsGSMString(content) {
		escapeCharsFound := gsmRegex.FindAllString(content, -1)
		length += len(escapeCharsFound)
	}

	return length
}
