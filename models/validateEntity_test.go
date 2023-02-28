package models

import (
	"encoding/json"
	"io/ioutil"
	u "p3/utils"
	"strings"
	"testing"
)

func TestValidateJsonSchema(t *testing.T) {
	testingEntities := []int{u.SITE, u.BLDG, u.ROOM, u.RACK, u.DEVICE, u.GROUP, u.BLDGTMPL}
	for _, entInt := range testingEntities {
		entStr := u.EntityToString(entInt)
		println("*** Testing " + entStr)
		var obj map[string]interface{}
		data, e := ioutil.ReadFile("schemas/" + entStr + "_schema.json")
		if e != nil {
			t.Error(e.Error())
		}
		json.Unmarshal(data, &obj)
		resp, ok := validateJsonSchema(entInt, obj["examples"].([]interface{})[0].(map[string]interface{}))
		if !ok {
			t.Errorf("Error validating json schema: %s", resp)
		}
	}
}

func TestErrorValidateJsonSchema(t *testing.T) {
	expectedErrors := map[string][]string{
		"site":     {"missing properties: 'domain'", "/attributes/reservedColor does not match pattern"},
		"building": {"missing properties: 'posXYUnit'", "/attributes/height expected string, but got number"},
		"room":     {"additionalProperties 'banana' not allowed", "/attributes/axisOrientation value must be one of"},
		"rack":     {"/attributes/posXYZ does not match pattern", "/attributes/heightUnit value must be one of"},
		"device":   {"missing properties: 'template'", "/description expected array, but got string"},
		"group":    {"/attributes missing properties: 'content'"},
	}
	testingEntities := []int{u.SITE, u.BLDG, u.ROOM, u.RACK, u.DEVICE, u.GROUP}
	for _, entInt := range testingEntities {
		entStr := u.EntityToString(entInt)
		println("*** Testing " + entStr)
		var obj map[string]interface{}
		data, e := ioutil.ReadFile("schemas/" + entStr + "_schema.json")
		if e != nil {
			t.Error(e.Error())
		}
		json.Unmarshal(data, &obj)
		resp, ok := validateJsonSchema(entInt, obj["examples"].([]interface{})[1].(map[string]interface{}))
		if ok {
			t.Errorf("Validated json schema that should have these errors: %v", expectedErrors[entStr])
		} else {
			if len(resp["errors"].([]string)) != len(expectedErrors[entStr]) {
				t.Errorf("Validation errors do not correspond expected errors: %v", expectedErrors[entStr])
			} else {
				for _, expected := range expectedErrors[entStr] {
					if !contains(resp["errors"].([]string), expected) {
						t.Errorf("Validation errors do not correspond expected errors: %v", expectedErrors[entStr])
					}
				}
			}
		}
	}
}

// helper function
func contains(slice []string, elem string) bool {
	for _, e := range slice {
		if strings.Contains(e, elem) {
			return true
		}
	}
	return false
}
