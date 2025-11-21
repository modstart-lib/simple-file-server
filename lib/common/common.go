package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"simple-file-server/lib/defs"
	"simple-file-server/lib/result"
	"strings"
	"time"
)

func LoadTimeZone() string {
	loc := time.Now().Location()
	if _, err := time.LoadLocation(loc.String()); err != nil {
		return "Asia/Shanghai"
	}
	return loc.String()
}

const (
	b  = uint64(1)
	kb = 1024 * b
	mb = 1024 * kb
	gb = 1024 * mb
)

func FormatBytes(bytes uint64) string {
	switch {
	case bytes < kb:
		return fmt.Sprintf("%dB", bytes)
	case bytes < mb:
		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(kb))
	case bytes < gb:
		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(mb))
	default:
		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(gb))
	}
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func Request(url string, data interface{}) defs.Response {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return result.GenerateError("json marshal error")
	}
	var fullUrl string
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		fullUrl = url
	} else {
		return result.GenerateError("url must start with http:// or https://")
	}
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return result.GenerateError("http request error")
	}
	//req.Header.Set("Api-Uuid", global.CONFIG.UUID)
	//req.Header.Set("Api-Device", platform.Name()+"/"+platform.Arch())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "simple-file-server")
	//if global.CONFIG.ApiToken != "" {
	//	req.Header.Set("Api-Token", global.CONFIG.ApiToken)
	//}
	//if global.CONFIG.LauncherKey != "" {
	//	req.Header.Set("Launcher-Key", global.CONFIG.LauncherKey)
	//}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result.GenerateError("http client error")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return result.GenerateError("http status error: " + string(resp.StatusCode) + " " + string(body))
	}
	var res defs.Response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return result.GenerateError("json decode error : " + string(body))
	}
	return res
}
