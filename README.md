# Forex Rates Fetcher for ITR filing in India for EUR/INR

For the ITR filing, if one holds Foreign Equity, one might probably need to calculate the cost of the acquired shares in INR, on that day. 
Income Tax Department of India, mentions the reference conversion rate should be the TT Buy rate as mentioned by SBI on the day of purchase (or sale).

How do we get the conversion rate on the day of purchase, which is usually in the past. There are few good Samaritans who get the daily SBI forex rates PDF and store in github, 
like: https://github.com/skbly7/sbi-tt-rates-historical/tree/master and https://www.sahil-gupta.in/posts/historical-sbi-forex-rates/

## What does this library do?
Given a list of dates, this code fetches you the TT Buy rate of EUR/INR for those dates. 
Where does this library fetch the historical data from: https://officialforexrates.com/ 

## How to run the code:
To run the code, one would need go environment, install golang in your system.
Input needed: create a csv file, with list of dates that you want to get the TT Buy rate for, and update the path in the main.go file (line number 29 - as of writing of this Readme file) also change the output file name, 
line number 218 - as of writing of this Readme file

run the code using command : `go run main.go`

## Cases this code helps:
1. Helps to get the historical TT Buy rate of EUR/INR
     i. which indirectly helps in calculating the cost of acquistion (especially if one buys the shares every month or SIPs)
     ii. Also this same code can be used to calculating the cost of equity held.
Ofcourse, you would need the help of excel to help calculate the total cost of shares (=SUM(=PRODUCT(quantity * cost of share * TT Buy Rate)))

## Limitations:
For now, i have only coded for what helped me, i only needed EUR/INR and logic does that only. There is lot of scope for improvements and enhancements... will try to update the code whenever possible or when a need arises.
