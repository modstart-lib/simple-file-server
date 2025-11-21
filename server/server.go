package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"simple-file-server/global"
	"simple-file-server/lib/common"
	"simple-file-server/lib/cron"
	"simple-file-server/lib/files"
	"simple-file-server/lib/response"
	"strconv"
	"strings"
)

func Start() {
	cron.Run()
	StartApi()
}

func StartApi() {

	if global.CONFIG.Debug {
		log.SetLevel(log.DebugLevel)
		log.Info("Debug mode enabled")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	log.Info("Server starting")

	r.GET("_admin/ping", ActionPing)
	r.POST("_admin/upload/multipart_init", ActionUploadMultipartInit)
	r.POST("_admin/upload/multipart_upload", ActionUploadMultipartUpload)
	r.POST("_admin/upload/multipart_end", ActionUploadMultipartEnd)
	r.POST("_admin/upload/abort", ActionUploadAbort)
	r.POST("_admin/upload", ActionUpload)
	r.POST("_admin/move", ActionMove)
	r.POST("_admin/delete", ActionDelete)
	r.POST("_admin/has", ActionHas)
	r.POST("_admin/size", ActionSize)
	r.POST("_admin/get", ActionGet)

	r.NoRoute(ActionServeFile)

	// create TempDir
	files.EnsureDir(global.CONFIG.DataDir, "0755")
	files.EnsureDir(global.CONFIG.TempDir, "0755")
	files.EnsureDir(global.CONFIG.TempDir+"/MultiPart", "0755")

	log.Info("Server listening on port ", global.CONFIG.Port)
	r.Run(fmt.Sprintf(":%d", global.CONFIG.Port))
}

func ActionPing(c *gin.Context) {
	response.GenerateSuccess(c, "ok")
}

func checkAdminToken(c *gin.Context) bool {
	token := c.GetHeader("admin-api-token")
	if token != global.CONFIG.ApiToken {
		response.GenerateError(c, "Invalid token")
		return false
	}
	return true
}

type FileDownloadInfo struct {
	Id    string `json:"id"`
	Url   string `json:"url"`
	Name  string `json:"name"`
	Md5   string `json:"md5"`
	Speed int    `json:"speed"`
}

type MultipartMeta struct {
	UploadID   string `json:"uploadId"`
	FilePath   string `json:"filePath"`
	TotalParts int    `json:"totalParts"`
	TotalSize  int64  `json:"totalSize"`
}

func ActionUploadMultipartInit(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	var req struct {
		FilePath   string `json:"filePath"`
		TotalParts int    `json:"totalParts"`
		TotalSize  int64  `json:"totalSize"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GenerateError(c, "Invalid request")
		return
	}
	uploadID := common.RandomString(32)
	meta := MultipartMeta{
		UploadID:   uploadID,
		FilePath:   req.FilePath,
		TotalParts: req.TotalParts,
		TotalSize:  req.TotalSize,
	}
	dir := global.CONFIG.TempDir + "/MultiPart/" + uploadID
	files.EnsureDir(dir, "0755")
	metaFile := dir + "/meta.json"
	data, _ := json.Marshal(meta)
	os.WriteFile(metaFile, data, 0644)
	response.GenerateSuccessWithData(c, "ok", gin.H{"uploadId": uploadID})
}

func ActionUploadMultipartUpload(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	uploadID := c.PostForm("uploadId")
	partNumberStr := c.PostForm("partNumber")
	partNumber, err := strconv.Atoi(partNumberStr)
	if err != nil {
		response.GenerateError(c, "Invalid partNumber")
		return
	}
	// check uploadID and partNumber
	if uploadID == "" || partNumber < 0 {
		response.GenerateError(c, "Invalid uploadId or partNumber")
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.GenerateError(c, "Invalid file")
		return
	}
	defer file.Close()
	dir := global.CONFIG.TempDir + "/MultiPart/" + uploadID
	if !files.FileExists(dir) {
		response.GenerateError(c, "UploadIDNotFound")
		return
	}
	partFile := dir + "/part" + strconv.Itoa(partNumber)
	out, err := os.Create(partFile)
	if err != nil {
		response.GenerateError(c, "Failed to create part file")
		return
	}
	defer out.Close()
	io.Copy(out, file)
	response.GenerateSuccess(c, "ok")
}

func ActionUploadMultipartEnd(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	var req struct {
		UploadID string `json:"uploadId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GenerateError(c, "Invalid request")
		return
	}
	// check uploadID and partNumber
	if req.UploadID == "" {
		response.GenerateError(c, "Invalid uploadId or partNumber")
		return
	}
	dir := global.CONFIG.TempDir + "/MultiPart/" + req.UploadID
	metaFile := dir + "/meta.json"
	data, err := os.ReadFile(metaFile)
	if err != nil {
		response.GenerateError(c, "Meta not found")
		return
	}
	var meta MultipartMeta
	json.Unmarshal(data, &meta)
	finalFile := filepath.Join(global.CONFIG.DataDir, meta.FilePath)
	files.EnsureDir(filepath.Dir(finalFile), "0755")
	out, err := os.Create(finalFile)
	if err != nil {
		response.GenerateError(c, "Failed to create final file")
		return
	}
	defer out.Close()
	for i := 1; i <= meta.TotalParts; i++ {
		partFile := dir + "/part" + strconv.Itoa(i)
		part, err := os.Open(partFile)
		if err != nil {
			response.GenerateError(c, "Part file missing")
			return
		}
		io.Copy(out, part)
		part.Close()
	}
	files.DeleteDir(dir)
	response.GenerateSuccessWithData(c, "ok", gin.H{
		"filePath": meta.FilePath,
	})
}

func ActionUpload(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.GenerateError(c, "Invalid file")
		return
	}
	defer file.Close()
	filePath := c.PostForm("filePath")
	if filePath == "" {
		response.GenerateError(c, "filePath is required")
		return
	}
	filePath = filepath.Join(global.CONFIG.DataDir, filePath)
	files.EnsureDir(filepath.Dir(filePath), "0755")
	out, err := os.Create(filePath)
	if err != nil {
		response.GenerateError(c, "Failed to create file")
		return
	}
	defer out.Close()
	io.Copy(out, file)
	response.GenerateSuccessWithData(c, "ok", gin.H{
		"filePath": filePath,
	})
}

func ActionServeFile(c *gin.Context) {
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/_admin/") {
		c.AbortWithStatus(404)
		return
	}
	fullPath := filepath.Join(global.CONFIG.DataDir, strings.TrimPrefix(path, "/"))
	if !files.FileExists(fullPath) {
		c.AbortWithStatus(404)
		return
	}
	info, err := os.Stat(fullPath)
	if err != nil || info.IsDir() {
		c.AbortWithStatus(404)
		return
	}
	c.File(fullPath)
}

func ActionMove(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	var req struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GenerateError(c, "Invalid request")
		return
	}
	fromPath := filepath.Join(global.CONFIG.DataDir, req.From)
	toPath := filepath.Join(global.CONFIG.DataDir, req.To)
	if !files.FileExists(fromPath) {
		response.GenerateError(c, "Source file not found")
		return
	}
	files.EnsureDir(filepath.Dir(toPath), "0755")
	err := os.Rename(fromPath, toPath)
	if err != nil {
		response.GenerateError(c, "Failed to move file")
		return
	}
	response.GenerateSuccess(c, "ok")
}

func ActionDelete(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GenerateError(c, "Invalid request")
		return
	}
	fullPath := filepath.Join(global.CONFIG.DataDir, req.Path)
	if !files.FileExists(fullPath) {
		response.GenerateError(c, "File not found")
		return
	}
	err := os.Remove(fullPath)
	if err != nil {
		response.GenerateError(c, "Failed to delete file")
		return
	}
	response.GenerateSuccess(c, "ok")
}

func ActionHas(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GenerateError(c, "Invalid request")
		return
	}
	fullPath := filepath.Join(global.CONFIG.DataDir, req.Path)
	exists := files.FileExists(fullPath)
	response.GenerateSuccessWithData(c, "ok", exists)
}

func ActionSize(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GenerateError(c, "Invalid request")
		return
	}
	fullPath := filepath.Join(global.CONFIG.DataDir, req.Path)
	if !files.FileExists(fullPath) {
		response.GenerateError(c, "File not found")
		return
	}
	info, err := os.Stat(fullPath)
	if err != nil {
		response.GenerateError(c, "Failed to get file info")
		return
	}
	response.GenerateSuccessWithData(c, "ok", gin.H{
		"size": info.Size(),
	})
}

func ActionGet(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(404, err)
		return
	}
	fullPath := filepath.Join(global.CONFIG.DataDir, req.Path)
	if !files.FileExists(fullPath) {
		c.AbortWithStatus(404)
		return
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Data(200, "application/octet-stream", data)
}

func ActionUploadAbort(c *gin.Context) {
	if !checkAdminToken(c) {
		return
	}
	var req struct {
		UploadID string `json:"uploadId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GenerateError(c, "Invalid request")
		return
	}
	if req.UploadID == "" {
		response.GenerateError(c, "Invalid uploadId")
		return
	}
	dir := global.CONFIG.TempDir + "/MultiPart/" + req.UploadID
	if !files.FileExists(dir) {
		response.GenerateError(c, "UploadIDNotFound")
		return
	}
	files.DeleteDir(dir)
	response.GenerateSuccess(c, "ok")
}
