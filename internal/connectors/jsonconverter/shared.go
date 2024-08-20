package jsonconverter

import (
	"encoding/json"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().
			SetDisplayName("API Key").
			SetDescription("Please use **test-key** as value for API Key").
			Build(),
	}).
	Build()

func convertTextToJSON(text string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//  func convertJSONToText(jsonObj map[string]interface{}) (string, error) {
//	jsonData, err := json.Marshal(jsonObj)
//	if err != nil {
//		return "", err
//	}
//	return string(jsonData), nil
//}
