package ottoTTS

import (
	"encoding/json"
	"github.com/gujial/ottoTTS/wavHandler"
	"github.com/mozillazg/go-pinyin"
	"log"
	"os"
	"strings"
	"unicode"
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
			if strings.HasPrefix(str, expr) {
				// 修正：计算 expr 的 rune 长度，而不是字节长度
				exprRuneLen := len([]rune(expr))
				return exprs.Otto, exprRuneLen
			}
		}
	}
	return "", 0
}

func stringToSlices(words string, expressionOverride bool) []wavHandler.Slice {
	var slices []wavHandler.Slice
	index := 0
	wordRunes := []rune(words) // 使用 rune 处理 UTF-8 字符

	for index < len(wordRunes) {
		char := wordRunes[index]

		// 处理表达式匹配
		if expressionOverride {
			log.Println("判断", string(wordRunes[index:])) // 日志输出剩余字符串
			matchedWords, length := expressionMatch(string(wordRunes[index:]))
			log.Println("匹配长度", length)

			if length != 0 {
				slices = append(slices, wavHandler.Slice{Category: "expressions", Content: matchedWords})
				index += length // 这里 length 是字符数，不是字节数
				continue
			}
		}

		// 处理中文字符
		if unicode.Is(unicode.Han, char) {
			a := pinyin.NewArgs()
			c := pinyin.Pinyin(string(char), a)
			slices = append(slices, wavHandler.Slice{Category: "characters", Content: c[0][0]})

			// 处理英文字母（单个）
		} else if unicode.IsLetter(char) {
			slices = append(slices, wavHandler.Slice{Category: "letters", Content: string(char)})

			// 处理数字（单个）
		} else if unicode.IsDigit(char) {
			slices = append(slices, wavHandler.Slice{Category: "numbers", Content: string(char)})

			// 处理标点、空格、符号等
		} else {
			slices = append(slices, wavHandler.Slice{Category: "others", Content: string(char)})
		}

		index++ // 这里 index 递增的是字符，不是字节
	}

	return slices
}

func Speech(s string, expressionOverride bool) ([]byte, error) {
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
