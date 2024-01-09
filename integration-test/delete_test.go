// test - интеграционные тесты. Кейсы можно расширять, проверки улучшать
package integrationtest

import (
	"io"
	"net/http"
	"testing"
)

func TestApp_Delete(t *testing.T) {
	cases := []struct {
		Name         string
		ImageId      string
		ErrorPresent bool
		StatusCode   int
	}{
		{
			Name:         "Normal1  ",
			ImageId:      "",
			ErrorPresent: false,
			StatusCode:   0,
		},
		{
			Name:         "Error400",
			ImageId:      "test0000",
			ErrorPresent: true,
			StatusCode:   http.StatusBadRequest,
		},
	}

	for _, cs := range cases {
		chartBody := cs.ImageId
		if cs.ImageId == "" {
			chartResp, _ := Create("350", "323", "", nil)
			defer chartResp.Body.Close()
			chartBodyByte, _ := io.ReadAll(chartResp.Body)
			chartBody = string(chartBodyByte)
		}

		resp, err := deleteChart(chartBody)
		if err != nil {
			t.Errorf("case.Name: %v, http.DELETE: %v", cs.Name, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			continue
		}

		if cs.ErrorPresent == true && resp.StatusCode == cs.StatusCode {
			deleteChart(chartBody)
			continue
		}

		t.Errorf("case.Name: %v. body: %v", cs.Name, "The error does not match")
	}
}
