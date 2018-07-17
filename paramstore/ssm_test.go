package paramstore

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/aws"
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

const (
	testParamName = "foo"
	testParamValue = "bar"
)

////////////////////////////
//
// mock service realization
//
////////////////////////////

type defaultSvcMock struct {
	SsmSvcDecorator
}

func (svc defaultSvcMock) GetParameter(i *ssm.GetParameterInput) (*ssm.GetParameterOutput, error){
	return &ssm.GetParameterOutput{
		Parameter:&ssm.Parameter{
			Value: aws.String(testParamValue),
			Name: aws.String(testParamName),
			Type: ssm.ParameterTypeSecureString,
			Version: aws.Int64(1),
		},
	}, nil
}

type svcMock_ParameterNotFound struct {
	SsmSvcDecorator
}

func (svc svcMock_ParameterNotFound) GetParameter(i *ssm.GetParameterInput) (*ssm.GetParameterOutput, error){
	return nil, errors.New(ssm.ErrCodeParameterNotFound)
}


func TestGetParam(t *testing.T) {
	cases := []struct {
		testCaseName string
		paramName 	 string
		SsmSvcDecorator
		expectedValue string
		expectedErrorMsg string
	}{
		{
			testCaseName: "happyPath",
			paramName: testParamName,
			SsmSvcDecorator: defaultSvcMock{},
			expectedValue: testParamValue,
			expectedErrorMsg: "",
		},
		{
			testCaseName: "parameterNotFound",
			paramName: testParamName,
			SsmSvcDecorator: svcMock_ParameterNotFound{},
			expectedValue: "",
			expectedErrorMsg: "ParameterNotFound",
		},
	}

	for i, c := range cases {
		t.Run(c.testCaseName, func(t *testing.T) {
			svcMock := c.SsmSvcDecorator
			value, err := GetParam(svcMock, c.paramName)
			failMsgTmpl := "TestGetParam testcase [%d], with name `%s` failed."
			if c.expectedErrorMsg == "" {
				assert.NoError(t, err, fmt.Sprintf(failMsgTmpl, i, c.testCaseName))
			} else {
				assert.EqualError(t, err, c.expectedErrorMsg, fmt.Sprintf(failMsgTmpl, i, c.testCaseName))
			}
			assert.Equal(t, c.expectedValue, value, fmt.Sprintf(failMsgTmpl, i, c.testCaseName))
		})
	}


}

