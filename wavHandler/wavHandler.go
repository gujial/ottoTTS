package wavHandler

import (
	"io"
	"os"

	"github.com/moutend/go-wav"
)

type Slice struct {
	Category string
	Content  string
}

// 读取 WAV 文件并返回音频数据和格式信息
func readWAVFile(filepath string) ([]byte, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// 生成合成的 WAV 音频
func GetSpeech(slices []Slice) ([]byte, error) {
	result := &wav.File{}
	for _, slice := range slices {
		var filePath string

		switch slice.Category {
		case "expressions":
			filePath = "./assets/sounds/multiple/" + slice.Content + ".wav"
		default:
			filePath = "./assets/sounds/single/" + slice.Content + ".wav"
		}

		matchFile, err := readWAVFile(filePath)
		if err != nil {
			return nil, err
		}

		matchWav := &wav.File{}
		wav.Unmarshal(matchFile, matchWav)
		result, _ = wav.New(matchWav.SamplesPerSec(), matchWav.BitsPerSample(), matchWav.Channels())
		io.Copy(result, matchWav)
	}

	return wav.Marshal(result)
}
