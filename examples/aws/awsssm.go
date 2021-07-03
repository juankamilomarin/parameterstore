package main

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	awsRegion = "your_aws_region"
)

var (
	awsSSM     *ssm.SSM
	awsSSMOnce sync.Once
)

func getSSMClient() *ssm.SSM {
	awsSSMOnce.Do(func() {
		config := &aws.Config{Region: aws.String(awsRegion)}
		sess := session.Must(session.NewSession(config))
		awsSSM = ssm.New(sess)
	})
	return awsSSM
}

type SSMParameterStore struct{}

func (ps SSMParameterStore) GetParams(paramMap map[string]string) error {
	ssmsvc := getSSMClient()
	input := &ssm.GetParametersInput{WithDecryption: aws.Bool(true)}
	for paramName := range paramMap {
		input.Names = append(input.Names, aws.String(paramName))
	}
	result, err := ssmsvc.GetParameters(input)
	if err != nil {
		return err
	}
	for _, p := range result.Parameters {
		paramMap[*p.Name] = *p.Value
	}
	return nil
}
