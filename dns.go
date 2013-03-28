package main
import (
  "io/ioutil"
  "log"
  "net"
  "net/http"
  "net/url"
  "strings"
  "time"
)
var (
    body = url.Values{
        "login_email":    {"swallow197@qq.com"},
        "login_password": {"*****"},
        "format":         {"json"},
        "domain_id":      {"2989344"},
        "record_id":      {"21001935"},
        "sub_domain":     {"www"},
        "record_line":    {"默认"},
    }
    current_ip     = ""
    check_interval = 30 * time.Second
    timeout = 60 * time.Second
)
func get_public_ip() (string, error) {
    conn, err := net.DialTimeout("tcp", "ns1.dnspod.net:6666", timeout)
    defer func() {
        if x := recover(); x != nil {
            log.Println("Can't get public ip", x)
        }
        if conn!=nil {
            conn.Close()
        }
    }()
    if err == nil {
        var bytes []byte
        deadline := time.Now().Add(timeout)
        err = conn.SetDeadline(deadline)
        if err!=nil {
            return "", err
        }
        bytes, err = ioutil.ReadAll(conn)
        if err == nil {
            return string(bytes), nil
        }
    }
    return "", err
}
func timeoutDialler(timeout time.Duration) func(net, addr string) (c net.Conn, err error) {
    return func(netw, addr string) (net.Conn, error) {
        c, err := net.DialTimeout(netw, addr, timeout)
        if err != nil {
            return nil, err
        }
        deadline := time.Now().Add(timeout)
        err = c.SetDeadline(deadline)
        if err != nil {
            return nil, err
        }
        return c, nil
    }
}
func update_dnspod(ip string) bool {
    client := &http.Client{
        Transport: &http.Transport{
            Dial: timeoutDialler(timeout),
        },
    }
    body := url.Values{
        "login_email":    {"swallow197@qq.com"},
        "login_password": {"****"},
        "format":         {"json"},
        "domain_id":      {"2989344"},
        "record_id":      {"21001935"},
        "sub_domain":     {"www"},
        "record_line":    {"默认"},
        "value":          {ip},
    }
    req, err := http.NewRequest("POST", "https://dnsapi.cn/Record.Ddns", strings.NewReader(body.Encode()))
    req.Header.Set("Accept", "text/json")
    req.Header.Set("Content-type", "application/x-www-form-urlencoded")
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        return false
    }
    bytes, _ := ioutil.ReadAll(resp.Body)
    log.Println(string(bytes))
    return resp.StatusCode == 200
}
func init() {
}
func main() {
    for {
        ip, err := get_public_ip()
        if ip != "" && err == nil {
            log.Println("got ip:" + ip)
            if ip != current_ip {
                log.Println("update dnspod with new ip:" + ip)
                if update_dnspod(ip) {
                    current_ip = ip
                }
            }
        } else {
            log.Println("error:", err)
        }
        time.Sleep(check_interval)
    }
}
