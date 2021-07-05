package utils

import (
    "bufio"
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "os"
    "reflect"
    "strings"
)


func Md5(target string) string {
    h := md5.New()
    h.Write([]byte(target))
    return hex.EncodeToString(h.Sum(nil))
}

//计算文件的md5，适用于本地文件计算
func FileMd5(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()
    md5hash := md5.New()
    if _, err := io.Copy(md5hash, f); err != nil {
        return "", err
    }
    return hex.EncodeToString(md5hash.Sum(nil)), nil
}



// 判断一个值是否在切片中
func InArray(need interface{}, haystack interface{}) (exists bool, index int) {
    exists = false
    index = -1
    switch reflect.TypeOf(haystack).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(haystack)
        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(need, s.Index(i).Interface()) == true {
                index = i
                exists = true
                return
            }
        }
    }
    return
}

// 下载远程文件到服务器
func DownloadRemoteFile(fileUrl, filePath string) error {
    response, err := http.Get(fileUrl)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    reader := bufio.NewReader(response.Body)
    file, err := os.Create(filePath)
    if err != nil {
        panic(err)
    }

    writer := bufio.NewWriter(file)
    _, err = io.Copy(writer, reader)
    if err != nil {
        return err
    }
    return nil
}


func Explode(delimiter, text string) []string {
    if len(delimiter) > len(text) {
        return strings.Split(delimiter, text)
    } else {
        return strings.Split(text, delimiter)
    }
}

func JsonEncode(data interface{}) (string, error) {
    jsons, err := json.Marshal(data)
    return string(jsons), err
}

func JsonDecode(data string) (map[string]interface{}, error) {
    var dat map[string]interface{}
    err := json.Unmarshal([]byte(data), &dat)
    return dat, err
}


func RemoveDuplicateElement(arr []string) []string {
    resArr := make([]string, 0)
    tmpMap := make(map[string]interface{})
    for _, val := range arr {
        //判断主键为val的map是否存在
        if _, ok := tmpMap[val]; !ok {
            resArr = append(resArr, val)
            tmpMap[val] = nil
        }
    }

    return resArr
}


func SubString(str string, start, end int) string {
    rs := []rune(str)
    length := len(rs)

    if start < 0 || start > length {
        panic("start is wrong")
    }

    if end < start || end > length {
        panic("end is wrong")
    }

    return string(rs[start:end])
}
