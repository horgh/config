/*
 * 2013-04-28
 * will@summercat.com
 *
 * a config file parser
 */

package config

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// readConfig reads a config file and returns the keys and values in
// a map. the config file syntax is:
// key = value
// lines may be commented if they begin with a '#' with only whitespace
// or no whitespace in front of the '#' character.
// lines may not have trailing '#' to be treated as comments.
func readConfig(path string) (map[string]string, error) {
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
		var value = strings.TrimSpace(parts[1])
		config[key] = value
	}
	return config, nil
}

// GetConfig reads a config file and then verifies that each of the requested
// keys is present.
// we do not convert types - the returned map is [string]string and I am not
// sure a good way to permit mixed types there.
func GetConfig(path string, keys []string) (map[string]string, error) {
	// first read in the config - every key will be associated with a value
	// which is a string.
	config, err := readConfig(path)
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		_, exists := config[key]
		if !exists {
			return nil, errors.New("missing key in config [" + key + "]")
		}
	}

	return config, nil
}
