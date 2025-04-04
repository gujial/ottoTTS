package wav

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type WAV struct {
	NumChannels   uint16
	SampleRate    uint32
	BitsPerSample uint16
	Data          []byte
}

func ReadWAV(data []byte) (*WAV, error) {
	if len(data) < 44 {
		return nil, errors.New("invalid WAV data")
	}

	numChannels := binary.LittleEndian.Uint16(data[22:24])
	sampleRate := binary.LittleEndian.Uint32(data[24:28])
	bitsPerSample := binary.LittleEndian.Uint16(data[34:36])

	// 查找 "data" 子块
	dataStart := bytes.Index(data, []byte("data"))
	if dataStart == -1 || dataStart+8 >= len(data) {
		return nil, errors.New("invalid data chunk")
	}

	// 获取 dataSize 并确保数据长度正确
	dataSize := binary.LittleEndian.Uint32(data[dataStart+4 : dataStart+8])
	if uint32(len(data)-(dataStart+8)) < dataSize {
		return nil, errors.New("invalid WAV data size")
	}

	return &WAV{
		NumChannels:   numChannels,
		SampleRate:    sampleRate,
		BitsPerSample: bitsPerSample,
		Data:          data[dataStart+8 : dataStart+8+int(dataSize)], // 正确偏移
	}, nil
}

func WriteWAV(wav *WAV) ([]byte, error) {
	buffer := new(bytes.Buffer)

	header := make([]byte, 44)
	copy(header[:4], "RIFF")
	binary.LittleEndian.PutUint32(header[4:8], uint32(36+len(wav.Data)))
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16)
	binary.LittleEndian.PutUint16(header[20:22], 1) // PCM 格式
	binary.LittleEndian.PutUint16(header[22:24], wav.NumChannels)
	binary.LittleEndian.PutUint32(header[24:28], wav.SampleRate)
	byteRate := wav.SampleRate * uint32(wav.NumChannels) * uint32(wav.BitsPerSample) / 8
	binary.LittleEndian.PutUint32(header[28:32], byteRate)
	blockAlign := wav.NumChannels * wav.BitsPerSample / 8
	binary.LittleEndian.PutUint16(header[32:34], blockAlign)
	binary.LittleEndian.PutUint16(header[34:36], wav.BitsPerSample)
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], uint32(len(wav.Data)))

	buffer.Write(header)
	buffer.Write(wav.Data)

	return buffer.Bytes(), nil
}

// SilentWAV 生成一段指定持续时间的静音 WAV 音频
func SilentWAV(numChannels uint16, sampleRate uint32, bitsPerSample uint16, durationSeconds float64) (*WAV, error) {
	if bitsPerSample%8 != 0 {
		return nil, errors.New("bitsPerSample must be a multiple of 8")
	}

	bytesPerSample := bitsPerSample / 8
	numSamples := uint32(float64(sampleRate) * durationSeconds)
	totalBytes := int(numSamples) * int(numChannels) * int(bytesPerSample)

	silence := make([]byte, totalBytes) // 默认为 0，表示静音

	return &WAV{
		NumChannels:   numChannels,
		SampleRate:    sampleRate,
		BitsPerSample: bitsPerSample,
		Data:          silence,
	}, nil
}

// ConcatenateWAVs 连接多个 WAV 文件
func ConcatenateWAVs(wavs []*WAV) (*WAV, error) {
	if len(wavs) == 0 {
		return nil, errors.New("no WAV files provided")
	}

	base := wavs[0]
	for _, w := range wavs[1:] {
		if w.NumChannels != base.NumChannels || w.SampleRate != base.SampleRate || w.BitsPerSample != base.BitsPerSample {
			return nil, errors.New("all WAV files must have the same format")
		}
		base.Data = append(base.Data, w.Data...)
	}

	return base, nil
}
