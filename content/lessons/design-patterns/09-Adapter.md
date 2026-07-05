# Adapter Pattern

The Adapter Pattern is a structural design pattern that allows incompatible interfaces to work together. It acts as a wrapper, translating calls from one interface into a format that the other interface can understand.

Think of a travel adapter: you have a US laptop plug, but you are in Europe. You cannot change the laptop, and you cannot change the wall socket. You insert an Adapter between them.

## 1. The Problem

You are building a Go application that saves data to AWS S3. 
Your internal `Uploader` service requires objects that implement a specific `Storage` interface.

```go
type Storage interface {
    UploadFile(filename string, data []byte) error
}
```

You import a 3rd-party library (like the official AWS SDK) to handle the actual S3 uploading. However, the AWS SDK doesn't have an `UploadFile` method. It has a completely different signature:

```go
// 3rd Party AWS SDK (You cannot modify this code!)
type AWSS3Client struct{}

func (c *AWSS3Client) PutObject(bucket, key string, body []byte) error {
    fmt.Println("Uploading to S3...")
    return nil
}
```

Your `Uploader` service refuses to accept the `AWSS3Client` because it doesn't satisfy the `Storage` interface!

## 2. The Solution (The Adapter)

We create a new struct (the Adapter) that explicitly implements *our* `Storage` interface, but secretly calls the *3rd-party* methods inside.

```go
// 1. The Adapter Struct
type S3Adapter struct {
    awsClient *AWSS3Client
    bucket    string
}

// 2. The Constructor
func NewS3Adapter(client *AWSS3Client, bucket string) *S3Adapter {
    return &S3Adapter{
        awsClient: client,
        bucket:    bucket,
    }
}

// 3. Implementing our internal Interface!
func (a *S3Adapter) UploadFile(filename string, data []byte) error {
    // 4. Translating the call to the incompatible 3rd-party signature
    return a.awsClient.PutObject(a.bucket, filename, data)
}
```

## 3. The Usage

Now, we can seamlessly pass the 3rd-party AWS client into our internal system, completely decoupled from the AWS implementation details.

```go
func main() {
    // The incompatible 3rd-party client
    aws := &AWSS3Client{}
    
    // Wrap it in our Adapter
    adapter := NewS3Adapter(aws, "my-bucket")
    
    // Now it perfectly satisfies our internal requirements!
    processData(adapter)
}

// Internal business logic
func processData(store Storage) {
    store.UploadFile("test.txt", []byte("hello"))
}
```

## 4. Architectural Value

The Adapter pattern is the bedrock of **Clean Architecture** (Hexagonal Architecture / Ports and Adapters). 

By wrapping 3rd-party libraries in Adapters, you ensure that your core business logic is never tainted by external dependencies. If you decide to switch from AWS S3 to Google Cloud Storage next year, you do not touch your `processData` business logic! You simply write a new `GCSAdapter` and plug it in.
