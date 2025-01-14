package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"myadmin/model"
	"myadmin/model/blog"
	"myadmin/model/vlog"
	"myadmin/model/vod"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

func DownloadFile(urlPath string, wg *sync.WaitGroup) error {
	defer wg.Done()

	// 判断文件是否存在
	fileRealPath := filepath.Join("/data", urlPath)
	if fsInfo, err := os.Stat(fileRealPath); !os.IsNotExist(err) {
		// 如果文件存在。
		if fsInfo.Size() > 0 {
			return nil
		}
	} else { // 避免文件夹不存在。
		err := os.MkdirAll(filepath.Dir(fileRealPath), 0755)
		if err != nil {
			log.Println("创建文件夹错误：", filepath.Dir(fileRealPath), err)
			return err
		}
	}

	client := http.Client{
		Timeout: time.Second * 10, // 设置超时时间为30秒
	}

	retryCount := 3 // 设置重试次数
	reqURL := ""
	for i := 0; i < retryCount; i++ {

		if i == 0 {
			reqURL = "http://source.ermeixk1128.com/" + urlPath
		} else if i == 1 {
			reqURL = "http://mogushipin.oss-cn-hangzhou.aliyuncs.com/" + urlPath
		} else if i == 2 {
			reqURL = "http://mogushipin.oss-accelerate.aliyuncs.com/" + urlPath

		}
		response, err := client.Get(reqURL)
		if err != nil {
			// log.Println("HTTP GET 错误：", urlPath, err)
			continue // 发生错误，进行下一次重试
		}

		if response.StatusCode != 200 {
			// log.Println("返回错误状态码 错误", nreqUrl, response.StatusCode)
			response.Body.Close()
			continue // 返回错误状态码，进行下一次重试
		}

		file, err := os.Create(fileRealPath)
		if err != nil {
			log.Println("os Create 错误：", err, reqURL)
			response.Body.Close()
			return err
		}

		_, err = io.Copy(file, response.Body)
		file.Close()
		response.Body.Close()
		if err != nil {
			// log.Println("io.Copy 错误：", reqURL, err)
			os.Remove(fileRealPath)
			continue // 发生复制错误，进行下一次重试
		}

		return nil // 成功下载文件，返回 nil
	}
	log.Println("下载失败", reqURL)
	return errors.New("下载文件失败") // 重试多次后仍然失败，返回错误
}

func DownloadFile2(urlPath string, wg *sync.WaitGroup) error {
	defer wg.Done()

	// 判断文件是否存在
	fileRealPath := filepath.Join("/data", urlPath)
	if fsInfo, err := os.Stat(fileRealPath); !os.IsNotExist(err) {
		// 如果文件存在。
		if fsInfo.Size() > 0 {
			return nil
		}
	} else { // 避免文件夹不存在。
		err := os.MkdirAll(filepath.Dir(fileRealPath), 0755)
		if err != nil {
			log.Println("创建文件夹错误：", filepath.Dir(fileRealPath), err)
			return err
		}
	}

	// 发起 HTTP GET 请求获取文件
	reqURL := "http://source.hnwwa.com/" + urlPath
	client := http.Client{
		Timeout: time.Second * 30, // 设置超时时间为30秒
	}

	response, err := client.Get(reqURL)
	if err != nil {
		log.Println("HTTP GET 错误：", urlPath, err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println("返回错误状态码 错误", reqURL, response.StatusCode)
		return err
	}

	file, err := os.Create(fileRealPath)
	if err != nil {
		log.Println("os Create 错误：", err, reqURL)
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		log.Println("io.Copy 错误：", fileRealPath, err)
		defer os.Remove(fileRealPath)
		return err
	}

	return nil
}

func DownloadM3U8(urlPath string, wg *sync.WaitGroup) error {
	defer wg.Done()

	wg.Add(1)
	DownloadFile(urlPath, wg)

	fileContent, err := ioutil.ReadFile(filepath.Join("/data", urlPath))
	if err != nil {
		log.Println("读取文件失败:", urlPath, err)
		return err
	}

	// 将字节切片转换为字符串
	bodyStr := string(fileContent)
	// 下载ts文件
	re := regexp.MustCompile(`index-.*\.ts`)
	matches := re.FindAllString(bodyStr, -1)
	var wgts sync.WaitGroup
	var wgtsCount = 0
	for _, match := range matches {
		tsPath := filepath.Join(filepath.Dir(urlPath), match)
		wgts.Add(1)
		wgtsCount += 1
		go DownloadFile(tsPath, &wgts)
		if wgtsCount%100 == 0 {
			wgts.Wait()
		}

	}
	wgts.Wait()
	return nil
}

func syncVod(stime, etime string) {
	var wg sync.WaitGroup
	var wgcount = 0

	var vods []vod.VodList
	model.DataBase.Model(vod.VodList{}).Where("job_status = 3 and created_at > ? and created_at < ?", stime, etime).Limit(1000).Select("id, cover, hls_path").Order("id asc").Find(&vods)
	log.Println("开始同步vod数据", stime, etime, len(vods))
	for _, vod := range vods {
		wgcount += 1
		if vod.Cover != "" {
			wg.Add(1)
			go DownloadFile(vod.Cover, &wg)
		}
		if vod.HlsPath != "" {
			wg.Add(1)
			go DownloadM3U8(vod.HlsPath, &wg)
		}
		if wgcount%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func syncVlog(stime, etime string) {
	var wg sync.WaitGroup
	var wgcount = 0
	var vlogs []vlog.VlogList
	model.DataBase.Model(vlog.VlogList{}).Where("job_status = 3 and created_at > ? and created_at <= ?", stime, etime).Find(&vlogs)
	log.Println("开始同步vlog数据", stime, etime, len(vlogs))
	for _, vlog := range vlogs {
		wgcount += 1
		if vlog.Cover != "" {
			wg.Add(1)
			go DownloadFile(vlog.Cover, &wg)
		}
		if vlog.HlsPath != "" {
			wg.Add(1)
			go DownloadM3U8(vlog.HlsPath, &wg)
		}
		// 避免耗尽资源
		if wgcount%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func syncBlog(stime, etime string) {
	var wg sync.WaitGroup
	var wgcount = 0
	var blogs []blog.BlogImage
	model.DataBase.Model(blog.BlogImage{}).Where("created_at > ? and created_at <= ?", stime, etime).Find(&blogs)
	log.Println("开始同步blog数据", stime, etime, len(blogs))
	for _, blog := range blogs { // 扫描当前行的字段值到结构体中
		wgcount += 1
		wg.Add(1)
		go DownloadFile(blog.Path, &wg)
		// 避免耗尽资源
		if wgcount%100 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func main() {
	stime := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 00:00:00"
	etime := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 23:59:59"
	syncVod(stime, etime)
	syncVlog(stime, etime)
	syncBlog(stime, etime)
}
