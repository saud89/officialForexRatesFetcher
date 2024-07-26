package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var urlString = "https://officialforexrates.com/"

func main() {
	sessionToken, authenticityToken, err := fetchSessionAndAuthenticityToken()
	//data, err := os.ReadFile("<sample Testing file>")
	if err != nil {
		fmt.Println(err)
		return
	}

	//for each row in csv, fetch the data from the website
	// format of input file: acquisition-date (in 2013-03-14 format)
	inputFile, err := os.Open("InputFile.csv")
	if err != nil {
		fmt.Println("Error while opening the csv file, error : ", err)
		return
	}
	defer inputFile.Close()

	csvReader := csv.NewReader(inputFile)

	processCSVData(csvReader, *sessionToken, *authenticityToken)

	//fetching the authenticityToken:

}

func fetchSessionAndAuthenticityToken() (*string, *string, error) {

	httpClient := http.Client{}
	req, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		fmt.Println("Error while creating http request")
		return nil, nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error while sending http request")
		return nil, nil, err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response body")
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error while sending http request. Status code is not 200")
		return nil, nil, err
	}

	sessionHeaderRaw := resp.Header.Get("Set-Cookie")
	fmt.Println("Session header is: " + sessionHeaderRaw)
	sessionHeader := extractSessionHeader(sessionHeaderRaw)

	authenticityToken, err := fetchAuthenticityToken(responseBody)
	if err != nil {
		fmt.Println("Error while fetching authenticity token")
		return nil, nil, err
	}
	return &sessionHeader, authenticityToken, nil

}

func extractSessionHeader(sessionHeaderRaw string) string {
	return strings.Split(sessionHeaderRaw, "session=")[1]
}

func fetchAuthenticityToken(htmlData []byte) (*string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlData))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	authenticityToken := doc.Find("input[name='authenticity_token']").AttrOr("value", "")
	fmt.Println("Authenticity token is: " + authenticityToken)
	return &authenticityToken, nil
}

func fetchConversionRate(sessionToken, authenticityToken, date string) (*string, error) {
	//fetch the conversion rate from the website
	//return the conversion rate
	httpClient := http.Client{}
	urlEncodedData := url.Values{}
	urlEncodedData.Set("authenticity_token", authenticityToken)
	urlEncodedData.Set("date", date)

	encodedData := urlEncodedData.Encode()
	req, err := http.NewRequest(http.MethodPost, urlString, strings.NewReader(encodedData))
	if err != nil {
		fmt.Println("Error while creating http request to fetch the conversion rate")
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "session="+sessionToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error while sending http request to fetch the conversion rate")
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response body")
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//fmt.Println("Doc is: " + doc.Text())
	var conversionRateInString string
	table := doc.Find("table.rates-table")
	if table.Length() > 0 {
		fmt.Println("Table found!")
		row := table.Find("tr").Eq(3)
		column := row.Find("td").Eq(0)
		fmt.Println(column.Text())
		conversionRateInString = column.Text()

	} else {
		fmt.Println("Table not found")
	}

	conversionRateToReturn, err := strconv.ParseFloat(conversionRateInString, 64)
	if err != nil {
		fmt.Println("Error while converting conversion rate to float")
		return nil, err
	}
	if conversionRateToReturn <= 0 {
		fmt.Println("Conversion rate is less than or equal to 0")
		return nil, fmt.Errorf("conversion rate is less than or equal to 0")
	}

	return &conversionRateInString, nil
}

func processCSVData(csvReader *csv.Reader, sessionToken string, authenticityToken string) {

	// for {
	// 	row, err := csvReader.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		fmt.Println("Error while reading the csv file")
	// 		return
	// 	}
	// 	fmt.Println(row)
	// 	date := row[0]
	// 	amountInForeignCurrency, err := strconv.ParseFloat(row[1], 64)
	// 	if err != nil {
	// 		fmt.Println("Error while converting amount to float")
	// 		return
	// 	}
	// 	currency := row[2]

	// 	conversionRate, err := fetchConversionRate(sessionToken, authenticityToken, date)
	// 	if err != nil {
	// 		fmt.Println("Error while fetching conversion rate")
	// 		return
	// 	}

	// 	amountInEUR := amountInForeignCurrency / *conversionRate
	// 	fmt.Println("Amount in EUR is: ", amountInEUR)
	// }
	outputToWrite := make([][]string, 0, 100)
	outputToWrite = append(outputToWrite, []string{"acquistion-date", "conversion-rate"})
	index := 0
	for {
		lineData, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if index > 0 {
			purchaseDate := strings.TrimSpace(lineData[0])
			conversionRate, err := fetchConversionRate(sessionToken, authenticityToken, purchaseDate)
			if err != nil {
				fmt.Println("Error while fetching conversion rate for purchase date : ", purchaseDate)
				fmt.Println("still continuing with the next record")
				continue
			}

			outputToWrite = append(outputToWrite, []string{lineData[0], *conversionRate})

		}

		index++

	}

	writeToFile(outputToWrite)
}

func writeToFile(outputData [][]string) {
	f, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("error in creating a file ,err : ", err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for index, dataToWrite := range outputData {
		if err := w.Write(dataToWrite); err != nil {
			fmt.Println("error while writing data at line number : ", index)
		}
	}
}
