package translate

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

type translation struct {
	from, to, text, wantText string
	wantError                error
}

func TestTranslateText(t *testing.T) {
	if err := os.Remove(languagesPath); err != nil && !os.IsNotExist(err) {
		t.Error("unexpected error removing languages.json")
	}
	tc := []translation{
		{"english", "fr", "hello", "salut", nil},
		{"", "french", "hello", "salut", nil},
		{"boulder", "fr", "hello", "", ErrInvalidLang},
		{" ", "boulder", "hello", "", ErrInvalidLang},
	}
	for _, c := range tc {
		t.Run(fmt.Sprintf("f:%s,t:%s", c.from, c.to), func(t *testing.T) {
			text, err := TranslateText(c.from, c.to, c.text)
			if c.wantText != text {
				t.Errorf("text got:%s.want:%s", text, c.wantText)
			}
			if !(errors.Is(err, c.wantError)) {
				t.Errorf("error got:%v.want:%v", err, c.wantError)
			}
		})
	}
	if err := os.Remove(languagesPath); err != nil && !os.IsNotExist(err) {
		t.Error("unexpected error removing languages.json")
	}
}
