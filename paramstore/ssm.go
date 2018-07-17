package paramstore

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

// SsmSvcDecorator hide aws internals and make it testable
type SsmSvcDecorator interface {
	GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
}

type SsmSvc struct {
	ssm.SSM
}

func (svc SsmSvc) GetParameter(i *ssm.GetParameterInput) (*ssm.GetParameterOutput, error){
	return svc.GetParameterRequest(i).Send()
}

// GetParam return decrypted value associated with the name or ssm.ErrCodeParameterNotFound if there is no such
func GetParam(svc SsmSvcDecorator, paramName string) (string, error){
	output, err := svc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(paramName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *output.Parameter.Value, err
}


func NewSvc(configs ...external.Config) (*SsmSvc, error) {
	cfg, err := external.LoadDefaultAWSConfig(configs)
	if err != nil {
		return nil, err
	}
	svc := ssm.New(cfg)
	svcDecorator := &SsmSvc{*svc}
	return svcDecorator, nil
}




