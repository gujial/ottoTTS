package wavHandler

import "os"

type Slice struct {
	Category string
	Content  string
}

func readSingleWordFile(word string) ([]byte, error) {
	data, err := os.ReadFile("./assets/sounds/single/" + word + ".wav")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func readMultiWordFile(word string) ([]byte, error) {
	data, err := os.ReadFile("./assets/sounds/multiple/" + word + ".wav")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetSpeech(slices []Slice) ([]byte, error) {
	var speech []byte
	for _, slice := range slices {
		switch slice.Category {
		case "expressions":
			speechSlice, err := readMultiWordFile(slice.Content)
			if err != nil {
				return nil, err
			}
			speech = append(speech, speechSlice...)
			break
		default:
			speechSlice, err := readSingleWordFile(slice.Content)
			if err != nil {
				return nil, err
			}
			speech = append(speech, speechSlice...)
			break
		}
	}
	return speech, nil
}
