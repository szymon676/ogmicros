# This project uses api from http://github.com/szymon676/stocks-api

## How it works? It fetches data from the TSLA stock (planning to change that) and watches if its recommendation is changing. If so, it sends this data to the processing-ms, which then sends an SMS and saves the time when it happened to Excel.

### Faq:
- What indicator does it use? 
one minute indicator and refreshes after every 5 seconds