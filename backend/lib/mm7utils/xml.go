package mm7utils

import "regexp"

func ExtractEntity(regex regexp.Regexp, soap string) string {
	matches := regex.FindStringSubmatch(soap)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}
