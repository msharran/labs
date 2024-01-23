package main

import (
	"fmt"

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

func recursiveMerge(src, dest interface{}) error {
	// if src is not a map[string]interface{}, then we
	// can directly merge it into dest
	srcMap, ok := src.(map[string]interface{})
	if !ok {
		dest = src
		return nil
	}

	// now we know that src is a map[string]interface{}
	// if dest is not a map[string]interface{}, then we
	// need to fail with an error
	destMap, ok := dest.(map[string]interface{})
	if !ok {
		return fmt.Errorf("dest is not a map[string]interface{}")
	}

	for k := range srcMap {
		switch srcMap[k].(type) {
		case map[string]interface{}:
			// check if the dest[k] is a map[string]interface{}. If it is
			// not, then we need to fail with an error

			if _, ok := destMap[k].(map[string]interface{}); !ok {
				return fmt.Errorf("dest[k] is not a map[string]interface{}")
			}

			// if the values of the keys are maps, then
			// we need to merge the keys of the maps recursively
			// as described above
			err := recursiveMerge(srcMap[k], destMap[k])
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

func mergeDependencies(src, dest map[string][]interface{}) error {
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

		// if the destDepList does not exist, then we need to add it
		destDepList, ok := destDep.([]map[string]interface{})
		if !ok {
			dest[srcDepName] = srcDep
			continue
		}

		// if the destDepList exists, then we need to merge the srcDepList
		// into the destDepList

		//
		// destDep, ok := dest[srcDepName]
		// // if the destDep does not exist, then we need to add it
		// if !ok {
		// 	dest[srcDepName] = srcDep
		// 	continue
		// }
		//
		// destDepList, ok := destDep.([]map[string]interface{})
		// // if the destDepList does not exist, then we need to add it
		// if !ok {
		// 	dest[srcDepName] = srcDep
		// 	continue
		// }
		//
		// if the destDepList exists, then we need to merge the srcDepList
		// into the destDepList

		// merge item name "_common" from srcDepList into all items in
		// destDepList except the item with name "_common" in destDepList
		if err := mergeCommonDependencyListItems(srcDepList, destDepList); err != nil {
			return err
		}

		// merge each item in srcDepList into the item with the same name
		// in destDepList
		if err := mergeDependencyListItems(srcDepList, destDepList); err != nil {
			return err
		}

	}
	return nil
}

func mergeDependencyListItems(srcDepList, destDepList []map[string]interface{}) error {
	// merge each item in srcDepList into the item with the same name
	// in destDepList
	return nil
}

func mergeCommonDepItemIntoOthers(srcDepList, destDepList []map[string]interface{}) error {
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
					err := recursiveMerge(srcDepItemValue, destDepItemValue)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}
	}

	return nil
}

func main() {
	// mergeFunc := koanf.WithMergeFunc(MergeFunc)
	mergeFunc := koanf.WithMergeFunc(func(src, dest map[string]interface{}) error {
		for k := range src {
			if k == "dependencies" {
				// get the dependencies from src and dest
				srcDeps, ok := src[k].(map[string][]interface{})
				if !ok {
					return fmt.Errorf("src dependencies is not a map")
				}

				destDeps, ok := dest[k].(map[string][]interface{})
				if !ok {
					return fmt.Errorf("dest dependencies is not a map")
				}

				// merge the dependencies from src into dest
				err := mergeDependencies(srcDeps, destDeps)
				if err != nil {
					return fmt.Errorf("failed to merge dependencies: %w", err)
				}
			}

			// if the key is not "dependencies", then we can merge it
			// directly
			dest[k] = src[k]
		}
		return nil
	})

	k := koanf.New(".")
	k.Load(file.Provider("config.yaml"), yaml.Parser(), mergeFunc)
	k.Load(file.Provider("config2.yaml"), yaml.Parser(), mergeFunc)

	o, err := k.Marshal(yaml.Parser())
	check(err)

	println(string(o))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
