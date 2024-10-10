package chip

import (
	"fmt"
	"testing"
	"time"
)

func TestFormat_FriTime(t *testing.T) {
	f := Format{}
	pastTime := time.Now().Add(-24 * 30 * 8 * time.Hour)
	fmt.Println(f.FriTime(pastTime))
}

func TestFormat_FriNumber(t *testing.T) {
	f := Format{}
	fmt.Println(f.FriNumber(120))
}

func TestFormat_ToHTTPS(t *testing.T) {
	f := Format{}
	fmt.Println(f.ToHTTPS("http://www.baidu.com"))
	fmt.Println(f.ToHTTPS("../detail.html"))
}

func TestFormat_StripScheme(t *testing.T) {
	f := Format{}
	fmt.Println(f.StripScheme("http://www.baidu.com"))
	fmt.Println(f.StripScheme("https://www.baidu.com"))
}
