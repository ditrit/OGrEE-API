package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	u "p3/utils"
	"strings"
	"testing"
)

func TestValidateJsonSchema(t *testing.T) {
	testingEntities := []int{u.SITE, u.BLDG, u.ROOM, u.RACK, u.DEVICE, u.GROUP, u.BLDGTMPL, u.OBJTMPL}
	for _, entInt := range testingEntities {
		entStr := u.EntityToString(entInt)
		println("*** Testing " + entStr)
		var obj map[string]interface{}
		data, e := ioutil.ReadFile("schemas/" + entStr + "_schema.json")
		if e != nil {
			t.Error(e.Error())
		}
		json.Unmarshal(data, &obj)
		if entInt == u.BLDGTMPL || entInt == u.OBJTMPL { // 3 good examples
			for i := 0; i < 3; i++ {
				resp, ok := validateJsonSchema(entInt, obj["examples"].([]interface{})[i].(map[string]interface{}))
				if !ok {
					t.Errorf("Error validating json schema: %s", resp)
				}
			}
		} else { // only first example is good
			resp, ok := validateJsonSchema(entInt, obj["examples"].([]interface{})[0].(map[string]interface{}))
			if !ok {
				t.Errorf("Error validating json schema: %s", resp)
			}
		}
	}
}

func TestErrorValidateJsonSchema(t *testing.T) {
	expectedErrors := map[string][]string{
		"site1":     {"missing properties: 'domain'", "/attributes/reservedColor does not match pattern"},
		"building1": {"missing properties: 'posXYUnit'", "/attributes/height expected string, but got number"},
		"room1":     {"additionalProperties 'banana' not allowed", "/attributes/axisOrientation value must be one of"},
		"rack1":     {"/attributes/posXYZ does not match pattern", "/attributes/heightUnit value must be one of"},
		"device1":   {"missing properties: 'template'", "/description expected array, but got string"},
		"group1":    {"/attributes missing properties: 'content'"},
		"obj_template3": {
			"/slug does not match pattern",
			"/attributes/vendor expected string, but got number",
			"if-then failed",
			"/slots/0/elemOrient value must be one of ",
			"/slots/1/elemPos maximum 3 items required, but found 4 items",
			"/slots/1/elemSize minimum 3 items required, but found 2 items",
			"/slots/1/elemOrient value must be one of",
			"/slots/2 missing properties: 'elemOrient'",
			"/slots/2/labelPos value must be one of ",
			"/slots/3/color does not match pattern",
		},
		"obj_template4": {
			"/components/0/elemPos minimum 3 items required, but found 0 items",
			"/components/1/elemOrient value must be one of",
			"/components/1/color does not match pattern",
			"/components/3/labelPos value must be one of",
			"if-else failed",
			"/slots/0/elemOrient value must be one of",
		},
		"bldg_template3": {
			"/sizeWDHm minimum 3 items required, but found 2 items",
			"/vertices/2 minimum 2 items required, but found 1 items",
			"/center minimum 2 items required, but found 0 items",
		},
	}
	testingEntities := []int{u.SITE, u.BLDG, u.ROOM, u.RACK, u.DEVICE, u.GROUP, u.BLDGTMPL, u.OBJTMPL}
	for _, entInt := range testingEntities {
		entStr := u.EntityToString(entInt)
		println("*** Testing " + entStr)
		var obj map[string]interface{}
		data, e := ioutil.ReadFile("schemas/" + entStr + "_schema.json")
		if e != nil {
			t.Error(e.Error())
		}
		json.Unmarshal(data, &obj)

		startExample := 1
		endExample := 1
		if entInt == u.BLDGTMPL {
			startExample = 3
			endExample = 3
		} else if entInt == u.OBJTMPL {
			startExample = 3
			endExample = 4
		}

		for i := startExample; i <= endExample; i++ {
			entObjName := fmt.Sprintf("%s%d", entStr, i)
			println(entObjName)
			resp, ok := validateJsonSchema(entInt, obj["examples"].([]interface{})[i].(map[string]interface{}))
			if ok {
				t.Errorf("Validated json schema that should have these errors: %v", expectedErrors[entObjName])
			} else {
				if len(resp["errors"].([]string)) != len(expectedErrors[entObjName]) {
					t.Errorf("Validation errors do not correspond expected errors:\n%v\nGot:\n%v", expectedErrors[entObjName], resp["errors"].([]string))
				} else {
					for _, expected := range expectedErrors[entObjName] {
						if !contains(resp["errors"].([]string), expected) {
							t.Errorf("Validation errors do not correspond expected errors:\n%v\nGot:\n%v", expectedErrors[entObjName], resp["errors"].([]string))
						}
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
