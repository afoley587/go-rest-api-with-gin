# Lightning Fast REST API With Go in Under 50 Lines Of Code

I've been using python for APIs for a while. I've bounced
between Django, FastAPI, Flask, and a slew of others as well,
but I have never use Go to do it. And since Go is like the standard
language to build super cool tools in my field (DevOps / SRE), I figured
it was a good time to give it a try.

So, on the docket today: Let's build a blazing fast REST API with Go using
the [gin gonic framework](https://github.com/gin-gonic/gin). 
The best part? It's going to be less than 50 lines of Go code (excluding comments!).

I've decided that I want my REST API to support file uploads, file downloads,
and a health check ping method. In more requirement-oriented terms:

* I want to be able to upload files with a `POST` to `/upload`
* I want to download files with a `GET` to `/download/<filename>`
* I want to be able to get ping-pong responses with a `GET` to `/ping`

Now, let's get coding!

# Building the API

## Dependencies
As a disclaimer, I will assume you have Go installed with a version of 1.16+!

Because this is Go, we need to initialize our module:

```shell
go mod init uploaddownload.go
go mod tidy
```

And because we plan on using the `gin` framework using `go get`:

```shell
go get -u github.com/gin-gonic/gin
```

Now our workspace is set up and all dependencies are ready for usage.

## Main Function
Let's open up `uploaddownload.go` and look at its `main` method:

```go
func main() {
  // Initialize a default gin router
  router := gin.Default()

  // Set the maximum memory used for multi-part uploads
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
```

The first thing we do is initialize our `gin` router. This really does all of the
heaavy lifting for our API. It will handle:

* Running the server
* Handling the appropriate requests
* Logging and even time tracing
* Much more

Next, because we do plan on handling file uploads, we set a max memory size for 
multi-part uploads. Finally, we add our REST endpoints.

These correspond 1-to-1 with the requirements we described in the previous section.

Let's break down one of these router lines:

```go
router.POST("/upload", upload)
```

So, this adds an endpoint to our server at `/upload`. It is of type `POST` and when it's hit,
the `gin` framework will call our `upload` method (shown below) with the calling context (such
as request information). If I wanted to create a client request with `cURL`, it might look like:

```shell
curl -X POST --form file="@some/file/pathf" http://localhost:8080/upload
```

## Router Methods
We've described our main router and how that's used, but it won't be of much
use to us if those function callbacks don't exist. Let's go through them in the
same order they are listed in the router, starting with `upload`:

```go
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
```

You'll notice the first thing we do is try to open the multi-part form upload from the 
requester. We are assuming there is a field called `file`. A caller might upload to this endpoint
using `curl -X POST --form file="@some/file/pathf" ...`. If that part is omitted from the
form, we will return an HTTP 400 to signify a bad request. Next, we save the file to disk in
`files/<filename>` using the `SaveUploadedFile` method from the `gin` framework. This
takes an uploaded file and saves it to the filename specified, and returns an error if that's
not possible.

Our `download` method is even simpler actually:

```go
func download(c *gin.Context) {  
  // Find the filename to download from the path param
  filename := c.Param("filename")
  // Join the name with the files directory
  fullpath := filepath.Join("files", filename)
  // return the file
  c.File(fullpath)
}
```

The first thing it's going to do is pull the filename from a path parameter using `c.Param`. A
caller might have requested the file as `/downloads/filenametodownload`, so `c.Param` would 
correspond to `filenametodownload`. We then prepend our file directory name (`files`) to the 
requested file and return the file with `c.File` provided by the `gin` framework.

Our last method is `ping`, which just returns a `JSON` representation of the `ping-pong`
status messages common in kubernetes applications:

```go
func ping(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "message": "pong",
  })
}
```

# Running
Now we can run it! In one terminal, you can run 

```shell
go run uploaddownload.go
```

And in another terminal, you can test uploads:

```shell
curl -X POST --form file="@some/path/to/a/file" http://localhost:8080/upload
```

Or test downloads:

```shell
curl -X GET http://localhost:8080/download/file-to-download
```

As an example, I wrote this test script that loops through
a few test files in my `test-upload-files` directory and uploads them
to the server. It then goes through the `files` directory, which is where the server
stores uploaded files, and then calls the download method on them:

```shell
# A directory of test files to upload
TEST_FILES_DIR=test-upload-files
# All of the files downloadable by the server
TEST_UPLOAD_DIR=files

for f in $(ls $TEST_FILES_DIR); do
  curl -X POST --form file="@$TEST_FILES_DIR/$f" http://localhost:8080/upload
  echo ""
done

for f in $(ls $TEST_UPLOAD_DIR); do
  curl -X GET http://localhost:8080/download/$f
  echo ""
done
```

The responses returned are:

```shell
File test1 uploaded successfully with fields as files/test1
File test2 uploaded successfully with fields as files/test2
hello world!
!dlrow olleh
```

# References
* All code can be found in [this github repo](https://github.com/afoley587/go-rest-api-with-gin)!
* More on the [gin gonic framework](https://github.com/gin-gonic/gin)