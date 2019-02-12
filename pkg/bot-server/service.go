package bot_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buildit/slackbot/pkg/config"
	"github.com/buildit/slackbot/pkg/poll"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var api = slack.New(config.Env.OauthToken)
var slackPoll = poll.Poll{}

func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}

func ListenAndServeSlash(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(config.Env.VerificationToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/poll":

		//Take the submitted parameters, normalize the text, and create a Slice of strings
		params := &slack.Msg{Text: s.Text}
		normalizedParams := strings.Map(poll.Normalize, params.Text)
		slicedParams := poll.SplitParameters(normalizedParams)
		fmt.Printf("Poll Submission detected with Message Paramters:%q\n", slicedParams)

		if len(slicedParams) < 1 {
			log.Printf("[ERROR] No Topic Provided for the submitted poll \n")
			w.WriteHeader(http.StatusInternalServerError)
		}
		if len(slicedParams) > 10 {
			log.Printf("[ERROR] Polling only supports up to 10 options \n")
			w.WriteHeader(http.StatusInternalServerError)
		}

		slackPoll = poll.CreatePoll(slicedParams)

		channelID, timestamp, err := api.PostMessage(s.ChannelID, slack.MsgOptionText(slicedParams[0], false), slack.MsgOptionAttachments(slackPoll.Attachment))

		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Poll '%s' successfully sent to channel %s at %s \n", slackPoll.Title, channelID, timestamp)
	}

}
func ListenAndServeInteractions(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unescape request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	var message slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Only accept message from slack with valid token
	if message.Token != config.Env.VerificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("Received Message: %s \n", jsonStr)
	callbackType := message.CallbackID
	fmt.Printf("CallbackType: %s \n", callbackType)

	switch callbackType {
	case "poll":
	}

}
func ListenAndServeEvents(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: config.Env.VerificationToken}))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			channelID, timeStamp, _ := api.PostMessage(ev.Channel, slack.MsgOptionText("Hello", false))
			fmt.Printf("Message successfully sent to channel %s at %s", channelID, timeStamp)
		}
	}

}
