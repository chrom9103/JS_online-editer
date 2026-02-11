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
