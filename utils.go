package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func LogInfo(msgid string, msgs ...interface{}) {
	logPrefix := time.Now().UTC().Format(time.RFC3339Nano) + " [" + msgid + "]"
	msgs = append([]interface{}{logPrefix}, msgs...)
	fmt.Println(msgs...)
}

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
	var foundLang, foundTemplate, foundKey bool
	if lang, foundLang := Translations[language]; foundLang {
		if tpl, foundTemplate := lang.(map[string]interface{})[templateName]; foundTemplate {
			if value, foundKey := tpl.(map[string]interface{})[key]; foundKey {
				return value.(string), nil
			}
		}
	}
	if foundTemplate == false {
		return "", errors.New("could not find template " + templateName + " in language " + language)
	}
	if foundLang == false {
		return "", errors.New("could not find language " + language)
	}
	if foundKey == false {
		return "", errors.New("could not find key " + key + " in template " + templateName + " in language " + language)
	}

	return "", errors.New("Unexpected exception")
}
