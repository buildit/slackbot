package poll

import (
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"strconv"
)

type Vote struct {
	SelectedOption string
	Voted          bool
}

type Poll struct {
	Title      string
	Options    []string
	Attachment slack.Attachment
}

var votes map[string]Vote

var NumberToWord = map[int]string{
	1:  "one",
	2:  "two",
	3:  "three",
	4:  "four",
	5:  "five",
	6:  "six",
	7:  "seven",
	8:  "eight",
	9:  "nine",
	10: "ten",
}

func AddVote(name string, currentVote Vote) {
	votes[name] = currentVote
}

func VoteWasCast(name string) bool {
	if votes[name].Voted {
		return true
	}
	return false
}

//Convert an integer value 1-10 to the string equivalent
func convert1to10(n int) (w string) {
	if n < 20 {
		w = NumberToWord[n]
		return
	}

	r := n % 10
	if r == 0 {
		w = NumberToWord[n]
	} else {
		w = NumberToWord[n-r] + "-" + NumberToWord[r]
	}
	return
}

//Normailize the text being used to create a poll. (ex.  Slack will respond text using Smart quotes, etc).
func Normalize(in rune) rune {
	switch in {
	case '“', '‹', '”', '›':
		return '"'
	case '‘', '’':
		return '\''
	}
	return in
}

func SplitParameters(inputString string) []string {
	var out []string
	pos := 0
	quoted := false

	for i, c := range inputString {
		switch c {
		case '"':
			quoted = !quoted
		case ' ':
			if !quoted {
				out = append(out, inputString[pos:i])
				pos = i + 1
			}
		}
	}

	if pos < len(inputString) {
		if quoted {
			log.Println("missing closing quote")
		}
		out = append(out, inputString[pos:])
	}
	return removeWrappedQuotes(out)
}

func removeWrappedQuotes(inputString []string) []string {
	var newout []string
	for _, value := range inputString {
		if len(value) > 0 && value[0] == '"' {
			value = value[1:]

		}
		if len(value) > 0 && value[len(value)-1] == '"' {
			value = value[:len(value)-1]
		}
		newout = append(newout, value)
	}
	return newout
}

func CreatePoll(slicedParams []string) Poll {

	attachedOptions := []slack.AttachmentAction{}
	var attachedOptionsText string
	for i, value := range slicedParams {
		if i > 0 { //processing options
			option := slack.AttachmentAction{
				Name:  value,
				Text:  fmt.Sprintf(":%s: %s", convert1to10(i), value),
				Type:  "button",
				Style: "default",
				Value: strconv.Itoa(i),
			}
			attachedOptions = append(attachedOptions, option)
			attachedOptionsText = attachedOptionsText + fmt.Sprintf(":%s: %s \n", convert1to10(i), value)
		}
	}
	option := slack.AttachmentAction{
		Name:  "actionCancel",
		Text:  "Delete Poll",
		Type:  "button",
		Style: "danger",
		Value: "cancel",
	}
	attachedOptions = append(attachedOptions, option)

	var attachment = slack.Attachment{
		Text:       attachedOptionsText,
		Color:      "#f9a41b",
		CallbackID: "poll",
		Actions:    attachedOptions,
	}

	newPoll := Poll{
		Title:      slicedParams[0],
		Options:    slicedParams,
		Attachment: attachment,
	}

	return newPoll
}
