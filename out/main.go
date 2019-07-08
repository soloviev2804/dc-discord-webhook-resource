package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	//VERSION -
	VERSION string
)

//Input - Struct that represents the input to out
type Input struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type Source struct {
	Discord Discord `json:"discord"`
}

type Params struct {
	MessageText string `json:"messageText"`
}

type Discord struct {
	Token     string `json:"token"`
	WebhookID string `json:"webhookId"`
}

//MetadataItem - metadata within output
type MetadataItem struct {
	Name  string
	Value string
}

//Output - represents output from out
type Output struct {
	Version struct {
		Time time.Time
	} `json:"version"`
	Metadata []MetadataItem
}

func main() {
	inbytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s - Error reading stdin - %s", VERSION, err.Error()))
		os.Exit(1)
	}
	output, err := Execute(VERSION, inbytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s - Error during execute - %s", VERSION, err.Error()))
		os.Exit(1)
	}
	fmt.Println(output)
}

func Execute(version string, input []byte) (string, error) {

	var indata Input

	err := json.Unmarshal(input, &indata)
	if err != nil {
		return "", errors.Wrap(err, "unmarshalling input")
	}

	err = validateConfiguration(indata)
	if err != nil {
		return "", errors.Wrap(err, "Invalid configuration")
	}

	source := indata.Source
	params := indata.Params
	discord, err := discordgo.New("Bot " + source.Discord.Token)
	defer discord.Close()

	text := replaceTokens(params.MessageText)
	webhookParams := discordgo.WebhookParams{
		Content:  text,
		Username: "Concourse CI",
	}
	err = discord.WebhookExecute(source.Discord.WebhookID, source.Discord.Token, false, &webhookParams)
	if err != nil {
		return "", errors.Wrap(err, "Can't execute webhook")
	}

	var outdata Output
	outdata.Version.Time = time.Now().UTC()
	outdata.Metadata = []MetadataItem{
		{Name: "messageText", Value: text},
	}
	outbytes, err := json.Marshal(outdata)
	if err != nil {
		return "", errors.Wrap(err, "Error Marshalling JSON:")
	}

	return string(outbytes), nil
}

func validateConfiguration(indata Input) error {
	if indata.Source.Discord.Token == "" {
		return errors.New(`missing required field "source.discord.token"`)
	}

	if indata.Source.Discord.WebhookID == "" {
		return errors.New(`missing required field "source.discord.webhookId"`)
	}

	if indata.Params.MessageText == "" {
		return errors.New(`missing required field "params.messageText" `)
	}

	return nil
}

func replaceTokens(sourceString string) string {
	var buildTokens = map[string]string{
		"${BUILD_ID}":            os.Getenv("BUILD_ID"),
		"${BUILD_NAME}":          os.Getenv("BUILD_NAME"),
		"${BUILD_JOB_NAME}":      os.Getenv("BUILD_JOB_NAME"),
		"${BUILD_PIPELINE_NAME}": os.Getenv("BUILD_PIPELINE_NAME"),
		"${ATC_EXTERNAL_URL}":    os.Getenv("ATC_EXTERNAL_URL"),
		"${BUILD_TEAM_NAME}":     os.Getenv("BUILD_TEAM_NAME"),
	}
	for k, v := range buildTokens {
		sourceString = strings.Replace(sourceString, k, v, -1)
	}
	return sourceString
}
