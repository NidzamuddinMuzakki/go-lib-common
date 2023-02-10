# Encryptiom

## Introduction
This package is used for encrypt and decrypt data. What's got in this package.
1. GenerateSalt - Generate salt for encryption and decryption
2. Encrypt - Encrypt data
3. Decrypt - Decrypt data

## Algorithm
- Blowfish

## Using Package

```go
validator := commonValidator.New()

encryption := commonEncryption.NewEncryption(
    validator,
    commonEncryption.WithAppName("service-name"),
)
```

### Using GenerateSalt
```go
salt := encryption.GenerateSalt("secret-key")
```

### Using Encrypt
```go
encryptionData := encryption.Encrypt("data to encrypt", salt)
// if you want to convert to string
encryptionDataToStr := hex.EncodeToString(encryptionData)
```

### Using Decrypt
```go
decryptionData := encryption.Decrypt("b277bd0b596fc04d41fc01e3b559675d", salt)
// convert to string
fmt.Printf("%s", decryptionData)
```