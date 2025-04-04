package ottoTTS

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
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

type config struct {
	ExpressionOverride bool `toml:"expression_override"`
	Debug              bool `toml:"Debug"`
}

var dict dictionary
var cfg config

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

func lettersMatch(l string) string {
	for _, letters := range dict.Letters {
		for _, letter := range letters.Expression {
			if strings.HasPrefix(letter, l) {
				return letters.Otto
			}
		}
	}
	return ""
}

func numbersMatch(n string) string {
	for _, numbers := range dict.Numbers {
		for _, number := range numbers.Expression {
			if strings.HasPrefix(n, number) {
				return numbers.Otto
			}
		}
	}
	return ""
}

func buildSlices(otto string, category string) []wavHandler.Slice {
	text := strings.Fields(otto)
	var slices []wavHandler.Slice
	for _, word := range text {
		slices = append(slices, wavHandler.Slice{
			Category: category,
			Content:  word,
		})
	}
	return slices
}

func stringToSlices(words string, expressionOverride bool) []wavHandler.Slice {
	var slices []wavHandler.Slice
	index := 0
	wordRunes := []rune(words) // 使用 rune 处理 UTF-8 字符

	for index < len(wordRunes) {
		char := wordRunes[index]

		// 处理表达式匹配
		if expressionOverride {
			matchedWords, length := expressionMatch(string(wordRunes[index:]))
			if cfg.Debug {
				log.Println("判断", string(wordRunes[index:])) // 日志输出剩余字符串
				log.Println("匹配长度", length)
			}

			if length != 0 {
				slices = append(slices, buildSlices(matchedWords, "expressions")...)
				index += length
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
			slices = append(slices, buildSlices(lettersMatch(string(char)), "letters")...)

			// 处理数字（单个）
		} else if unicode.IsDigit(char) {
			slices = append(slices, buildSlices(numbersMatch(string(char)), "numbers")...)

			// 处理标点、空格、符号等
		} else {
			slices = append(slices, wavHandler.Slice{Category: "others", Content: string(char)})
		}

		index++ // 这里 index 递增的是字符，不是字节
	}

	return slices
}

func Speech(s string) ([]byte, error) {
	return wavHandler.GetSpeech(stringToSlices(s, cfg.ExpressionOverride))
}

func loadConfig() {
	_, err := toml.DecodeFile("./config.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Debug {
		log.Println("调试模式已开启")
	}

	if cfg.ExpressionOverride {
		log.Println("启用otto短语覆盖")
	}
}

func InitializeTTS() {
	log.Println("初始化otto文字转语音引擎")

	log.Println("加载配置文件")
	loadConfig()
	log.Println("配置文件加载完成")

	log.Println("加载otto词典")
	var err error
	dict, err = getDictionary()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("otto词典加载完成")

	log.Println("otto文字转语音引擎初始化完成")
}
