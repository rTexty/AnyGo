package main

import (
	"errors"
	"strings"
	"unicode"
)

func Unpack(str string) (string, error) {
    var escape bool = false
    var lastRune rune
    var hasLast bool
    var builder strings.Builder

    for _, r := range str {
        if escape == true {
            builder.WriteRune(r)
            lastRune = r
            hasLast = true
            escape = false
        } else if unicode.IsDigit(r) && hasLast == true {
            count := int(r-'0') - 1
            for range count {
                builder.WriteRune(lastRune)
                }
            hasLast = false
            lastRune = r
        } else if string(r) == "\\"{
            escape = true
        } else {
            builder.WriteRune(r)
            lastRune = r
            hasLast = true
        }
    }
    if escape == true{
        return "", errors.New("String finishes with slash")
    }
    out := builder.String()
    return out, nil
}

func main() {

}
