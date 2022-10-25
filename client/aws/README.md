# AWS

## Introduction
This package is used for configuration request to aws. 
What's got in this package.
1. UploadFileInByte - Upload a file using the bytes type

## Using Package
## Using Package
```go
    moladin_evo := NewS3(
        context,
        validator.New(), // import from go-lib-common/validator
        WithSentry(sentry), // is required field | import from go-lib-common/sentry
        WithAwsS3Region(awsRegion), // is required field  | The region to send requests to. A full list of regions is found in the "Regions and Endpoints" document.
        WithAwsS3AccessKeyID(accessKeyId), // is required field | AWS Access key ID
        WithAwsS3SecretAccessKey(accessKey), // is required field | AWS Secret Access Key
        WithAwsS3Arn(arn), // is optional field | AWS S3 Arn
        WithAwsS3ACL(acl), // is required field | The canned ACL to apply to the object.
        WithAwsS3BucketName(bucketName), // is required field | name of the bucket to which the PUT operation was initiated.
        WithAwsS3PresignTimeInMinutes(presignTime), // is required field | Presign returns the request's signed URL
    )
```

### Using UploadFileInByte
UploadFileInByte has 4 parameters
1. context
2. path - the path of the folder to be placed the file
3. fileName - the name of the file to be uploaded
4. data - data object to be uploaded

```go

aws.UploadFileInByte(ctx,"myBucket/subBucket","data.csv",byteData)
```

