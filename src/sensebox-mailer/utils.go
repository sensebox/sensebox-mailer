package main

import (
	"errors"
	"os"
	"strconv"
)

func getStringFromEnv(key string) (string, error) {
	str := os.Getenv(envPrefix + key)
	if len(str) == 0 {
		return "", errors.New("Please add " + envPrefix + key + " to your environment")
	}
	return str, nil
}

func getBytesFromEnv(key string) ([]byte, error) {
	str, err := getStringFromEnv(key)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func getIntFromEnv(key string) (int, error) {
	str, err := getStringFromEnv(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("Environment key " + envPrefix + key + " is not parseable as integer")
	}
	return i, nil
}

func getTranslation(language string, templateName string, key string) (string, error) {
	if lang, ok := Translations[language]; ok {
		if ok == false {
			return "", errors.New("could not find language " + language)
		}
		if tpl, ok := lang.(map[string]interface{})[templateName]; ok {
			if ok == false {
				return "", errors.New("could not find template " + templateName + " in language " + language)
			}
			if value, ok := tpl.(map[string]interface{})[key]; ok {
				if ok == false {
					return "", errors.New("could not find key " + key + " in template " + templateName + " in language " + language)
				}
				return value.(string), nil
			}
		}
	}
	return "", errors.New("shouldn't happen")
}
