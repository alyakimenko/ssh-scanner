package interactive

func SSHInteractive(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	answers = make([]string, len(questions))
	for n, _ := range questions {
		answers[n] = "root"
	}

	return answers, nil
}
