package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	ccnuv1 "github.com/MuxiKeStack/be-api/gen/proto/ccnu/v1"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func (c *ccnuService) Login(ctx context.Context, studentId string, password string) (bool, error) {
	// 划分一下本科生、研究生
	var (
		client *http.Client
		err    error
	)
	if len(studentId) > 5 && studentId[4] == '2' {
		// 本科生
		client, err = c.loginUndergraduateClient(ctx, studentId, password)
	} else {
		client, err = c.loginPostgraduateClient(ctx, studentId, password)
	}
	return client != nil, err
}

func (c *ccnuService) client() *http.Client {
	j, _ := cookiejar.New(&cookiejar.Options{})
	return &http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Jar:     j,
		Timeout: c.timeout,
	}
}

func hexToBigInt(hexStr string) (*big.Int, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	bigInt := new(big.Int)
	bigInt.SetBytes(bytes)
	return bigInt, nil
}

func (c *ccnuService) loginPostgraduateClient(ctx context.Context, studentId string, password string) (*http.Client, error) {
	return &http.Client{}, nil
	//modulus, exponent := c.getPublicKey()
	//// 将modulus和exponent从hex转换为big.Int
	//modulus, err := hexToBigInt(modulusHex)
	//if err != nil {
	//	fmt.Println("Invalid modulus:", err)
	//	return
	//}
	//
	//exponent, err := hexToBigInt(exponentHex)
	//if err != nil {
	//	fmt.Println("Invalid exponent:", err)
	//	return
	//}
	//
	//// 创建RSA公钥
	//publicKey := &rsa.PublicKey{
	//	N: modulus,
	//	E: int(exponent.Int64()), // 注意：E一般是小整数，如65537
	//}
	//
	//// 要加密的明文数据
	//plaintext := []byte("your-password")
	//
	//// 使用RSA/OAEP加密
	//ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
	//if err != nil {
	//	fmt.Println("Error encrypting:", err)
	//	return
	//}
	//
	//// 将加密后的数据转换为Base64
	//ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)
}

func (c *ccnuService) getPublicKey() (modulus string, exponent string) {
	req, err := http.NewRequest("GET", "https://grd.ccnu.edu.cn/yjsxt/xtgl/login_getPublicKey.html?time=1726134051870&_=1726133922646", nil)
	if err != nil {
		return "", ""
	}

	// 添加请求头
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://grd.ccnu.edu.cn/yjsxt/xtgl/login_slogin.html?time=1726133573723")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Microsoft Edge";v="128"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	type Result struct {
		Modulus  string `json:"modulus"`
		Exponent string `json:"exponent"`
	}
	res := &Result{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", ""
	}
	return res.Modulus, res.Exponent
}

func (c *ccnuService) loginUndergraduateClient(ctx context.Context, studentId string, password string) (*http.Client, error) {
	params, err := c.makeAccountPreflightRequest()
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("username", studentId)
	v.Set("password", password)
	v.Set("lt", params.lt)
	v.Set("execution", params.execution)
	v.Set("_eventId", params._eventId)
	v.Set("submit", params.submit)

	request, err := http.NewRequest("POST", "https://account.ccnu.edu.cn/cas/login;jsessionid="+params.JSESSIONID, strings.NewReader(v.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")
	request.WithContext(ctx)

	client := c.client()
	resp, err := client.Do(request)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			return nil, ccnuv1.ErrorNetworkToXkError("网络异常")
		}
		return nil, err
	}
	if len(resp.Header.Get("Set-Cookie")) == 0 {
		return nil, ccnuv1.ErrorInvalidSidOrPwd("学号或密码错误")
	}
	return client, nil
}

type ClientKey struct{} // 用于 context 的键

// 将 http.Client 添加到 context 中
func (c *ccnuService) addClientToContext(ctx context.Context, client *http.Client) context.Context {
	return context.WithValue(ctx, ClientKey{}, client)
}

// 从 context 中获取 http.Client
func (c *ccnuService) getClientFromContext(ctx context.Context) *http.Client {
	client, ok := ctx.Value(ClientKey{}).(*http.Client)
	if !ok {
		return nil // 这里可以处理默认逻辑或错误
	}
	return client
}
