package common

//  utils
import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"strconv"
	"strings"
	"time"
)

func PanicIf(err error) {

	if err != nil {

		panic(err)
	}
}

func ParseInt(value string) int {

	if value == "" {
		return 0
	}

	val, _ := strconv.Atoi(value)
	return val
}

func IntString(value int) string {

	return strconv.Itoa(value)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetLocalAddr() string {
	conn, err := net.Dial("udp", "localhost:80")
	if err != nil {
		return ""
	}

	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func GetIPFromUrl(url string) string {

	url = strings.Split(url, "//")[1]
	url = strings.Split(url, "/")[0]
	if strings.Contains(url, ":") {

		return strings.Split(url, ":")[0]
	}

	return url

}

func Now() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}
