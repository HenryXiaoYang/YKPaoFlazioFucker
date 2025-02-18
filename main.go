package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	concurrentClients = 10 // 并发客户端数量
)

// 线程安全的随机数生成器
func generatePassword(length int, r *rand.Rand) string {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[r.Intn(len(characters))]
	}
	return string(b)
}

// 每个客户端独立运行的goroutine
func sendRequest(wg *sync.WaitGroup, clientID int) {
	defer wg.Done()

	// 为每个goroutine创建独立的随机源
	src := rand.NewSource(time.Now().UnixNano() + int64(clientID))
	r := rand.New(src)

	// 复用HTTP客户端（支持连接池）
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
	}

	for { // 持续发送请求
		// 生成动态凭证
		part1 := r.Intn(18) + 10
		part2 := r.Intn(401) + 100
		username := fmt.Sprintf("s%d%d@ykpaoschool.cn", part1, part2)
		password := generatePassword(7, r)

		// 构建表单数据
		data := url.Values{}
		data.Set("version", "225-299")
		data.Set("__ca", "603")
		data.Set("i", "10")
		data.Set("attributi", "|")
		data.Set("valori", fmt.Sprintf("%s|%s", username, password))
		data.Set("destinatari", "")
		data.Set("f", "generamail")
		data.Set("lingua", "it")
		data.Set("hook", "Lw==")

		// 创建请求
		req, err := http.NewRequest(
			"POST",
			"https://ykpaoschool-onmicrosoft.flazio.com/manager/includer.php",
			strings.NewReader(data.Encode()),
		)
		if err != nil {
			fmt.Println(fmt.Sprintf("线程%d 错误:%v", clientID, err))
			continue
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(fmt.Sprintf("线程%d 错误:%v", clientID, err))
			continue
		}

		// 读取响应（根据实际需要处理）
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("线程%d 错误:%v", clientID, err))
			continue
		}

		fmt.Println(fmt.Sprintf("线程%d 随机用户名:%s 随机密码:%s 成功:%s", clientID, username, password, string(body)))
	}
}

func main() {
	fmt.Println("让我们一起用假数据塞满黑客的数据库！✊")
	fmt.Println("并发客户端数量:%d", concurrentClients)
	fmt.Println("返回值成功为1则代表成功")
	fmt.Println("3秒后开始发送请求...")
	time.Sleep(3 * time.Second)
	fmt.Println("开始发送请求...")

	var wg sync.WaitGroup

	// 启动并发客户端
	for i := 0; i < concurrentClients; i++ {
		wg.Add(1)
		go sendRequest(&wg, i)
	}

	// 保持主线程运行
	wg.Wait()
}
