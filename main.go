package main

import (
	"log"
	"spider/database/mysql"
	"spider/database/redis"
	"spider/extract"
	"sync"
)

const (
	MAXGOROUTINES = 4
	RedisKey      = "URLKEY"
)

func main() {
	contents := []extract.Content{
		extract.Content{"http://www.zhihu.com", "http://www.zhihu.com"},
		extract.Content{"http://www.baidu.com", "http://www.baidu.com"},
		extract.Content{"http://www.bilibili.com", "http://www.bilibili.com"},
	}
	receive := make(chan []extract.Content, 10)
	// wg 控制所有处理goroutine的退出，
	// 如controller、idetifier和saver
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {

		controller(contents, receive)
		wg.Done()
	}()

	go func() {
		identifier(receive)
	}()
	wg.Wait()
	log.Println("All controller is done")
	close(receive)

}

// identifier 负责将传入的数据分类鉴别，如果在Redis已经存在的，就舍弃掉
// 如果不存在，就储存在Redis中和MySQL中，传入saver中
func identifier(receive <-chan []extract.Content) {
	save := make(chan extract.Content, 10)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		saver(save)
		wg.Done()
	}()

	for cs := range receive {
		for _, c := range cs {
			ok, err := redis.HKExists(RedisKey, c.URL)
			if err != nil {
				log.Println(err)
				continue
			}
			if ok {
				log.Println(c.URL, "was reached already!")
			} else {
				log.Println(c.URL, "is not reached")
				save <- c
			}
		}
	}
	wg.Wait()
	close(save)
}

// saver 负责将上层中的数据存入MySQL和Redis数据库
func saver(data <-chan extract.Content) {
	for c := range data {
		redis.HSetNX(RedisKey, c.URL)
		mysql.Save(c.URL, "")
		// fmt.Println(c.URL)
	}
}

// controller 负责向worker中分配任务，一个controller中最多可以运行MAXGOROUTINES个worker
// controller 把上一层的receive传给worker，使worker直接把结果传给上一层
// controller等待所有子worker 都退出之后再退出
func controller(contents []extract.Content, receive chan []extract.Content) {
	worktoken := make(chan struct{}, MAXGOROUTINES)
	var wg sync.WaitGroup

	for _, content := range contents {
		// wg.Add必须放在go外面，不然会导致直接退出
		wg.Add(1)
		go func(content extract.Content, receive chan []extract.Content) {

			worktoken <- struct{}{}
			worker(content, receive)
			<-worktoken
			wg.Done()
		}(content, receive)
	}
	wg.Wait()
	close(worktoken)
}

// worker 负责网页中链接的提取，把结果传入到data中
func worker(content extract.Content, data chan<- []extract.Content) {
	log.Println("Worker processing...")
	cs, err := extract.Extract(content)
	if err != nil {
		return
	}
	// fmt.Println(cs)
	data <- cs
	log.Println("Worker procceed")
}
