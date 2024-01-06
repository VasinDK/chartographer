// test - интеграционные тесты. Кейсы можно расширять, проверки улучшать
package test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
	"this_module/internal/usecase"
)

const (
	host = "http://localhost:8080"
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
			Name:         "Normal",
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
			ImageAddres:  "./sample.bmp",
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
			ImageAddres:  "./sample.bmp",
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

func Create(width, height, ContentType string, Body io.Reader) (*http.Response, error) {
	return http.Post(host+"/v1/chartas/?width="+width+"&height="+height+"", ContentType, Body)
}

func deleteChart(body string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", host+"/v1/chartas/"+body+"/", nil)
	return client.Do(req)
}
