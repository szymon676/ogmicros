# This project uses api from http://github.com/szymon676/stocks-api

## How it works? It fetches data from the stock api and watches if its recommendation is changing. If so, it sends this data to the processing-ms, which then sends an SMS and saves the time when it happened to Excel.

## Running:
- run the aggregation-ms by : go run main.go --stockname=stock name of your choice (NADASQ)
- run the processing-ms by : go run main.go, but if you want to have access to sending sms you will need to provide twilio auth dependencies.

### Faq:
- What indicator does it use? 
one minute indicator and refreshes after every 5 seconds