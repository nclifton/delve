package biz

import (
	"fmt"
	"math"
	"regexp"
	"unicode/utf16"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/lib/number"
	"github.com/burstsms/mtmo-tp/backend/lib/stringutil"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

type SMSOptions struct {
	TrackLink        bool
	TrackLinkDomain  string
	OptOutLinkDomain string
	MaxParts         int
}

func ParseMobileCountry(mobile string, country string) (string, string, error) {
	return number.ParseMobileCountry(mobile, country)
}

func GetCountryFromPhone(mobile string) (string, error) {
	return number.GetCountryFromPhone(mobile)
}

func IsValidSMS(message string, opt SMSOptions) (int, error) {
	// TODO unsub url length calculation
	// this would be knowing what char count to use for [unsubscribe] in the message

	// message length
	smsCount := CountSMSParts(message, opt)
	if smsCount > opt.MaxParts {
		return smsCount, errorlib.ErrInvalidSMSTooManyParts
	}

	return smsCount, nil
}

func IsValidSender(sender *senderpb.Sender, address, country string) error {
	if sender == nil {
		return errorlib.ErrInvalidSenderNotFound
	}
	if sender.Address != address {
		return errorlib.ErrInvalidSenderAddress
	}
	if !stringutil.Includes(sender.Channels, "sms") {
		return errorlib.ErrInvalidSenderChannel
	}
	return nil
}

var gsmRegex = regexp.MustCompile(`[|^€{}[\]~]`)

// isGSMRune tests input for the GSM 3.38 character set using default alphabet and extension table. This set was manually generated from the table:
// https://en.wikipedia.org/wiki/GSM_03.38#GSM_7-bit_default_alphabet_and_extension_table_of_3GPP_TS_23.038_.2F_GSM_03.38
// The cases are ordered to generate a hit early by detecting frequently used characters to not-so-frequently used.
func isGSMRune(r rune) (isGSM bool, isGMExtended bool) {
	switch {
	case 'a' <= r && r <= 'z':
		return true, false
	case 'A' <= r && r <= 'Z':
		return true, false
	case '0' <= r && r <= '9':
		return true, false
	case r == ' ' || r == '\n' || r == '\r':
		return true, false
	case r == '!' || r == '"' || r == '#' || r == '@' || r == '_' || r == '%' || r == '&' || r == '\'' || r == '(' || r == ')' || r == '*' || r == '+' || r == ',' || r == '-' || r == '.' || r == '/' || r == ':' || r == ';' || r == '<' || r == '=' || r == '>' || r == '?' || r == '$':
		return true, false
	case r == '£' || r == '¥' || r == 'è' || r == 'é' || r == 'ù' || r == 'ì' || r == 'ò' || r == 'Ç' || r == 'Ø' || r == 'ø' || r == 'Å' || r == 'å' || r == 'Δ' || r == 'Φ' || r == 'Γ' || r == 'Λ' || r == 'Ω' || r == 'Π' || r == 'Ψ' || r == 'Σ' || r == 'Θ' || r == 'Ξ' || r == 'Æ' || r == 'æ' || r == 'ß' || r == 'É' || r == '¤' || r == '¡' || r == 'Ä' || r == 'Ö' || r == 'Ñ' || r == 'Ü' || r == '§' || r == '¿' || r == 'ä' || r == 'ö' || r == 'ñ' || r == 'ü' || r == 'à':
		return true, false
	case r == '\f' || r == '[' || r == '\\' || r == ']' || r == '^' || r == '{' || r == '}' || r == '|' || r == '~' || r == '€':
		return true, true // esc chars (dbl size)
	}

	return false, false
}

// IsGSMString is a helper wrapper for testing a whole string contains GSM Chars
func IsGSMString(s string) bool {
	for _, r := range s {
		gsm, _ := isGSMRune(r)
		if !gsm {
			return false
		}
	}

	return true
}

type Split struct {
	Charset        string // GSM or Unicode
	Length         int    // Total length of message
	Bytes          int    // Total bytes of message
	CountParts     int    // Number of SMS in message
	Parts          []Sms  // SMS parts
	RemainingChars int    // Remaining char in current SMS
}

type Sms struct {
	Content string
	Bytes   int
	Length  int
}

const (
	smsBytes      int = 140
	singleGsm7    int = 160
	singleUnicode int = 70
)

var (
	multiGSM7    int
	multiUnicode int
)

func (m *Split) appendSMS(sms string, bytes int, length int) {
	if bytes > 0 {
		m.Parts = append(m.Parts, Sms{sms, bytes, length})
		m.Length += length
		m.Bytes += bytes
		m.CountParts = len(m.Parts)

		if m.Charset == "GSM" {
			m.RemainingChars = singleGsm7 - m.Bytes
			if len(m.Parts) > 1 {
				m.RemainingChars = multiGSM7 - m.Parts[len(m.Parts)-1].Bytes
			}
		} else {
			m.RemainingChars = singleUnicode - (m.Bytes / 2)
			if len(m.Parts) > 1 {
				m.RemainingChars = multiUnicode - (m.Parts[len(m.Parts)-1].Bytes / 2)
			}
		}
	}
}

func SplitSMSParts(content string) (*Split, error) {

	isGSM := IsGSMString(content)
	bytes := 0
	length := 0
	curSMS := ""

	split := &Split{Charset: "GSM"}

	multiGSM7 = ((smsBytes * 8) - (6 * 8)) / 7
	multiUnicode = ((smsBytes * 8) - (6 * 8)) / 16

	for _, char := range content {
		if isGSM {
			_, gsmExt := isGSMRune(char)
			if gsmExt {
				if bytes == multiGSM7-1 {
					split.appendSMS(curSMS, bytes, length)
					bytes = 0
					length = 0
					curSMS = ""
				}
				// Add escape code
				bytes++
			}
			bytes++
			length++
		} else {
			split.Charset = "Unicode"
			if isHighSurrogate(char) {
				if bytes == (multiUnicode-1)*2 {
					split.appendSMS(curSMS, bytes, length)
					bytes = 0
					length = 0
					curSMS = ""
				}
				bytes += 2
			}
			bytes += 2
			length++
		}

		curSMS += string(char)

		if (isGSM && bytes == multiGSM7) || (!isGSM && bytes == (multiUnicode*2)) {
			split.appendSMS(curSMS, bytes, length)
			bytes = 0
			length = 0
			curSMS = ""
		}

	}

	split.appendSMS(curSMS, bytes, length)

	if (isGSM && len(split.Parts) > 1 && split.Bytes <= singleGsm7) || (!isGSM && len(split.Parts) > 1 && split.Bytes <= (singleUnicode*2)) {
		split.Parts[0].Content += split.Parts[1].Content
		split.Parts[0].Bytes += split.Parts[1].Bytes
		split.Parts[0].Length += split.Parts[1].Length
		split.Parts = split.Parts[:len(split.Parts)-1]
		split.CountParts = 1
		if isGSM {
			split.RemainingChars = singleGsm7 - split.Bytes
		} else {
			split.RemainingChars = (singleUnicode * 2) - split.Bytes
		}
	}

	return split, nil
}

func isHighSurrogate(r rune) bool {
	r1, _ := utf16.EncodeRune(r)
	return r1 >= 0xD800 && r1 <= 0xDBFF
}

// CountSMSParts returns the sms count of a message body
func CountSMSParts(content string, opt SMSOptions) int {
	var start, cutOff int = 161, 153 // GSM by default

	if !IsGSMString(content) { // If not GSM string, use unicode
		start = 71
		cutOff = 67
	}

	contentLength := contentLength(content, opt)
	if contentLength >= start {
		return int(math.Ceil(float64(contentLength) / float64(cutOff)))
	}

	return 1
}

var urlRegex = regexp.MustCompile(`(http|https)\://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,6}(/\S*)?`)
var optoutRegex = regexp.MustCompile(`\[opt-out-link\]`)

// contentLength returns length of a message body considering gsm escape chars
func contentLength(content string, opt SMSOptions) int {
	if opt.TrackLink {
		// Replace all urls in content with a template matching the length
		// of a shortlink
		content = urlRegex.ReplaceAllString(content, fmt.Sprintf(`%s/12345678`, opt.TrackLinkDomain))
	}
	// Check for opt out link tag
	content = optoutRegex.ReplaceAllString(content, fmt.Sprintf(`%s/12345678`, opt.OptOutLinkDomain))

	var length = len([]rune(content))

	if IsGSMString(content) {
		escapeCharsFound := gsmRegex.FindAllString(content, -1)
		length += len(escapeCharsFound)
	}

	return length
}
