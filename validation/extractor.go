package validation

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func ExtractMetadataRaw(raw []byte) (map[string]map[string][]interface{}, error) {
	meta := make(map[string][]map[string]interface{})

	if err := json.Unmarshal(raw, &meta); err != nil {
		return nil, fmt.Errorf("failed extracting meta: %w", err)
	}

	out := make(map[string]map[string][]interface{})

	for fileName, data := range meta {
		bpTypes := make(map[string]string)
		for i, obj := range data {
			if i == 0 && obj["Type"] != "BlueprintGeneratedClass" {
				break
			}

			if obj["Type"] == "BlueprintGeneratedClass" {
				superName := obj["SuperStruct"].(map[string]interface{})["ObjectName"].(string)
				_, objName := splitName(superName)
				bpTypes[obj["Name"].(string)] = objName
				continue
			}

			if obj["Properties"] != nil {
				classType := obj["Type"].(string)
				if _, ok := ignoredClasses[classType]; ok {
					continue
				}

				if _, ok := out[fileName]; !ok {
					out[fileName] = make(map[string][]interface{})
				}

				typ := bpTypes[classType]
				if typ == "" {
					typ = classType
				}

				out[fileName][typ] = append(out[fileName][typ], rewriteRecursive(obj["Properties"]))
			}
		}
	}

	return out, nil
}

var objNameRegex = regexp.MustCompile(`^(.+?)'(.+?)'$`)

func splitName(n string) (string, string) {
	matches := objNameRegex.FindStringSubmatch(n)
	return matches[1], matches[2]
}

func rewriteRecursive(obj interface{}) interface{} {
	switch b := obj.(type) {
	case map[string]interface{}:
		if mapHas("CultureInvariantString", b) {
			return b["CultureInvariantString"]
		} else if mapHas("ObjectName", b) && mapHas("ObjectPath", b) {
			_, val := splitName(b["ObjectName"].(string))
			return val
		} else if mapHas("AssetPathName", b) && mapHas("SubPathString", b) {
			return b["AssetPathName"]
		} else {
			newOut := make(map[string]interface{})
			for k, v := range b {
				newOut[k] = rewriteRecursive(v)
			}
			return newOut
		}
	case []interface{}:
		newOut := make([]interface{}, len(b))
		for i, v := range b {
			newOut[i] = rewriteRecursive(v)
		}
		return newOut
	}
	return obj
}

func mapHas(key string, mp map[string]interface{}) bool {
	_, ok := mp[key]
	return ok
}
