package api

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"ghost-images/config"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type GhostError struct {
	Message string
}

func (e *GhostError) Error() string {
	return e.Message
}

var (
	VALID_HTTP_METHODS = []string{"GET", "POST", "PUT", "DELETE"}
)

func MakeGhostRequest(endpoint string, headers map[string]string, ctx context.Context, isResource bool, httpMethod string, jsonData map[string]interface{}) (map[string]interface{}, error) {
	httpMethod = strings.ToUpper(httpMethod)
	if !isValidHTTPMethod(httpMethod) {
		return nil, fmt.Errorf("invalid HTTP method: %s", httpMethod)
	}

	url := constructURL(endpoint)
	req, err := createRequest(httpMethod, url, jsonData)
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("Error creating request: %v", err)}
	}

	setHeaders(req, headers)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("HTTP error accessing Ghost API: %v", err)}
	}
	defer resp.Body.Close()

	if err := handleResponseStatus(httpMethod, resp); err != nil {
		return nil, err
	}

	return parseResponseBody(resp)
}

func MakeGhostMultipartRequest(endpoint string, headers map[string]string, ctx context.Context, isResource bool, httpMethod string, multipartLocalPath string) (map[string]interface{}, error) {

	file, err := os.Open(multipartLocalPath)
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("파일 열기 오류: %v", err)}
	}
	defer file.Close()

	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)

	part, err := writer.CreateFormFile("file", filepath.Base(multipartLocalPath))
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("폼 파일 생성 오류: %v", err)}
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("파일 내용 복사 오류: %v", err)}
	}

	writer.WriteField("ref", multipartLocalPath)

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, httpMethod, constructURL(endpoint), &reqBody)
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("요청 생성 오류: %v", err)}
	}

	setHeaders(req, headers)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("Ghost API 접근 HTTP 오류: %v", err)}
	}
	defer resp.Body.Close()

	if err := handleResponseStatus(httpMethod, resp); err != nil {
		return nil, err
	}

	return parseResponseBody(resp)
}

func MakeGhostMultipartRequestCurl(endpoint string, headers map[string]string, ctx context.Context, isResource bool, httpMethod string, multipartLocalPath string) (map[string]interface{}, error) {
	token, err := getToken(config.Config.GHOST_STAFF_API_KEY)
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("토큰 생성 오류: %v", err)}
	}

	// 백틱을 사용하여 여러 줄의 문자열을 올바르게 포맷팅
	curlCommand := fmt.Sprintf(`curl -X %s -H "Content-Type: multipart/form-data" -H "Authorization: Ghost %s" -H "Accept-Version: v5.109" -F "file=@%s" -F "ref=%s" %s`,
		httpMethod,
		token,
		multipartLocalPath,
		multipartLocalPath,
		constructURL(endpoint))

	// 명령어 실행
	// windows
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", curlCommand)
	} else {
		cmd = exec.Command("sh", "-c", curlCommand)
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("명령어 실행 오류: %v", err)}
	}

	return map[string]interface{}{
		"output": string(output),
	}, nil
}

func isValidHTTPMethod(method string) bool {
	for _, validMethod := range VALID_HTTP_METHODS {
		if validMethod == method {
			return true
		}
	}
	return false
}

func constructURL(endpoint string) string {
	baseURL := strings.TrimRight(config.Config.GHOST_API_URL, "/") + "/ghost/api/admin"
	return baseURL + "/" + strings.Trim(endpoint, "/")
}

func createRequest(method, url string, jsonData map[string]interface{}) (*http.Request, error) {
	if method == "POST" || method == "PUT" {
		jsonBytes, _ := json.Marshal(jsonData)
		return http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	}
	return http.NewRequest(method, url, nil)
}

func setHeaders(req *http.Request, headers map[string]string) error {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if headers["Accept-Version"] == "" {
		req.Header.Set("Accept-Version", "v5.109")
	}
	if headers["Authorization"] == "" {
		token, err := getToken(config.Config.GHOST_STAFF_API_KEY)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Ghost "+token)
	}
	return nil
}

func handleResponseStatus(method string, resp *http.Response) error {
	if method == "DELETE" && resp.StatusCode == 204 {
		return nil
	}
	var body []byte
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ = io.ReadAll(resp.Body)
		return &GhostError{Message: fmt.Sprintf("HTTP error accessing Ghost API: Status %d, %s", resp.StatusCode, string(body))}
	}
	return nil
}

func parseResponseBody(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("Error reading response body: %v", err)}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &GhostError{Message: fmt.Sprintf("Error parsing JSON response: %v", err)}
	}

	return result, nil
}

func getToken(staffAPIKey string) (string, error) {
	token, err := generateToken(staffAPIKey, "/admin/")
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateToken(staffAPIKey string, audience string) (string, error) {
	if audience == "" {
		audience = "/admin/"
	}

	parts := split(staffAPIKey, ":")
	if len(parts) != 2 {
		return "", errors.New("STAFF_API_KEY must be in the format 'id:secret'")
	}

	keyID, secret := parts[0], parts[1]
	if keyID == "" || secret == "" {
		return "", errors.New("both key ID and secret are required")
	}

	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return "", errors.New("invalid secret format - must be hexadecimal")
	}

	now := time.Now().UTC()
	exp := now.Add(5 * time.Minute)

	payload := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": exp.Unix(),
		"aud": audience,
		"sub": keyID,
		"typ": "ghost-admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token.Header["kid"] = keyID

	tokenString, err := token.SignedString(secretBytes)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func split(s, sep string) []string {
	var parts []string
	for len(s) > 0 {
		idx := indexOf(s, sep)
		if idx < 0 {
			parts = append(parts, s)
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	return parts
}

func indexOf(s, sep string) int {
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}
