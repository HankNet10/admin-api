package ip

import (
	"log"
)

// Result 归属地信息
type Result struct {
	IP      string
	Country string
	Area    string
}

var qqWry QQwry

func init() {
	IPData.FilePath = "./resource/qqwry.dat"
	// startTime := time.Now().UnixNano()
	res := IPData.InitIPData()
	if v, ok := res.(error); ok {
		log.Panic(v)
	}
	// endTime := time.Now().UnixNano()
	//log.Printf("IP 库加载完成 共加载:%d 条 IP 记录, 所花时间:%.1f ms\n", IPData.IPNum, float64(endTime-startTime)/1000000)
	qqWry = NewQQwry()
}

func Find(ip string) Result {

	return qqWry.Find(ip)
}
