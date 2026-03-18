package controller

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 从全局数据源中随机抽取指定数量的图片数据，组装成一个首页数据结构，然后返回给前端
// JSON 响应。通常用于网站的首页推荐流，目的是让每次刷新页面时展示的内容都不一样，
// 具有随机性的同时控制返回的数据量。

// home 定义返回给前端的JSON结构
type home struct {
	Banners []string `json:"banners"` // 轮播图列表
	Images1 []string `json:"images1"` // 第一组图片
	Images2 []string `json:"images2"` // 第二组图片
}

// bannerDataList 定义好的全局轮播图数据源
var bannerDataList = []string{
	"tes1",
	"test2",
	"test3",
}

// imgDataList 定义好的全局图片数据源
var imgDataList = []string{
	"tes1",
	"test2",
	"test2",
	"tes1",
	"test2",
	"test2",
	"tes1",
	"test2",
	"test2",
	"tes1",
	"test2",
	"test2",
	"tes1",
	"test2",
	"test2",
	"tes1",
	"test2",
	"test2",
	"tes1",
	"test2",
	"test2",
	"tes1",
	"test2",
	"test2",
}

// Home 组装首页数据，通过gin返回给前端
func (*Controller) Home(c *gin.Context) {
	// 1. 处理轮播图片
	bannerNum := 3 // 抽取三张作为轮播图
	// 1.1 生成索引列表
	indexList := make([]int, len(bannerDataList))
	for i := range bannerDataList {
		indexList[i] = i
	}
	// 1.2 随机打乱数据源并截取前bannerNum个索引
	list := randList(indexList, bannerNum)
	// 1.3 根据随机索引抽取具体的图片URL
	bannerList := make([]string, bannerNum)
	for i := range list {
		bannerList[i] = bannerDataList[i]
	}
	// 2. 处理推荐图片
	imgNum := 10 // 抽取十张作为图片
	// 2.1 重置indexList，生成索引列表
	indexList = make([]int, len(imgDataList))
	for i := range imgDataList {
		indexList[i] = i
	}
	// 2.2 随机打乱数据源并截取前imgNum个索引
	list = randList(indexList, imgNum)
	// 2.3 根据随机索引抽取具体的图片URL
	imgList := make([]string, imgNum)
	for i := range list {
		imgList[i] = imgDataList[i]
	}
	// 3. 数据组装与拆分
	h := &home{}
	h.Banners = bannerList
	h.Images1 = imgList[:5]
	h.Images2 = imgList[5:]
	// 4. 返回响应
	c.JSON(http.StatusOK, h)
}

// randList 打乱顺序逻辑
func randList(indexList []int, num int) []int {
	list := make([]int, num)
	for i := 0; i < num; i++ {
		l := len(indexList)
		index := rand.Intn(l)
		list[i] = indexList[index]
		indexList = append(indexList[:index], indexList[index+1:]...)
	}
	return list
}
