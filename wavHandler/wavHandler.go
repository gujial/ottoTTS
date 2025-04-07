package wavHandler

import (
	"github.com/gujial/ottoTTS/wav"
	"os"
)

type Slice struct {
	Category string
	Content  string
}

func generateApproxNames(content string) []string {
	var result []string

	// 简单规则：逐步去掉最后一个字符（退化为前缀匹配）
	for i := len(content) - 1; i > 0; i-- {
		result = append(result, content[:i])
	}

	return result
}

func sliceToWav(slice Slice) (*wav.WAV, error) {
	if slice.Category == "others" {
		return nil, nil
	}

	basePath := "./assets/sounds/"
	filename := slice.Content + ".wav"
	filePath := basePath + filename

	// 先尝试直接读取
	if data, err := tryReadWav(filePath); err == nil {
		return data, nil
	}

	// 近似匹配尝试
	approxNames := generateApproxNames(slice.Content)
	for _, alt := range approxNames {
		altPath := basePath + alt + ".wav"
		if data, err := tryReadWav(altPath); err == nil {
			return data, nil
		}
	}

	return nil, nil // 如果还是没有匹配，返回 nil 用静音
}

func tryReadWav(filePath string) (*wav.WAV, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return wav.ReadWAV(bytes)
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
		} else {
			silentWav, err := wav.SilentWAV(
				matchWavs[0].NumChannels,
				matchWavs[0].SampleRate,
				matchWavs[0].BitsPerSample,
				0.2,
			)
			if err != nil {
				return nil, err
			}
			matchWavs = append(matchWavs, silentWav)
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
