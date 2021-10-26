# cx_scroll_query
 
## Description
I created the cx_scroll_query tool to retrieve logs from Coralogix teams faster and easier. It executes HTTP requests to Coralogix ES API (https://coralogix.com/tutorials/elastic-api/)

<br />

## List of files
- config.txt - Set your ES key and the cluster once.<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; - ES_KEY - It is your ES key which you can find at https://\<your team url\>/#/settings/account/api_key.<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; - CLUSTER - It can be EU (coralogix.com), IN (app.coralogix.in), US (app.coralogix.us). By default: EU.<br />
- cx_scroll_query - the tool compiled on MacOS<br />
- cx_scroll_query.exe - the tool compiled on Windows.<br />
- main.go - the source file<br />
- query.txt - a file with ES query. It contains an example<br />

<br />

## Usage
1. Download the following files: cx_scroll_query (or cx_scroll_query.exe), config.txt and query.txt to one directory.
2. Open the terminal on Mac or Command Prompt (cmd) on Windows and change to the directory where you downloaded files.
3. Only on Mac: Add executable permission to the script. Execute:
```
chmod +x cx_scroll_query
```
4. Update query.txt with your Elastic query.
5. Update config.txt with your ES key and cluster.
6. Run the script:<br/>
On Mac:
```
./cx_scroll_query
```
On Windows:
```
cx_scroll_query.exe
```
7. It will create new directory query_result_files and will save files with results there.

<br />

## Development
If you need to compile the tool for your platform
#### Requirements
* ``Go`` version >= 1.17

#### Build
1. Create a directory (for example: cx_scroll_query)
2. Download files.
3. Execute:
```
$ make
```
