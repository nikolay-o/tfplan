package diff

import (
	"fmt"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/sters/yaml-diff/yamldiff"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func Diff(file string) {
	planBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error getting before values: %v\n", err)
	}

	plan := tfjson.Plan{}
	err = plan.UnmarshalJSON(planBytes)
	if err != nil {
		log.Fatalf("error unmarshall plan file: %v\n", err)
	}

	for _, resource := range plan.ResourceChanges {
		if resource.Address == "module.charts.helm_release.charts[\"argocd\"]" {
			fmt.Println(resource.Address)
			before, err := getState(resource.Change.Before.(map[string]interface{})["values"].([]interface{}))
			if err != nil {
				log.Fatalf("error getting before values: %v\n", err)
			}
			after, err := getState(resource.Change.After.(map[string]interface{})["values"].([]interface{}))
			if err != nil {
				log.Fatalf("error getting after values: %v\n", err)

			}

			mergedBefore, err := yaml.Marshal(before)
			beforeData, err := yamldiff.Load(string(mergedBefore))

			mergedAfter, err := yaml.Marshal(after)
			afterData, err := yamldiff.Load(string(mergedAfter))

			for _, diff := range yamldiff.Do(beforeData, afterData) {
				fmt.Print(diff.Dump())
			}
		}
	}

}

// getState func will return map of before or after changes
func getState(valuesFilesSlice []interface{}) (state map[string]interface{}, err error) {
	for _, valuesFileContent := range valuesFilesSlice {
		vf := map[string]interface{}{}
		err = yaml.Unmarshal([]byte(valuesFileContent.(string)), &vf)
		if err != nil {
			return map[string]interface{}{}, err
		}
		state = mergeMaps(state, vf)
	}
	return state, nil
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
