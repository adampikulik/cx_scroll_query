# cx_scroll_query
 
## Description
---
I created the cx_scroll_query tool to retrieve logs from Coralogix teams faster and easier. It executes HTTP requests to Coralogix ES API (https://coralogix.com/tutorials/elastic-api/)


## List of files
---
 config.txt - Set your ES key and the cluster once.<br />
    \t ES_KEY - It is your ES key which you can find at https://\<your team url\>/#/settings/account/api_key. 
    \t CLUSTER - It can be EU (coralogix.com), IN (app.coralogix.in), US (app.coralogix.us). By default: EU.
 
 cx_scroll_query - the tool compiled on MacOS
 main.go - the source file
 query.txt - a file with ES query. It contains an example
