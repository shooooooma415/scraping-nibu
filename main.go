package main

import (
	"fmt"
	"log"
	"sync"
	"scraping-nibu/usecase"
)

const baseURL = "https://www.hinatazaka46.com"
const totalPages = 12 // ページ数

func main() {
	var wg sync.WaitGroup

	// 34ページを並列処理
	for i := 11; i < totalPages; i++ {
		wg.Add(1) // Goroutine を追加
		go func(pageNum int) {
			defer wg.Done() // Goroutine が完了したらカウントを減らす

			url := fmt.Sprintf("https://www.hinatazaka46.com/s/official/diary/member/list?ima=0000&page=%d&ct=16&cd=member", pageNum)
			htmlContent := usecase.ScrapeWebsite(url)

			if htmlContent == "" {
				log.Printf("Failed to fetch webpage content for page %d", pageNum)
				return
			}

			htmlContent = usecase.FixRelativeURLs(htmlContent, baseURL)
			htmlContent = usecase.AddBaseTag(htmlContent, baseURL)

			fileName := fmt.Sprintf("blog_%02d.pdf", pageNum) // 2桁のページ番号をファイル名にする
			usecase.SaveToPDF(htmlContent, fileName)
			fmt.Printf("PDF saved: %s\n", fileName)
		}(i) // i を Goroutine に渡す
	}

	wg.Wait() // すべての Goroutine の処理を待つ
	fmt.Println("All PDFs saved successfully!")
}