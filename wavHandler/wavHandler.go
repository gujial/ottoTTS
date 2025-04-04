package wavHandler

import (
	"github.com/gujial/ottoTTS/wav"
	"os"
)

type Slice struct {
	Category string
	Content  string
}

func sliceToWav(slice Slice) (*wav.WAV, error) {
	var filePath string
	switch slice.Category {
	case "others":
		return nil, nil

	default:
		filePath = "./assets/sounds/" + slice.Content + ".wav"
	}

	matchFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	matchWav, err := wav.ReadWAV(matchFile)
	if err != nil {
		return nil, err
	}

	return matchWav, nil
}

// GetSpeech 生成合成的 WAV 音频
func GetSpeech(slices []Slice) ([]byte, error) {
	var matchWavs []*wav.WAV
	for _, slice := range slices {
		if len(slice.Content) == 0 {
			continue
		}

		matchWav, err := sliceToWav(slice)
		if err != nil {
			return nil, err
		}

		if matchWav != nil {
			matchWavs = append(matchWavs, matchWav)
		}
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
