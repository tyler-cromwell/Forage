package twilio

type MockTwilio struct {
	From                   string
	To                     string
	OverrideComposeMessage func(int, int, string) string
	OverrideSendMessage    func(string, string, string) (string, error)
}

func (mtc *MockTwilio) ComposeMessage(quantity, quantityExpired int, url string) string {
	if mtc.OverrideComposeMessage != nil {
		return mtc.OverrideComposeMessage(quantity, quantityExpired, url)
	} else {
		return ""
	}
}

func (mtc *MockTwilio) SendMessage(phoneFrom, phoneTo, message string) (string, error) {
	if mtc.OverrideSendMessage != nil {
		return mtc.OverrideSendMessage(phoneFrom, phoneTo, message)
	} else {
		return "", nil
	}
}
