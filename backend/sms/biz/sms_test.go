package biz_test

import (
	"testing"

	"github.com/burstsms/mtmo-tp/backend/sms/biz"
)

const (
	phoneNumber         = "61404123456"
	invalidNumber       = "invalidnumber"
	gsmString           = "[Destiny guides our fortunes more favorably than we could have expected. Look there, Sancho Panza, my friend, and see those thirty or so wild giants, with whom I intend to do battle and kill each and all of them, so with their stolen booty we can begin to enrich ourselves. This is nobel, righteous warfare, for it is wonderfully useful to God to have such an evil race wiped from the face of the earth.' 'What giants?' Asked Sancho Panza. 'The ones you can see over there,' answered his master, 'with the huge arms, some of which are very nearly two leagues long.' 'Now look, your grace,' said Sancho, 'what you see over there aren't giants, but windmills, and what seems to be arms are just their sails, that go around in the wind and turn the millstone.' 'Obviously,' replied Don Quijote, 'you don't know much about adventures.]"
	gsmStringWithLinks  = "341 http://www.moose.com/ 125689xkftXCKijmqlnD4Pm3xvZTLrHkyQa2mSqKO9D3bzlLLUKYfKsM4vcGidc3cTZC3iTPKAkwSaiPbwTik9hnCKLBJNboj9opYLW0NpQJhr1qIuW0cKAmzjPrsIRzec4CJw0"
	gsmStringWithOptOut = "341 [opt-out-link] 4444125689xkftXCKijmqlnD4Pm3xvZTLrHkyQa2mSqKO9D3bzlLLUKYfKsM4vcGidc3cTZC3iTPKAkwSaiPbwTik9hnCKLBJNboj9opYLW0NpQJhr1qIuW0cKAmzjPrsIRzec4CJw0"
	nonGsmString        = "石室诗士施氏，嗜狮，誓食十狮。氏时时适市视狮。十时，适十狮适市。 是时，适施氏适市。氏视是十狮，恃矢势，使是十狮逝世。氏拾是十狮尸，适石室。石室湿，氏使侍拭石室。石室拭，氏始试食是十狮尸。食时，始识是十狮，实十石狮尸。试释是事。"
)

func TestGetCountryFromPhone(t *testing.T) {
	t.Run("australian number returns au", func(t *testing.T) {
		result, err := biz.GetCountryFromPhone(phoneNumber)
		if err != nil {
			t.Fatal("unexpected failure getting country from phone number:", err)
		}
		if result != "au" {
			t.Fatal("unexpected country returned from phone number")
		}
	})

	t.Run("return err if number is not valid", func(t *testing.T) {
		_, err := biz.GetCountryFromPhone(invalidNumber)
		if err == nil {
			t.Fatal("expected to get an error for invalid phone number")
		}
	})
}

func TestIsValidSMS(t *testing.T) {
	t.Run("valid sms returns no error", func(t *testing.T) {
		count, err := biz.IsValidSMS(nonGsmString, biz.SMSOptions{TrackLink: false, MaxParts: 4})
		if err != nil {
			t.Fatal("unexpected error when testing valid SMS message")
		}
		if count != 2 {
			t.Fatal("unexpected number of parts returned:", count)
		}
	})

	t.Run("invalid sms returns false", func(t *testing.T) {
		count, err := biz.IsValidSMS(gsmString, biz.SMSOptions{TrackLink: false, MaxParts: 4})
		if err != biz.ErrInvalidSMSTooManyParts {
			t.Fatal("expected to get error about sms having too many parts:", err)
		}
		if count != 6 {
			t.Fatal("unexpected number of parts returned:", count)
		}
	})
}

func TestIsGSMString(t *testing.T) {
	t.Run("returns true if string is gsm", func(t *testing.T) {
		result := biz.IsGSMString(gsmString)
		if result != true {
			t.Fatal("expected gsm string to return true")
		}
	})

	t.Run("returns false if string is non-gsm", func(t *testing.T) {
		result := biz.IsGSMString(nonGsmString)
		if result != false {
			t.Fatal("expected non-gsm string to return false")
		}
	})
}

func TestCountSMSParts(t *testing.T) {
	t.Run("count parts of gsm message", func(t *testing.T) {
		result := biz.CountSMSParts(gsmString, biz.SMSOptions{TrackLink: false, TrackLinkDomain: `track.lnk`})
		if result != 6 {
			t.Fatal("unexpected number of gsm message parts:", result)
		}
	})

	t.Run("count parts of non-gsm message", func(t *testing.T) {
		result := biz.CountSMSParts(nonGsmString, biz.SMSOptions{TrackLink: false, TrackLinkDomain: `track.lnk`})
		if result != 2 {
			t.Fatal("unexpected number of non gsm message parts:", result)
		}
	})

	t.Run("count parts of gsm message with link", func(t *testing.T) {
		result := biz.CountSMSParts(gsmStringWithLinks, biz.SMSOptions{TrackLink: false, TrackLinkDomain: `track.lnk`})
		if result != 2 {
			t.Fatal("unexpected number of gsm message with link parts and shortlink disabled:", result)
		}

		shortResult := biz.CountSMSParts(gsmStringWithLinks, biz.SMSOptions{TrackLink: true, TrackLinkDomain: `track.lnk`})
		if shortResult != 1 {
			t.Fatal("unexpected number of gsm message with link parts and shortlink enabled:", shortResult)
		}
	})

	t.Run("count parts of gsm message with opt out tag", func(t *testing.T) {
		result := biz.CountSMSParts(gsmStringWithOptOut, biz.SMSOptions{TrackLink: false, OptOutLinkDomain: `opted.out`})
		if result != 2 {
			t.Fatal("unexpected number of gsm message with link parts and shortlink disabled:", result)
		}
	})

}
