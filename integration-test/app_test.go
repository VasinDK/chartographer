// test - интеграционные тесты. Кейсы можно расширять, проверки улучшать
package integrationtest

import (
	"io"
	"net/http"
)

const (
	host = "http://localhost:8080"
)

func Create(width, height, ContentType string, Body io.Reader) (*http.Response, error) {
	return http.Post(host+"/v1/chartas/?width="+width+"&height="+height+"", ContentType, Body)
}

func deleteChart(body string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", host+"/v1/chartas/"+body+"/", nil)
	return client.Do(req)
}
