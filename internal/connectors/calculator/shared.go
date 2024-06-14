package calculator

import (
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type CalculatorOperationProps struct {
	FirstNumber  float64 `json:"firstNumber"`
	SecondNumber float64 `json:"secondNumber"`
}

var inputFields = map[string]*sdkcore.AutoFormSchema{
	"firstNumber": autoform.NewNumberField().
		SetDisplayName("First number").
		Build(),
	"secondNumber": autoform.NewNumberField().
		SetDisplayName("Second number").
		Build(),
}
