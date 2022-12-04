package main

import (
  "net/http"
  "path/filepath"
  "github.com/gin-gonic/gin"
)

func upload(c *gin.Context) {  
  // Pull the file field out of the form
  file, err := c.FormFile("file")
  if err != nil {
    c.String(http.StatusBadRequest, "get form err: %s", err.Error())
    return
  }
  
  // Join the filename from the form with the files directory so
  // the files can be saved there
  filename := filepath.Join("files",filepath.Base(file.Filename))
  // save the file
  if err := c.SaveUploadedFile(file, filename); err != nil {
    c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
    return
  }
  
  c.String(http.StatusOK, "File %s uploaded successfully with fields as %s", file.Filename, filename)
}

func download(c *gin.Context) {  
  // Find the filename to download from the path param
  filename := c.Param("filename")
  // Join the name with the files directory
  fullpath := filepath.Join("files", filename)
  // return the file
  c.File(fullpath)
}

func ping(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "message": "pong",
  })
}


func main() {
  // Initialize a default gin router
  router := gin.Default()

  // Set the maximum memort used for multi-part uploads
  router.MaxMultipartMemory = 8 << 20 // 8 MiB

  // Adding a:
  //    * POST /upload
  //    * GET /download/filename
  //    * GET /ping
  router.POST("/upload", upload)
  router.GET("/download/:filename", download)
  router.GET("/ping", ping)
  router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}