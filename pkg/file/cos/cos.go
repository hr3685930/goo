package cos

import (
    "goo/pkg/file"
    "github.com/tencentyun/cos-go-sdk-v5"
    "net/http"
    "net/url"
    "time"
)

type cosFile struct {
    secretID  string
    secretKey string
    region    string
    bucket    string
    Timeout   time.Duration
}

func NewCosFile(secretID string, secretKey string, region string, bucket string, timeout time.Duration) *cosFile {
    return &cosFile{secretID: secretID, secretKey: secretKey, region: region, bucket: bucket, Timeout: timeout}
}

func (c *cosFile) InitCos() {
    u, _ := url.Parse("https://" + c.bucket + ".cos." + c.region + ".myqcloud.com")
    b := &cos.BaseURL{BucketURL: u}
    file.Fs = cos.NewClient(b, &http.Client{
        Timeout: c.Timeout,
        Transport: &cos.AuthorizationTransport{
            SecretID:  c.secretID,
            SecretKey: c.secretKey,
        },
    })

}
