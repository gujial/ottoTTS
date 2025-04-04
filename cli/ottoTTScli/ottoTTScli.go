package main

import (
	"fmt"
	"github.com/gujial/ottoTTS"
	"log"
	"os"
)

func saveWav(str string) {
	ottoTTS.InitializeTTS()
	b, err := ottoTTS.Speech(str)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("写入文件中")
	file, err := os.OpenFile("otto.wav", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	count, err := file.Write(b)
	if err != nil {
		panic(err)
	}
	log.Println("共写入", count, "bytes")
}

func main() {
	var str string

	if len(os.Args) < 2 {
		fmt.Println("输入需要转换的字符串")
		_, err := fmt.Scanf("%s", &str)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		str = os.Args[1]
	}

	saveWav(str)
}
