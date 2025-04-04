package wavHandler

import (
	"github.com/gujial/ottoTTS/wav"
	"os"
)

type Slice struct {
	Category string
	Content  string
}

// GetSpeech 生成合成的 WAV 音频
func GetSpeech(slices []Slice) ([]byte, error) {
	var matchWavs []*wav.WAV
	for _, slice := range slices {
		var filePath string

		if len(slice.Content) == 0 {
			continue
		}

		switch slice.Category {
		case "expressions":
			filePath = "./assets/sounds/multiple/" + slice.Content + ".wav"
		case "others":
			continue

		default:
			filePath = "./assets/sounds/single/" + slice.Content + ".wav"
		}

		matchFile, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		matchWav, err := wav.ReadWAV(matchFile)
		if err != nil {
			return nil, err
		}

		matchWavs = append(matchWavs, matchWav)
	}

	ConcatenatedWav, err := wav.ConcatenateWAVs(matchWavs)
	if err != nil {
		return nil, err
	}

	resultWav, err := wav.WriteWAV(ConcatenatedWav)
	if err != nil {
		return nil, err
	}
	return resultWav, nil
}
