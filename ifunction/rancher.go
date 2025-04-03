package ifunction

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Access struct {
	AccessKey string
	SecretKey string
}

type Rancher struct {
	access  *Access
	baseUrl string
	client  *http.Client
	conn    *tls.Conn
}

func NewRancher(auth *Access, baseUrl string) *Rancher {
	// 创建自定义 TLS 配置，动态获取服务器证书
	conn, err := tls.Dial("tcp", strings.Split(baseUrl, "//")[1], &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	certs := conn.ConnectionState().PeerCertificates
	caCertPool := x509.NewCertPool()
	for _, cert := range certs {
		caCertPool.AddCert(cert)
	}

	// 配置客户端
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	return &Rancher{
		access:  auth,
		baseUrl: baseUrl,
		conn:    conn,
		client:  client,
	}
}

func buildUrl(baseUrl string, uri string, params url.Values) (requestUrl string) {
	addr := strings.Split(baseUrl, ":")
	requestUrl = addr[0] + ":" + addr[1] + uri
	if params != nil {
		queryString := params.Encode()
		if queryString != "" {
			requestUrl = requestUrl + "?" + queryString
		}
	}
	return
}

func buildGetReq(access *Access, rancherURL string) *http.Request {
	// 创建请求
	req, err := http.NewRequest("GET", rancherURL, nil)
	if err != nil {
		log.Fatal("创建请求失败:", err)
	}
	encode := base64.StdEncoding.EncodeToString([]byte(access.AccessKey + ":" + access.SecretKey))
	req.Header.Set("Authorization", "Basic "+encode)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func (rancher *Rancher) GetPod(clusterID string, namespace string, podName string) Pod {
	// 目标端点：获取单个 Pod 的信息
	req := buildGetReq(rancher.access, buildUrl(rancher.baseUrl, fmt.Sprintf("/k8s/clusters/%s/v1/pods/%s/%s", clusterID, namespace, podName), nil))
	// 发送请求
	resp, err := rancher.client.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	// 关闭资源
	defer closeResource(rancher.conn, resp)
	// 检查状态码
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("请求失败，状态码: %d, 响应: %s\n", resp.StatusCode, err)
	}
	var pod Pod
	if err := json.Unmarshal(body, &pod); err != nil {
		log.Fatalln("解析 JSON 失败:", err)
	}
	return pod
}

func (rancher *Rancher) GetPodList(clusterID string, namespace string) PodList {
	req := buildGetReq(rancher.access, buildUrl(rancher.baseUrl, fmt.Sprintf("/k8s/clusters/%s/v1/pods?namespace=%s", clusterID, namespace), nil))
	resp, err := rancher.client.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	defer closeResource(rancher.conn, resp)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("读取响应失败:", err)
	}
	var pods PodList

	if err := json.Unmarshal(body, &pods); err != nil {
		log.Fatal(err)
	}

	return pods
}

func (rancher *Rancher) GetClusters(namespace string) ClusterResponse {
	req := buildGetReq(rancher.access, buildUrl(rancher.baseUrl, "/v3/clusters", nil))
	resp, err := rancher.client.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	defer closeResource(rancher.conn, resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("读取响应失败:", err)
	}
	var clusterResponse ClusterResponse

	if err := json.Unmarshal(body, &clusterResponse); err != nil {
		log.Fatal(err)
	}
	datas := clusterResponse.Data

	if len(datas) == 0 {
		log.Fatal("未找到集群")
	}
	if len(namespace) > 0 {
		for _, data := range datas {
			if data.Name == namespace {
				return ClusterResponse{
					Data: []Cluster{
						{
							ID:   data.ID,
							Name: data.Name,
						},
					},
				}
			}
		}
	} else {
		return clusterResponse
	}
	return ClusterResponse{}
}

func closeResource(conn *tls.Conn, resp *http.Response) {
	resp.Body.Close()
	conn.Close()
}
