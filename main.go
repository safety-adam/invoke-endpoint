package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
)

var region string = "ap-southeast-2"
var endpointName string = "hf-llm-falcon-7b-instruct-bf16-2024-02-15-11-20-55-820"

func main() {
	for {
		input := getInput()
		resp := invoke(input)
		fmt.Println(resp)
	}
}

func getInput() string {
	// Get input from the console
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">")
	input, _ := reader.ReadString('\n')
	return input

}

func invoke(prompt string) string {
	s := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	sagemaker := sagemakerruntime.New(s)
	endpointName := aws.String(endpointName)

	jsonData := map[string]interface{}{
		"inputs": prompt,
		"parameters": map[string]interface{}{
			"max_new_tokens":   100,
			"top_k":            10,
			"return_full_text": false,
			"do_sample":        false,
		},
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	payload := &sagemakerruntime.InvokeEndpointInput{
		EndpointName: endpointName,
		ContentType:  aws.String("application/json"),
		Body:         jsonBytes,
	}
	resp, err := sagemaker.InvokeEndpoint(payload)

	if err != nil {
		panic(err)
	}

	return string(resp.Body)
}
