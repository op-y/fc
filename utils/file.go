package utils

import (
    "strings"
)

// Get extension name of file
func GetExt(filename string) string {
    parts := strings.Split(filename, ".")
    if len(parts) >= 2 {
        return parts[len(parts)-1]
    }
    return filename
}

// Filter certain extension name
func FilterExt(filename, ext string) bool {
    return strings.HasSuffix(filename, ext)
}

