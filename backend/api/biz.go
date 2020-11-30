package api

func checkValidSender(sender string, configured []string) bool {
	for _, possibleSender := range configured {
		validSender := possibleSender == sender
		if validSender {
			return true
		}
	}
	return false
}
