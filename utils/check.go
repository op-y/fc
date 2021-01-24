package utils

import (
    "crypto/md5"
    "crypto/sha256"
    "encoding/hex"
)

func CheckDigest(message, digest, vtype string) bool {
    switch vtype {
    case "md5":
        return CheckMD5Sum(message, digest)
    case "sha256":
        return CheckSHA256Sum(message, digest)
    default:
        return false
    }
}

func CheckMD5Sum(message, digest string) bool {
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(message))
    signature := hex.EncodeToString(md5Ctx.Sum(nil))

    if signature == digest {
        return true
    } else {
        return false
    }
}

func CheckSHA256Sum(message, digest string) bool {
    sha256Ctx := sha256.New()
    sha256Ctx.Write([]byte(message))
    signature := hex.EncodeToString(sha256Ctx.Sum(nil))

    if signature == digest {
        return true
    } else {
        return false
    }
}

