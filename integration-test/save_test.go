// test - интеграционные тесты. Кейсы можно расширять, проверки улучшать
package integrationtest

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestApp_save(t *testing.T) {
	cases := []struct {
		Name         string
		Width        string
		Height       string
		X            string
		Y            string
		ImageAddres  string
		ErrorPresent bool
		ErrorText    string
		StatusCode   int
	}{
		{
			Name:         "normal1",
			Width:        "300",
			Height:       "300",
			X:            "10",
			Y:            "0",
			ImageAddres:  "../docs/img/i/sample.bmp",
			ErrorPresent: false,
			ErrorText:    "",
			StatusCode:   500,
		},
		{
			Name:         "normal2",
			Width:        "500",
			Height:       "500",
			X:            "20",
			Y:            "50",
			ImageAddres:  "../docs/img/i/sample.bmp",
			ErrorPresent: false,
			ErrorText:    "",
			StatusCode:   500,
		},
	}

	for _, cs := range cases {
		chartResp, _ := Create("250", "223", "", nil)
		// defer chartResp.Body.Close()
		chartBodyByte, _ := io.ReadAll(chartResp.Body)
		chartBody := string(chartBodyByte)
		chartResp.Body.Close()

		dstFile, err := os.Open(cs.ImageAddres)
		if err != nil {
			t.Errorf("case.Name: %v, save.os.Open: %v", cs.Name, err.Error())
		}
		defer dstFile.Close()

		requestBody := &bytes.Buffer{}
		writer := multipart.NewWriter(requestBody)

		imagePart, err := writer.CreateFormFile("aaa", "sample.bmp")
		if err != nil {
			t.Errorf("case.Name: %v, save.writer.CreateFormFile: %v", cs.Name, err.Error())
		}

		_, err = io.Copy(imagePart, dstFile)
		if err != nil {
			t.Errorf("case.Name: %v, save.io.Copy: %v", cs.Name, err.Error())
		}

		writer.Close()

		resp, err := http.Post(host+"/v1/chartas/"+chartBody+"/?x="+cs.X+"&y="+cs.Y+"&width="+cs.Width+"&height="+cs.Height+"",
			writer.FormDataContentType(),
			requestBody,
		)
		if err != nil {
			t.Errorf("case.Name: %v, save.http.Post: %v", cs.Name, err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			deleteChart(chartBody)
			continue
		}

		if cs.ErrorPresent && resp.StatusCode == cs.StatusCode {
			deleteChart(chartBody)
			continue
		}

		t.Errorf("case.Name: %v, save.StatusCode: %v", cs.Name, resp.StatusCode)

		deleteChart(chartBody)
	}
}
