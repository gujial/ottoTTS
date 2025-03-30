package main

import (
	"fmt"
	"github.com/gujial/ottoTTS"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ottoTTScli <str>")
		return
	}

	ottoTTS.InitializeTTS()
	str := os.Args[1]
	b, err := ottoTTS.Speech(str, true)
	if err != nil {
		panic(err)
	}

	fmt.Println("写入文件中")
	file, err := os.OpenFile("otto.wav", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	count, err := file.Write(b)
	if err != nil {
		panic(err)
	}
	fmt.Println("共写入", count, "bytes")
}
