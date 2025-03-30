package ottoTTS

import (
	"encoding/json"
	"github.com/mozillazg/go-pinyin"
	"log"
	"os"
	"ottoTTS/wavHandler"
	"strings"
)

type expression struct {
	Expression []string `json:"expression"`
	Otto       string   `json:"otto"`
}

type dictionary struct {
	Expressions []expression `json:"expressions"`
	Letters     []expression `json:"letters"`
	Numbers     []expression `json:"numbers"`
}

var dict dictionary

func getDictionary() (dictionary, error) {
	dictFile, err := os.ReadFile("./assets/dictionary.json")
	if err != nil {
		return dictionary{}, err
	}

	var dictionary dictionary
	err = json.Unmarshal(dictFile, &dictionary)
	if err != nil {
		return dictionary, err
	}

	return dictionary, nil
}

func expressionMatch(str string) (string, int) {
	for _, exprs := range dict.Expressions {
		for _, expr := range exprs.Expression {
			if strings.Contains(str, expr) {
				return exprs.Otto, len(expr)
			}
		}
	}
	return "", 0
}

func stringToSlices(words string, expressionOverride bool) []wavHandler.Slice {
	var slices []wavHandler.Slice
	index := 0
	for index < len(words) {
		if expressionOverride {
			matchedWords, length := expressionMatch(words[index:])
			if length != 0 {
				slice := wavHandler.Slice{Category: "expressions", Content: matchedWords}
				slices = append(slices, slice)
				index += length
				continue
			}
		}
		a := pinyin.NewArgs()
		c := pinyin.Pinyin(string(words[index]), a)
		slice := wavHandler.Slice{Category: "characters", Content: c[0][0]}
		slices = append(slices, slice)
		index++
	}

	return slices
}

func speech(s string, expressionOverride bool) ([]byte, error) {
	return wavHandler.GetSpeech(stringToSlices(s, expressionOverride))
}

func InitializeTTS() {
	log.Println("初始化otto文字转语音引擎")

	log.Println("加载otto词典")
	var err error
	dict, err = getDictionary()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("otto词典加载完成")

	log.Println("otto文字转语音引擎初始化完成")
}
