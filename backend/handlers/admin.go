package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func ListRuns(c *gin.Context) {
	dir := os.Getenv("RUNS_DIR")
	if dir == "" {
		dir = "./tmp/runs"
	}

	if !filepath.IsAbs(dir) {
		if wd, err := os.Getwd(); err == nil {
			dir = filepath.Join(wd, dir)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var runs []gin.H
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			continue
		}
		runs = append(runs, gin.H{
			"name":    e.Name(),
			"size":    fi.Size(),
			"modTime": fi.ModTime().In(time.FixedZone("JST", 9*60*60)).Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, runs)
}

func GetRunFile(c *gin.Context) {
	name := c.Param("name")
	// reject paths containing separators
	if name == "" || name != filepath.Base(name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}

	dir := os.Getenv("RUNS_DIR")
	if dir == "" {
		dir = "./tmp/runs"
	}
	if !filepath.IsAbs(dir) {
		if wd, err := os.Getwd(); err == nil {
			dir = filepath.Join(wd, dir)
		}
	}

	full := filepath.Join(dir, name)
	if _, err := os.Stat(full); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.File(full)
}

// DeleteRunFiles は指定されたファイルを一括削除する
func DeleteRunFiles(c *gin.Context) {
	var req struct {
		Files []string `json:"files" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: files array required"})
		return
	}

	if len(req.Files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files specified"})
		return
	}

	dir := os.Getenv("RUNS_DIR")
	if dir == "" {
		dir = "./tmp/runs"
	}
	if !filepath.IsAbs(dir) {
		if wd, err := os.Getwd(); err == nil {
			dir = filepath.Join(wd, dir)
		}
	}

	var deleted []string
	var errors []gin.H

	for _, name := range req.Files {
		// パストラバーサル防止
		if name == "" || name != filepath.Base(name) {
			errors = append(errors, gin.H{"name": name, "error": "invalid filename"})
			continue
		}

		full := filepath.Join(dir, name)
		if _, err := os.Stat(full); err != nil {
			errors = append(errors, gin.H{"name": name, "error": "file not found"})
			continue
		}

		if err := os.Remove(full); err != nil {
			errors = append(errors, gin.H{"name": name, "error": err.Error()})
			continue
		}

		deleted = append(deleted, name)
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": deleted,
		"errors":  errors,
	})
}
