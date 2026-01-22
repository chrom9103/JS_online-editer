package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ExecuteRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language"`
	ClientID string `json:"clientId"`
}

type ExecuteResponse struct {
	Success bool     `json:"success"`
	Output  []Output `json:"output"`
	Error   string   `json:"error,omitempty"`
}

type Output struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SandboxRequest struct {
	Code    string `json:"code"`
	Timeout int    `json:"timeout"`
}

type SandboxResponse struct {
	Success bool     `json:"success"`
	Output  []Output `json:"output"`
	Error   string   `json:"error,omitempty"`
}

// saveRunToFile saves the code to a file under ./tmp/runs
func saveRunToFile(code string, clientID string, ipHash string) error {
	dir := "./tmp/runs"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// サニタイズ: ファイル名に使えるように clientID の非英数字を '_' に置換
	sanitizeRe := regexp.MustCompile(`[^A-Za-z0-9_-]`)
	safeID := sanitizeRe.ReplaceAllString(clientID, "_")

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	max := int64(-1)
	// ファイル名形式: run-<safeID>-<4hex>.js
	re := regexp.MustCompile(`^run-` + regexp.QuoteMeta(safeID) + `-([0-9a-fA-F]{4})\.js$`)
	for _, e := range entries {
		name := e.Name()
		m := re.FindStringSubmatch(name)
		if len(m) == 2 {
			if n, err := strconv.ParseInt(m[1], 16, 64); err == nil {
				if n > max {
					max = n
				}
			}
		}
	}

	next := max + 1
	if next < 0 {
		next = 0
	}
	filename := fmt.Sprintf("run-%s-%04x.js", safeID, next)
	full := filepath.Join(dir, filename)

	f, err := os.Create(full)
	if err != nil {
		return err
	}
	defer f.Close()

	jst := time.FixedZone("JST", 9*60*60)
	header := fmt.Sprintf("// ClientID: %s\n// IPHash: %s\n// Start: %s\n\n", clientID, ipHash, time.Now().In(jst).Format(time.RFC3339))
	if _, err := f.WriteString(header + code); err != nil {
		return err
	}

	return nil
}

func ExecuteCode(c *gin.Context) {
	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExecuteResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// 保存先を作成し、コードを保存する
	userIP := c.ClientIP()
	h := sha256.Sum256([]byte(userIP))
	ipHash := hex.EncodeToString(h[:])
	if err := saveRunToFile(req.Code, req.ClientID, ipHash); err != nil {
		c.JSON(http.StatusInternalServerError, ExecuteResponse{
			Success: false,
			Error:   "Failed to save run: " + err.Error(),
		})
		return
	}

	sandboxURL := os.Getenv("SANDBOX_SERVICE_URL")
	if sandboxURL == "" {
		sandboxURL = "http://localhost:3000"
	}

	sandboxReq := SandboxRequest{
		Code:    req.Code,
		Timeout: 10000, // 10秒（ミリ秒）
	}

	jsonData, err := json.Marshal(sandboxReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ExecuteResponse{
			Success: false,
			Error:   "Failed to marshal request: " + err.Error(),
		})
		return
	}

	// HTTPクライアント（タイムアウト設定）
	client := &http.Client{
		Timeout: 15 * time.Second, // サンドボックスより少し長めに設定
	}

	resp, err := client.Post(sandboxURL+"/execute", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ExecuteResponse{
			Success: false,
			Error:   "Sandbox service unavailable: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ExecuteResponse{
			Success: false,
			Error:   "Failed to read response: " + err.Error(),
		})
		return
	}

	var sandboxResp SandboxResponse
	if err := json.Unmarshal(body, &sandboxResp); err != nil {
		c.JSON(http.StatusInternalServerError, ExecuteResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse sandbox response: %v, body: %s", err, string(body)),
		})
		return
	}

	c.JSON(http.StatusOK, ExecuteResponse{
		Success: sandboxResp.Success,
		Output:  sandboxResp.Output,
		Error:   sandboxResp.Error,
	})
}
