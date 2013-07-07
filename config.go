/*
 * 2013-04-28
 * will@summercat.com
 *
 * a config file parser.
 *
 * A note on usage:
 * Due to the fact that we use the reflect package, you must pass in
 * the struct for which you want to parse config keys using all
 * exported fields, or this config package cannot set those fields.
 *
 * As well, key names in the config file itself are currently case
 * sensitive and must match the struct field name.
 *
 * For an example of using this package, see the test(s).
 *
 * For the types that we support parsing out of the struct, refer to
 * the populateConfig() function.
 */

package config

import (
	"bufio"
	"errors"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// readConfigFile reads a config file and returns the keys and values in
// a map. the config file syntax is:
// key = value
// lines may be commented if they begin with a '#' with only whitespace
// or no whitespace in front of the '#' character.
// lines may not have trailing '#' to be treated as comments.
func readConfigFile(path string) (map[string]string, error) {
	if len(path) == 0 {
		return nil, errors.New("invalid path")
	}

	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(fi)

	var config map[string]string = make(map[string]string)
	for {
		// XXX: what encoding is this defaulting to?
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		var parts = strings.SplitAfterN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		var key = strings.TrimSpace(parts[0])
		// XXX: if key has '=' or ' ' in it then we will trim it off
		//   since this does not restrict us to only the final two
		//   characters.
		//   though we likely don't need to support that anyway...
		key = strings.TrimRight(key, "= ")
		var value = strings.TrimSpace(parts[1])
		config[key] = value
	}
	return config, nil
}

// populateConfig takes values read from a config and uses them to fill the
// struct. the values will be converted to the struct's types as necessary.
func populateConfig(config interface{}, rawValues map[string]string) error {
	var v reflect.Value = reflect.ValueOf(config)
	// XXX: why is this needed?
	//   see the 'laws of reflection' article section on structs.
	var vElem reflect.Value = v.Elem()
	// we need the type of this struct so we can retrieve member names.
	var vType reflect.Type = vElem.Type()
	// iterate over every field of the struct.
	for i := 0; i < vElem.NumField(); i++ {
		var f reflect.Value = vElem.Field(i)
		var fieldName = vType.Field(i).Name
		// we require that we have read a value for the struct's field.
		rawValue, ok := rawValues[fieldName]
		if !ok {
			return errors.New("Field " + fieldName + " not found in config")
		}

		// we support a subset of types ('kinds' in reflect) currently.
		if f.Kind() == reflect.Int64 {
			// convert the string to an int64.
			converted, err := strconv.ParseInt(rawValue, 10, 64)
			if err != nil {
				return err
			}
			f.SetInt(converted)
			continue
		}

		if f.Kind() == reflect.Uint64 {
			converted, err := strconv.ParseUint(rawValue, 10, 64)
			if err != nil {
				return err
			}
			f.SetUint(converted)
			continue
		}

		if f.Kind() == reflect.String {
			f.SetString(rawValue)
			continue
		}

		return errors.New("Unhandled field kind: " + f.Kind().String())
	}
	return nil
}

// GetConfig reads a config file and populates a struct with what is read.
// we use the reflect package to populate the struct from the config.
// currently every member of the struct must have had a value set in the
// config. that is, every config option is required.
func GetConfig(path string, config interface{}) error {
	// we don't need to parameter check path or keys. why?
	// path will get checked when we read the config.
	// we do not need to check anything with the config as it is up to the
	// caller to ensure that they gave us a struct with members they want
	// parsed out of a config.

	// first read in the config - every key will be associated with a value
	// which is a string.
	rawValues, err := readConfigFile(path)
	if err != nil {
		return err
	}

	// fill the struct with the values read from the config.
	err = populateConfig(config, rawValues)
	if err != nil {
		return err
	}

	return nil
}
