package tool

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

const NvShenDomain = "https://www.invshen.net"

var InvShenHeader = map[string]string{
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36",
	"referer":    "https://www.invshen.net/",
}
var InvShenHeader2 = http.Header{
	"User-Agent":       []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"},
	"referer":          []string{"https://www.invshen.net/"},
	"sec-ch-ua":        []string{"\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\""},
	"sec-ch-ua-mobile": []string{"?0"},
}

func SearchNvShen(kw string) []string {
	url := NvShenDomain + "/gallery/"
	resp, err := Get(url, InvShenHeader)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	list := doc.Find(".tag_div li a")
	//var result []string
	listUrl := ""
	list.Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		name := s.Text()
		if strings.Contains(name, kw) {
			listUrl = NvShenDomain + href
		}
		return
	})
	if listUrl == "" {
		return nil
	}
	return filterList(listUrl)
}

func filterList(url string) []string {
	resp, err := Get(url, InvShenHeader)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	list := doc.Find(".galleryli_link img")
	var firstImg string
	var detailUrl string
	var images []string
	index := int(RandInt(0, list.Length()))
	for _, attr := range list.Nodes[index].Attr {
		if attr.Key == "data-original" {
			images = append(images, attr.Val)
			firstImg = attr.Val
		}
	}
	for _, attr := range doc.Find(".galleryli_link").Nodes[index].Attr {
		if attr.Key == "href" {
			detailUrl = NvShenDomain + attr.Val
		}
	}
	totalCount := getTotalCount(detailUrl)
	for i := 1; i <= totalCount; i++ {
		imgUrl := strings.Replace(firstImg, "cover/0.jpg", "s/"+fmt.Sprintf("%03d", i)+".jpg", 1)
		images = append(images, imgUrl)
	}
	max := totalCount
	if totalCount >= 10 {
		max -= 10
	}
	randIndex := RandInt(0, max)
	result := images[randIndex : randIndex+10]
	return result
}

func getTotalCount(url string) int {
	resp, err := Get(url, InvShenHeader)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	list := doc.Find("#dinfo span")
	var total int
	list.Each(func(i int, s *goquery.Selection) {
		text := strings.ReplaceAll(s.Text(), " ", "")
		count := strings.ReplaceAll(text, "张照片", "")
		total, _ = strconv.Atoi(count)
		return
	})
	return total
}
