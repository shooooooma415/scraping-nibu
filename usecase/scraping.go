package usecase

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gocolly/colly"
)



func ScrapeWebsite(url string) string {
	c := colly.NewCollector()

	var htmlContent string

	c.OnHTML("html", func(e *colly.HTMLElement) {
		var err error
		htmlContent, err = e.DOM.Html()
		if err != nil {
			log.Println("Error getting HTML:", err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	return htmlContent
}

// 相対パスの画像やリンクを完全なURLに変換
func FixRelativeURLs(html, baseURL string) string {
	// 画像URLを修正
	imgRegex := regexp.MustCompile(`(<img[^>]+src=["'])(/[^"']+)`)
	html = imgRegex.ReplaceAllString(html, `$1`+baseURL+`$2`)

	// CSSやJSのURLを修正
	linkRegex := regexp.MustCompile(`(<link[^>]+href=["'])(/[^"']+)`)
	html = linkRegex.ReplaceAllString(html, `$1`+baseURL+`$2`)

	scriptRegex := regexp.MustCompile(`(<script[^>]+src=["'])(/[^"']+)`)
	html = scriptRegex.ReplaceAllString(html, `$1`+baseURL+`$2`)

	// 内部リンクを修正
	aRegex := regexp.MustCompile(`(<a[^>]+href=["'])(/[^"']+)`)
	html = aRegex.ReplaceAllString(html, `$1`+baseURL+`$2`)

	return html
}

// HTMLに <base href="..."> を追加して、相対パスを自動補正
func AddBaseTag(html, baseURL string) string {
	baseTag := `<base href="` + baseURL + `">`
	if strings.Contains(html, "<head>") {
		return strings.Replace(html, "<head>", "<head>\n"+baseTag, 1)
	} else {
		return "<head>" + baseTag + "</head>\n" + html
	}
}

// HTMLをPDFに変換して保存
func SaveToPDF(htmlContent, outputFilename string) {
	// PDF生成ツールのインスタンス作成
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// ページ設定
	page := wkhtmltopdf.NewPageReader(strings.NewReader(htmlContent))
	page.EnableLocalFileAccess.Set(true) // ローカルファイルのアクセスを有効化
	page.NoStopSlowScripts.Set(false)    // JavaScript 実行を許可
	page.JavascriptDelay.Set(2000)       // JS遅延を追加（必要に応じて調整）

	// PDFへ追加
	pdfg.AddPage(page)

	// PDF作成
	err = pdfg.Create()
	if err != nil {
		log.Fatal("Failed to create PDF:", err)
	}

	// ファイルとして保存
	err = os.WriteFile(outputFilename, pdfg.Bytes(), 0644)
	if err != nil {
		log.Fatal("Failed to save PDF:", err)
	}
}
