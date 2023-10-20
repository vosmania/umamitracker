package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	Title           string
	Description     string
	DiscountPercent int
}

func main() {
	// Replace 'URL' with the URL of the shop website containing the discount information.
	url := "https://umami.ee/soodus-2/"

	// Send an HTTP GET request to the website.
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Check if the request was successful (status code 200).
	if response.StatusCode != 200 {
		log.Fatalf("Failed to retrieve the webpage. Status code: %d", response.StatusCode)
	}

	// Parse the HTML content of the page using goquery.
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Create a slice to store product information.
	var products []Product

	// Find all product containers with class "white-bg."
	doc.Find("div.white-bg").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h4").Text()
		description := s.Find("p").Text()
		discountText := s.Find("span.discount_percentage").Text()

		// Trim whitespace and percent sign.
		discountPercentStr := strings.Trim(discountText, " %")

		// Check if discount percentage is empty or not a valid number.
		if discountPercentStr == "" {
			log.Printf("Discount percentage is empty for product #%d: %s\n", i, title)
			return
		}

		discountPercent, err := strconv.Atoi(discountPercentStr)
		if err != nil {
			log.Printf("Error parsing discount percentage for product #%d: %v\n", i, err)
			return
		}

		// Append product information to the slice.
		products = append(products, Product{
			Title:           title,
			Description:     description,
			DiscountPercent: discountPercent,
		})
	})

	// Sort products by discount percentage in descending order.
	sort.Slice(products, func(i, j int) bool {
		return products[i].DiscountPercent > products[j].DiscountPercent
	})

	// Display sorted product information.
	for _, product := range products {
		fmt.Printf("Title: %s\n", product.Title)
		fmt.Printf("Description: %s\n", product.Description)
		fmt.Printf("Discount Percentage: %d%%\n", product.DiscountPercent)
		fmt.Println("-------------------------")
	}
}
