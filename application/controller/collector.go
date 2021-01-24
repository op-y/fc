package controller

import (
    "log"
	"net/http"

	"github.com/op-y/fc/backend/alioss"
	"github.com/op-y/fc/config"
	"github.com/op-y/fc/utils"

	"github.com/gin-gonic/gin"
)

type LogUploadInput struct {
	ShopID    string `form:"shopid" json:"shopid"`
	DeviceID  string `form:"deviceid" json:"deviceid"`
	VerSaaS   string `form:"versaas" json:"versaas"`
	AppName   string `form:"applicationname" json:"applicationname"`
	Timestamp string `form:"timestamp" json:"timestamp"`
	Digest    string `form:"digest" json:"digest"`
}

func (in LogUploadInput) GetMessage(filename, salt string) string {
	return in.ShopID + in.DeviceID + in.AppName + in.Timestamp + filename + salt
}

func (in LogUploadInput) GetKey(prefix, filename string) string {
	return prefix + "/" + in.ShopID + "/" + in.DeviceID + "/" + in.AppName + "/" + in.Timestamp + "/" + filename
}

type DBUploadInput struct {
	GroupID    string `form:"groupid" json:"groupid"`
	ShopID     string `form:"shopid" json:"shopid"`
	BackupDate string `form:"backUpDate" json:"backUpDate"`
	Module     string `form:"module" json:"module"`
	Digest     string `form:"digest" json:"digest"`
}

func (in DBUploadInput) GetMessage(filename, salt string) string {
	return in.GroupID + in.ShopID + in.BackupDate + in.Module + filename + salt
}

func (in DBUploadInput) GetKey(prefix, filename string) string {
	return prefix + "/" + in.GroupID + "/" + in.ShopID + "/" + in.BackupDate + "/" + in.Module + "/" + filename
}

type ResponseFile struct {
	Code     int    `json:"code"`
	Status   string `json:"Status"`
	FileInfo string `json:"FileInfo"`
	Message  string `json:"message"`
}

// upload pos log
func UploadLog(c *gin.Context) {
	var input LogUploadInput
	var msg ResponseFile

	if c.ShouldBind(&input) != nil {
		msg.Code = http.StatusBadRequest
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "failed to bind parameters"
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	header, err := c.FormFile("uploadfile")
	if err != nil {
		msg.Code = http.StatusBadRequest
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "failed to get upload file"
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	if ok := utils.FilterExt(header.Filename, ".zip") || utils.FilterExt(header.Filename, ".log"); !ok {
		msg.Code = http.StatusBadRequest
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "only zip or log file is allowed"
		c.JSON(http.StatusBadRequest, msg)
		return
	}

    // verify
	if config.Cfg.Verification.Enable {
		message := input.GetMessage(header.Filename, config.Cfg.Verification.Salt)
		if !utils.CheckDigest(message, input.Digest, config.Cfg.Verification.Method) {
			msg.Code = http.StatusBadRequest
			msg.Status = "fail"
			msg.FileInfo = ""
			msg.Message = "unexpect digest"
			c.JSON(http.StatusBadRequest, msg)
			return
		}
	}

	file, err := header.Open()
	if err != nil {
		msg.Code = http.StatusInternalServerError
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "file error"
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer file.Close()

	// upload
    var objectKey string
	if config.Cfg.Backend == "oss" {
        objectKey = input.GetKey(config.Cfg.PosLog.Prefix, header.Filename)
        if err := alioss.Client.PutObject(objectKey, file); err != nil {
            msg.Code = http.StatusInternalServerError
            msg.Status = "fail"
            msg.FileInfo = ""
            msg.Message = "failed to put object"
            log.Printf("failed to put object: %s", err.Error())
            c.JSON(http.StatusInternalServerError, msg)
            return
        }
        log.Printf("put object %s done", objectKey)
	    msg.Code = http.StatusOK
	    msg.Status = "success"
	    msg.FileInfo = objectKey
	    msg.Message = "upload log file successfully"
	    c.JSON(http.StatusOK, msg)
	    return
	} else {
		msg.Code = http.StatusNotImplemented
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "backend implementing..."
		c.JSON(http.StatusNotImplemented, msg)
		return
	}
}

// upload pos db backup
func UploadDB(c *gin.Context) {
	var input DBUploadInput
	var msg ResponseFile

	if c.ShouldBind(&input) != nil {
		msg.Code = http.StatusBadRequest
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "failed to bind parameters"
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	header, err := c.FormFile("uploadfile")
	if err != nil {
		msg.Code = http.StatusBadRequest
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "failed to get upload file"
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	if ok := utils.FilterExt(header.Filename, ".zip") || utils.FilterExt(header.Filename, ".sql"); !ok {
		msg.Code = http.StatusBadRequest
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "only zip or sql file is allowed"
		c.JSON(http.StatusBadRequest, msg)
		return
	}

    // verify
	if config.Cfg.Verification.Enable {
		message := input.GetMessage(header.Filename, config.Cfg.Verification.Salt)
		if !utils.CheckDigest(message, input.Digest, config.Cfg.Verification.Method) {
			msg.Code = http.StatusBadRequest
			msg.Status = "fail"
			msg.FileInfo = ""
			msg.Message = "unexpect digest"
			c.JSON(http.StatusBadRequest, msg)
			return
		}
	}

	file, err := header.Open()
	if err != nil {
		msg.Code = http.StatusInternalServerError
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "file error"
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer file.Close()

	// upload
    var objectKey string
	if config.Cfg.Backend == "oss" {
        objectKey = input.GetKey(config.Cfg.DbBackup.Prefix, header.Filename)
        if err := alioss.Client.PutObject(objectKey, file); err != nil {
            msg.Code = http.StatusInternalServerError
            msg.Status = "fail"
            msg.FileInfo = ""
            msg.Message = "failed to put object"
            log.Printf("failed to put object: %s", err.Error())
            c.JSON(http.StatusInternalServerError, msg)
            return
        }
        log.Printf("put object %s done", objectKey)
	    msg.Code = http.StatusOK
	    msg.Status = "success"
	    msg.FileInfo = objectKey
	    msg.Message = "upload db backup successfully"
	    c.JSON(http.StatusOK, msg)
	    return
	} else {
		msg.Code = http.StatusNotImplemented
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "backend implementing..."
		c.JSON(http.StatusNotImplemented, msg)
		return
	}
}

// mock
func DownloadBill(c *gin.Context) {
	var msg ResponseFile
	filename := c.DefaultQuery("filename", "")
	msg.Code = http.StatusOK
	msg.Status = "success"
	msg.FileInfo = filename
	msg.Message = "mock api"
	c.JSON(http.StatusOK, msg)
	return
}

// upload bill file
func UploadBill(c *gin.Context) {
	var msg ResponseFile
	header, err := c.FormFile("uploadfile")
	if err != nil {
		msg.Code = http.StatusBadRequest
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "failed to get upload file"
		c.JSON(http.StatusBadRequest, msg)
		return
	}
	if ok := utils.FilterExt(header.Filename, ".zip"); !ok {
		msg.Code = http.StatusBadRequest
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "only zip file is allowed"
		c.JSON(http.StatusBadRequest, msg)
		return
	}
	file, err := header.Open()
	if err != nil {
		msg.Code = http.StatusInternalServerError
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "file error"
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer file.Close()

	// upload
    var objectKey string
    var location string
	if config.Cfg.Backend == "oss" {
        objectKey = config.Cfg.BillFile.Prefix + "/" + header.Filename
		location = config.Cfg.BillFile.Target + objectKey
        if err := alioss.Client.PutObject(objectKey, file); err != nil {
            msg.Code = http.StatusInternalServerError
            msg.Status = "fail"
            msg.FileInfo = ""
            msg.Message = "failed to put object"
            log.Printf("failed to put object: %s", err.Error())
            c.JSON(http.StatusInternalServerError, msg)
            return
        }
        log.Printf("put object %s done", objectKey)
	    msg.Code = http.StatusOK
	    msg.Status = "succ"
	    msg.FileInfo = location
	    msg.Message = "upload bill successfully"
	    c.JSON(http.StatusOK, msg)
	    return
	} else {
		msg.Code = http.StatusNotImplemented
		msg.Status = "fail"
		msg.FileInfo = ""
		msg.Message = "backend implementing..."
		c.JSON(http.StatusNotImplemented, msg)
		return
	}
}

func AppRoutes(r *gin.Engine) {
	// old routes
	r.POST("/doc", UploadLog)
	r.POST("/dbBackup", UploadDB)
	r.GET("/upload", DownloadBill)
	r.POST("/upload", UploadBill)
}
