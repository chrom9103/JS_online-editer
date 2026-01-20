package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type ExecuteRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language"`
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

func ExecuteCode(c *gin.Context) {
	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExecuteResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
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
