# GCP

## Introduction
This package is used for configuration request to gcp. 
What's got in this package.
1. UploadFileInByte - Upload a file using the bytes type

## Using Package
## Using Package
```go
    gcsServiceAccount := gcp.ServiceAccountKeyJSON{
		Type:                    Type,
		ProjectId:               ProjectId,
		PrivateKeyId:            PrivateKeyId,
		PrivateKey:              PrivateKey,
		ClientEmail:             ClientEmail,
		ClientId:                ClientId,
		AuthUri:                 AuthUri,
		TokenUri:                TokenUri,
		AuthProviderX509CertUrl: AuthProviderX509CertUrl,
		ClientX509CertUrl:       ClientX509CertUrl,
    } // Generated from GCP Console Service Account Keys | is required
    moladin_evo := NewGCP(
        context,
        validator.New(), // import from go-lib-common/validator
        WithSentry(sentry), // is required field | import from go-lib-common/sentry

        gcp.WithServiceAccountKeyJSON(gcsServiceAccount), // is required field |  Generated from GCP Console Service Account Keys
        gcp.WithSignedUrlTimeInMinutes(cfg.Cold.GCSSignedUrlTimeInMinutes), // is required field | Set Signed URL Expires Time in Minutes
        gcp.WithBucketName(cfg.Cold.GCSBucketName)), // Name of bucket 
    )
```

### Using UploadFileInByte
UploadFileInByte has 3 parameters
1. context
2. fileName - the name of the file to be uploaded
3. data - data object to be uploaded

Cloud Storage operates with a flat namespace, which means that folders don't actually exist within Cloud Storage. If you create an object named `folder1/file.txt` in the bucket `your-bucket`, the path to the object is `your-bucket/folder1/file.txt`, but there is no folder named `folder1`; instead, the string `folder1` is part of the object's name.

```go

gcp.UploadFileInByte(ctx,"myBucket/subBucket/data.csv",byteData)
```

