package main

import (
	"fmt"
	"log"
	"scraping-nibu/usecase"
)

const baseURL = "https://www.hinatazaka46.com"

func main() {
	url := "https://www.hinatazaka46.com/s/official/diary/member/list?ima=0000&page=0&ct=16&cd=member"
	htmlContent := usecase.ScrapeWebsite(url)

	if htmlContent == "" {
		log.Fatal("Failed to fetch webpage content")
	}

	// 画像・リンクのURLを修正
	htmlContent = usecase.FixRelativeURLs(htmlContent, baseURL)

	// `<base>` タグを追加（相対パス対策）
	htmlContent = usecase.AddBaseTag(htmlContent, baseURL)

	usecase.SaveToPDF(htmlContent, "output.pdf")
	fmt.Println("PDF saved successfully!")
}
