package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
)

var questions = []*survey.Question{
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Choose a type:",
			Options: []string{"function", "variable"},
			Default: "function",
		},
		Validate: survey.Required,
	},
	{
		Name: "overview",
		Prompt: &survey.Multiline{
			Message: "Enter an overview:",
		},
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		Type     string
		Overview string
	}{}

	err := survey.Ask(questions, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.openai.com/v1/completions",
		bytes.NewBuffer(
			[]byte(`{
				"model": "text-davinci-003",
				"prompt": "Say this is a test",
				"temperature": 0,
				"max_tokens": 70
			}`)),
	)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req.Header.Add("Authorization", "Bearer ")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	resultStruct := struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Text         string `json:"text"`
			Index        int    `json:"index"`
			Logprobs     string `json:"logprobs"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}{}

	if err := json.Unmarshal(result, &resultStruct); err != nil {
		fmt.Println(err)
		return
	} else if len(resultStruct.Choices) == 0 {
		fmt.Println("No choices")
		return
	}

	fmt.Println(resultStruct.Choices[0].Text)

}