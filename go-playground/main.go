package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// dependencies:
//
//	s3:
//	  - name: test-s3-bucket
//	    authz:
//	    - resource: snape.service
//	      accessc
//	        - read
//	        - write
//	    - resource: locksmith.service
//	      access:
//	        - read
//	        - write
//	  - name: test-s3-bucket-2

//	s3:
//	  - name: test-s3-bucket
//	    authz:
//	    - resource: snape.service
//	      access:
//	        - read
//	    - resource: locksmith.service
//	      access:
//	        - read
//	        - write
//	  - name: test-s3-bucket-2

// RecursiveMapMerge merges two maps recursively.
// - if the values of the keys are maps, then we need to merge the
//   keys of the maps recursively
// - if the values of the keys are lists or primitives, then we need to
//   replace the value in dest with the value in src
// - if the values of the keys are different types, then we need
//   to fail with an error
// - if the keys are missing in dest, then we need to add the keys
//   from src to dest
// - if the keys are missing in src, then we need to skip them
//   (i.e. do nothing)

var ErrIsNotMap = fmt.Errorf("value is not a map")

func mergeMapKeysRecursively(src, dest interface{}) error {
	// if src is not a map[string]interface{}, then we
	// can directly merge it into dest
	srcMap, ok := src.(map[string]interface{})
	if !ok {
		return fmt.Errorf("src is not a map[string]interface{}: %w", ErrIsNotMap)
	}

	if dest == nil {
		return fmt.Errorf("dest is nil: %w", ErrIsNotMap)
	}

	destMap, ok := dest.(map[string]interface{})
	if !ok {
		return fmt.Errorf("dest is not a map[string]interface{}: %w", ErrIsNotMap)
	}

	for k := range srcMap {
		switch srcMap[k].(type) {
		case map[string]interface{}:
			// if the key is not in dest, then we need to add it

			// if the values of the keys are maps, then
			// we need to merge the keys of the maps recursively
			// as described above
			err := mergeMapKeysRecursively(srcMap[k], destMap[k])
			if err != nil {
				return err
			}

		default:
			// if the values of the keys are primitives or lists, then
			// we need to replace the value in dest with the value in src
			destMap[k] = srcMap[k]
		}
	}
	return nil
}

func mergeCommonDependencyListItems(srcDepList, destDepList []map[string]interface{}) error {
	return nil
}

func mergeDependencies(src, dest map[string][]map[string]interface{}) error {
	// src and dest are both maps which represent the "dependencies" key
	// in the yaml file
	// loop through the keys in src
	// each key is a service name and the value is a list of maps
	// each map in the list has a "name" key. If the value of the "name"
	// key is "_common", then we need to merge the keys of that map
	// into all the other maps in the dest list except the map with the
	// "_common" name. Then delete both src and dest maps with the "_common"
	// name.
	// If the value of the "name" key is missing, then we need to fail with
	// an error.
	// If the value of the "name" key is not "_common", then we need to
	// merge the keys of that map into the map in dest with the same name
	// recursively.
	for srcDepName, srcDep := range src {

		// if the destDep does not exist, then we need to add it
		destDep, ok := dest[srcDepName]
		if !ok {
			dest[srcDepName] = srcDep
			continue
		}

		// now we know that destDep exists, Example, "s3" in dest exists

		// if there is a "_common" item in srcDep, then we need to merge
		// the keys of this element into all the elements of the []interface{}
		// in dest except the element with the "_common" name in dest (if it exists)
		if err := mergeCommonSrcDepIntoAllDstDep(srcDep, destDep); err != nil {
			return err
		}

		// merge each item in srcDep into the item with the same name
		// in destDep
		if err := mergeMatchingSrcDepIntoDstDep(srcDep, destDep); err != nil {
			return err
		}

		// remove the "_common" item from both srcDep and destDep
		if err := removeCommonItemFromDeps(srcDep); err != nil {
			return fmt.Errorf("failed to remove common item from srcDep: %w", err)
		}

		if err := removeCommonItemFromDeps(destDep); err != nil {
			return fmt.Errorf("failed to remove common item from destDep: %w", err)
		}

	}
	return nil
}

func removeCommonItemFromDeps(depList []map[string]interface{}) error {
	for i, depItem := range depList {
		depItemName, ok := depItem["name"].(string)
		if !ok {
			return fmt.Errorf("name not found for %s", depItem)
		}

		if depItemName == "_common" {
			depList = append(depList[:i], depList[i+1:]...)
		}
	}
	return nil
}

func mergeMatchingSrcDepIntoDstDep(srcDepList, destDepList []map[string]interface{}) error {
	// merge each item in srcDepList into the item with the same name
	// in destDepList

	for _, srcDepItem := range srcDepList {
		srcDepItemName, ok := srcDepItem["name"].(string)
		if !ok {
			return fmt.Errorf("name not found for %s", srcDepItem)
		}

		// if the name is "_common", then we need to skip it
		if srcDepItemName == "_common" {
			continue
		}

		// if the name is not "_common", then we need to merge it
		// into the item with the same name in destDepList
		for _, destDepItem := range destDepList {
			destDepItemName, ok := destDepItem["name"].(string)
			if !ok {
				return fmt.Errorf("name not found for %s", destDepItem)
			}

			if destDepItemName == srcDepItemName {
				fmt.Fprintln(os.Stderr, "found matching name in destDepItem", srcDepItemName, destDepItemName)
				// - if srcDepItemKey and destDepItemKey are the same
				//   then recursively merge it by calling mergeItemRecursively
				for srcDepItemKey, srcDepItemValue := range srcDepItem {
					if srcDepItemKey == "name" {
						continue
					}

					destDepItemValue, ok := destDepItem[srcDepItemKey]
					if !ok {
						// if srcDepItemKey is not in destDepItem, then
						// add it
						destDepItem[srcDepItemKey] = srcDepItemValue
						continue
					}

					// - if the values of the keys are maps, then
					// we need to merge the keys of the maps recursively
					// - if the values of the keys are lists or primitives, then
					// we need to replace the value in dest with the value in src
					// - if the values of the keys are different types, then we need
					// to fail with an error

					fmt.Fprintln(os.Stderr, "found matching key in destDepItem", srcDepItemKey, destDepItemValue, ok)
					err := mergeMapKeysRecursively(srcDepItemValue, destDepItemValue)
					if err != nil {
						if errors.Is(err, ErrIsNotMap) {
							destDepItem[srcDepItemKey] = srcDepItemValue
							continue
						}
						return err
					}

				}
			}
		}
	}

	return nil
}

func mergeCommonSrcDepIntoAllDstDep(srcDepList, destDepList []map[string]interface{}) error {
	for _, srcDepItem := range srcDepList {
		srcDepItemName, ok := srcDepItem["name"].(string)
		if !ok {
			return fmt.Errorf("name not found for %s", srcDepItem)
		}

		// if the name is "_common", then we need to merge
		// the keys of this element into all the elements
		// of the []interface{} in dest except the element with
		// the "_common" name in dest (if it exists)

		if srcDepItemName == "_common" {
			fmt.Fprintln(os.Stderr, "found _common in srcDepItem")
			for srcDepItemKey, srcDepItemValue := range srcDepItem {
				if srcDepItemKey == "name" {
					continue
				}

				for _, destDepItem := range destDepList {
					destDepItemName, ok := destDepItem["name"].(string)
					if !ok {
						return fmt.Errorf("name not found for %s", destDepItem)
					}

					if destDepItemName == "_common" {
						continue
					}

					// - if srcDepItemKey and destDepItemKey are the same
					//   then recursively merge it
					destDepItemValue, ok := destDepItem[srcDepItemKey]
					if !ok {
						destDepItem[srcDepItemKey] = srcDepItemValue
						continue
					}

					// - if srcDepItemKey is not in destDepItem, then
					//   add it
					err := mergeMapKeysRecursively(srcDepItemValue, destDepItemValue)
					if err != nil {
						if errors.Is(err, ErrIsNotMap) {
							destDepItem[srcDepItemKey] = srcDepItemValue
							continue
						}
						return err
					}
				}
			}
			return nil
		}
	}

	return nil
}

func parseDependencies(deps interface{}) (map[string][]map[string]interface{}, error) {
	srcDeps, ok := deps.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("dependencies is not a map")
	}

	out := make(map[string][]map[string]interface{}, len(srcDeps))
	for k := range srcDeps {
		dd, ok := srcDeps[k].([]interface{})
		if !ok {
			return nil, fmt.Errorf("dependency %s is not a list", k)
		}

		for _, d := range dd {
			dd, ok := d.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("dependency %s is not a []map[string]interface{}", k)
			}

			out[k] = append(out[k], dd)
		}

	}
	return out, nil
}

func main() {
	// mergeFunc := koanf.WithMergeFunc(MergeFunc)
	mergeFunc := koanf.WithMergeFunc(func(src, dest map[string]interface{}) error {
		for k := range src {
			if k == "dependencies" {
				// get the dependencies from src and dest
				srcDeps, err := parseDependencies(src[k])
				if err != nil {
					fmt.Fprintln(os.Stderr, "failed to parse src dependencies")
					continue
				}

				if _, ok := dest[k]; !ok {
					dest[k] = src[k]
					continue
				}

				destDeps, err := parseDependencies(dest[k])
				if err != nil {
					fmt.Fprintln(os.Stderr, "failed to parse dest dependencies")
					continue
				}

				// merge the dependencies from src into dest
				err = mergeDependencies(srcDeps, destDeps)
				if err != nil {
					fmt.Fprintln(os.Stderr, "failed to merge dependencies: ", err)
					continue
				}

				continue
			}

			// if the key is not "dependencies", then we can merge it
			// directly
			dest[k] = src[k]
		}
		return nil
	})

	k := koanf.New(".")
	fmt.Fprintln(os.Stderr, "loading config.yaml")
	k.Load(file.Provider("config.yaml"), yaml.Parser(), mergeFunc)
	fmt.Fprintln(os.Stderr, "loading config2.yaml")
	k.Load(file.Provider("config2.yaml"), yaml.Parser(), mergeFunc)
	fmt.Fprintln(os.Stderr, "loading config3.yaml")
	k.Load(file.Provider("config3.yaml"), yaml.Parser(), mergeFunc)

	o, err := k.Marshal(yaml.Parser())
	check(err)

	fmt.Println(string(o))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
