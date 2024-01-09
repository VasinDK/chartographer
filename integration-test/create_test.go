// test - интеграционные тесты. Кейсы можно расширять, проверки улучшать
package integrationtest

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"this_module/internal/usecase"
)

func TestApp_create(t *testing.T) {
	cases := []struct {
		Name         string
		Width        string
		Height       string
		ContentType  string
		Body         io.Reader
		ErrorPresent bool
		ErrorText    string
	}{
		{
			Name:         "Normal ",
			Width:        "250",
			Height:       "223",
			ContentType:  "",
			Body:         nil,
			ErrorPresent: false,
			ErrorText:    "",
		},
		{
			Name:         "WidthMax",
			Width:        "25000",
			Height:       "223",
			ContentType:  "",
			Body:         nil,
			ErrorPresent: true,
			ErrorText:    "The image dimensions exceed the allowed ones",
		},
	}

	for _, cs := range cases {
		resp, err := Create(cs.Width, cs.Height, cs.ContentType, cs.Body)
		if err != nil {
			t.Errorf("Create: %v", err.Error())
		}
		defer resp.Body.Close()

		bodyByte, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Create.io.ReadAll: %v", err.Error())
		}
		body := string(bodyByte)
		if resp.StatusCode == http.StatusCreated && len(body) == usecase.LenIdImg {
			deleteChart(body)
			continue
		}
		if cs.ErrorPresent == true && strings.TrimSpace(body) == cs.ErrorText {
			continue
		}

		t.Errorf("case.Name: %v. body: %v", cs.Name, body)

		deleteChart(body)
	}
}
