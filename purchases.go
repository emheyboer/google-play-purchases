package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// panics on a non-nil err
func check(err error) {
	if err != nil {
		panic(err)
	}
}

type history []struct {
	OrderHistory struct {
		// BillingInstrument interface{}
		// OtherInstruments  interface{}
		// BillingContact    interface{}
		// AssociatedContact interface{}
		// UserLocale        interface{}
		OrderId      string
		CreationTime string
		IpAddress    string
		IpCountry    string
		TotalPrice   string
		Tax          string
		RefundAmount string
		Preorder     bool
		LineItem     []struct {
			Doc struct {
				DocumentType string
				Title        string
			}
			Quantity float64
		}
	}
}

func processFile(file *zip.File) {
	contents, err := file.Open()
	check(err)
	defer contents.Close()

	dec := json.NewDecoder(contents)

	var data history
	dec.Decode(&data)

	fmt.Printf("%10s,%8s,%13s, %s\n", "date", "price", "kind", "name")

	for _, purchase := range data {
		order := purchase.OrderHistory

		if order.TotalPrice != "$0.00" && order.RefundAmount == "$0.00" {
			price, err := strconv.ParseFloat(order.TotalPrice[1:], 64)
			check(err)

			item := ""
			kind := ""
			sep := ""
			for _, lineItem := range order.LineItem {
				item += sep + lineItem.Doc.Title
				kind += sep + lineItem.Doc.DocumentType
				sep = "; "
			}

			date := order.CreationTime
			dateIndex := strings.Index(order.CreationTime, "T")
			if dateIndex != -1 {
				date = order.CreationTime[:dateIndex]
			}

			fmt.Printf("%s,%8s,%13s, %s\n", date, fmt.Sprintf("$%.2f", price), kind, item)
		}
	}
}

func help() {
	name := os.Args[0]
	fmt.Println("a tool for extracting a list of google play purchases")
	fmt.Println("from a google takeout file\n")
	fmt.Println("Usage:")
	fmt.Println(name, "[google takeout zip]")
	fmt.Println("Example:")
	fmt.Println(name, "~/Downloads/takeout.zip")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		help()
	}

	zipFile, err := zip.OpenReader(os.Args[1])
	check(err)
	defer zipFile.Close()

	for _, file := range zipFile.File {
		if file.Name == "Takeout/Google Play Store/Order History.json" {
			processFile(file)
			break
		}
	}
}
