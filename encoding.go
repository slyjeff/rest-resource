package resource

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"sort"
)

func (fd FormattedData) MarshalJSON() ([]byte, error) {
	if reflect.TypeOf(fd.Value).Kind() == reflect.String {
		return []byte("\"" + fd.FormattedString() + "\""), nil
	}

	return []byte(fd.FormattedString()), nil
}

func (l *Link) MarshalJSON() ([]byte, error) {
	json := "{\"href\": \"" + l.Href + "\"}"

	if l.Verb != "GET" {
		json = addToJson(json, "verb", quoted(l.Verb))
	}

	if l.IsTemplated {
		json = addToJson(json, "templated", "true")
	}

	if len(l.Parameters) > 0 {
		parametersJson := "{}"
		for _, parameter := range l.Parameters {
			parameterValues := "{}"

			if parameter.DefaultValue != "" {
				parameterValues = addToJson(parameterValues, "default", quoted(parameter.DefaultValue))
			}

			if parameter.ListOfValues != "" {
				parameterValues = addToJson(parameterValues, "listOfValues", quoted(parameter.ListOfValues))
			}

			if parameter.DataType != "" {
				parameterValues = addToJson(parameterValues, "dataType", quoted(parameter.DataType))
			}

			parametersJson = addToJson(parametersJson, parameter.Name, parameterValues)
		}
		json = addToJson(json, "parameters", parametersJson)
	}

	return []byte(json), nil
}

func quoted(s string) string {
	return "\"" + s + "\""
}

func addToJson(json, name, value string) string {
	nameValue := "\"" + name + "\":" + value
	if json == "{}" {
		return json[:len(json)-1] + nameValue + "}"
	}
	return json[:len(json)-1] + "," + nameValue + "}"
}

//goland:noinspection GoMixedReceiverTypes
func (r Resource) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	tokens := make([]xml.Token, 0)
	tokens = append(tokens, xml.StartElement{Name: xml.Name{Local: "resource"}})

	keys := make([]string, 0)
	for k := range r.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tokens = addXmlTokens(tokens, k, r.Values[k])
	}

	tokens = append(tokens, xml.EndElement{Name: xml.Name{Local: "resource"}})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	return e.Flush()
}

func addXmlTokens(tokens []xml.Token, k string, v interface{}) []xml.Token {
	tokens = append(tokens, xml.StartElement{Name: xml.Name{Local: k}})

	if formattedData, ok := v.(FormattedData); ok {
		v = formattedData.FormattedString()
		tokens = append(tokens, xml.CharData(fmt.Sprint(v)))
	} else if slice, ok := v.([]MappedData); ok {
		tokens = addSliceXmlTokens(tokens, slice)
	} else if md, ok := v.(MappedData); ok {
		tokens = addMapDataXmlTokens(tokens, md)
	} else {
		tokens = append(tokens, xml.CharData(fmt.Sprint(v)))
	}

	return append(tokens, xml.EndElement{Name: xml.Name{Local: k}})
}

func addSliceXmlTokens(tokens []xml.Token, slice []MappedData) []xml.Token {
	for _, md := range slice {
		tokens = append(tokens, xml.StartElement{Name: xml.Name{Local: "Value"}})
		tokens = addMapDataXmlTokens(tokens, md)
		tokens = append(tokens, xml.EndElement{Name: xml.Name{Local: "Value"}})
	}
	return tokens
}

func addMapDataXmlTokens(tokens []xml.Token, md MappedData) []xml.Token {
	keys := make([]string, 0)
	for k := range md {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tokens = addXmlTokens(tokens, k, md[k])
	}

	return tokens
}
