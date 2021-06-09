package tool

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var url string
var ImgPath = "/images/meizitu/"

var Header = map[string]string{
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36",
	"referer":    "https://www.mzitu.com/page/",
}
var MzituHeader = http.Header{
	"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36"},
	"referer":    []string{"https://www.mzitu.com/page/"},
}

func SearchMzitu(kw string) []string {
	url = fmt.Sprintf("https://www.mzitu.com/search/%s", kw)
	resp, err := Get(url, Header)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	list := doc.Find("#pins").Find("span a")
	if list.Length() == 0 {
		return nil
	}
	index := int(RandInt(0, list.Length()))
	var result []string
	list.Each(func(i int, s *goquery.Selection) {
		if i != index {
			return
		}
		detailUrl, _ := s.Attr("href")
		result = getDetail(detailUrl)
		return
	})
	if result == nil {
		list.Each(func(i int, s *goquery.Selection) {
			if result != nil {
				return
			}
			detailUrl, _ := s.Attr("href")
			result = getDetail(detailUrl)
			return
		})
	}
	result = SliceUnique(result)
	return result
}

func getDetail(url string) []string {
	resp, err := Get(url, Header)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	img := doc.Find(".main-image img")
	firstImg, _ := img.Attr("src")
	lastDoc := doc.Find(".pagenavi a")
	lastIndex := lastDoc.Length() - 2
	totalCount := 0
	lastDoc.Each(func(i int, s *goquery.Selection) {
		if i == lastIndex {
			lastUrl, _ := s.Attr("href")
			urlArr := strings.Split(lastUrl, "/")
			totalCount, _ = strconv.Atoi(urlArr[len(urlArr)-1])
		}
	})
	var images []string
	for i := 1; i <= totalCount; i++ {
		imgUrl := strings.Replace(firstImg, "01.jpg", fmt.Sprintf("%02d", i)+".jpg", 1)
		images = append(images, imgUrl)
	}
	max := totalCount
	if totalCount >= 10 {
		max -= 10
	}
	index := RandInt(0, max)
	result := images[index : index+10]
	return result
}
