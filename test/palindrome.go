package main

import "unicode"

func IsPalindrome(s string) bool {
	left, right := 0, len(s)-1
    for left < right {
        //left pointer: skip all char except letter only
        if !unicode.IsLetter(rune(s[left])) && !unicode.IsDigit(rune(s[left])){
            left++
            continue
        }
        //right pointer: skip all char except letter only
        if !unicode.IsLetter(rune(s[right])) && !unicode.IsDigit(rune(s[right])){
            right--
            continue
        }

        if unicode.ToLower(rune(s[left])) != unicode.ToLower(rune(s[right])){
            return false
        }

        left++
        right--
    }
    return true
}