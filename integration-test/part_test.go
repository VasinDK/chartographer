// test - интеграционные тесты. Кейсы можно расширять, проверки улучшать
package integrationtest

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestApp_part(t *testing.T) {
	cases := []struct {
		Name         string
		Width        string
		Height       string
		X            string
		Y            string
		ImageId      string
		ErrorPresent bool
		ErrorText    string
		StatusCode   int
	}{
		{
			Name:         "Normal13 ",
			Width:        "200",
			Height:       "200",
			X:            "0",
			Y:            "0",
			ImageId:      "test1",
			ErrorPresent: false,
			ErrorText:    "",
		},
		{
			Name:         "Normal2",
			Width:        "500",
			Height:       "100",
			X:            "10",
			Y:            "10",
			ImageId:      "test1",
			ErrorPresent: false,
			ErrorText:    "",
		},
	}

	for _, cs := range cases {
		resp, err := http.Get(host + "/v1/chartas/" + cs.ImageId + "/?x=" + cs.X + "&y=" + cs.Y + "&width=" + cs.Width + "&height=" + cs.Height + "")
		if err != nil {
			t.Errorf("case.Name: %v, http.Get: %v", cs.Name, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			continue
		}

		bodyByte, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("case.Name: %v, io.ReadAll: %v", cs.Name, err.Error())
		}
		body := string(bodyByte)

		if cs.ErrorPresent == true && strings.TrimSpace(body) == cs.ErrorText {
			continue
		}

		t.Errorf("case.Name: %v. body: %v", cs.Name, body)
	}
}
