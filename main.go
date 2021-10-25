package main

import (
	"io/ioutil"
  "bufio"
  "io"
	"fmt"
	"strings"
	"os"
	"net/http"
	"encoding/json"
  "strconv"
  "math"
  "path/filepath"
)

// Load configuration from config.txt file to variables
var config, err = ReadConfig(`config.txt`)

func loadQuery() ( query string, size int) {
  file, err := os.Open("query.txt")
  if err != nil {
    panic(err)
  }
  defer file.Close()
  var strs []string
  buf := make([]byte, 1024)
  for {
    n, err := file.Read(buf)
    //fmt.Println(n, err, buf[:n])
    strs = append(strs, string(buf[:n]))
    //fmt.Println(string(buf[:n]))
    if err == io.EOF {
      break
    }
  }
  loadedquery := (strings.Join(strs, ""))

  var qsize map[string]interface{}
  json.Unmarshal([]byte(loadedquery), &qsize)
  size_str := fmt.Sprint(qsize["size"])
  size, err = strconv.Atoi(size_str)
  if err != nil {
    fmt.Println("Size not specified in the query. The default is 10000")
    size = 0
  }
  //fmt.Printf("size: %v \n", qsize["size"])
  //fmt.Println(strs) - it is in []

  return loadedquery, size
}

type Config map[string]string

func ReadConfig(filename string) (Config, error) {
     // init with some bogus data
  config := Config{
    "CLUSTER":"EU",
    "ES_KEY":"11111111-1111-1111-1111-111111111111",
  }
  if len(filename) == 0 {
    return config, nil
  }
  file, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  
  reader := bufio.NewReader(file)
  
  for {
    line, err := reader.ReadString('\n')
    
    // check if the line has = sign
             // and process the line. Ignore the rest.
    if equal := strings.Index(line, "="); equal >= 0 {
      if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
        value := ""
//        fmt.Println("key: ", key)
        if len(line) > equal {
          value = strings.TrimSpace(line[equal+1:])
        }
                // assign the config map
//                fmt.Println("value: ", value)
        config[key] = value
      }
    }
    if err == io.EOF {
      break
    }
    if err != nil {
      return nil, err
    }
  }
    if os.Getenv("ES_KEY") != "" {
    config["ES_KEY"] = os.Getenv("ES_KEY")
    }

  return config, nil
}

// Delete all files from the directory
func RemoveContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    fmt.Println("Deleting old files if exist.")
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return err
        }
    }
    return nil
}

func saveToFile(data []byte, filename string) {
  // Create a file
  out, err := os.Create(filename)
  if err != nil {
    // panic?
  }
  defer out.Close()

  // Write a data type []byte to the file
  _, err = out.Write(data)
  if err != nil {
    fmt.Println(err)
  //  return
  }
  // Issue a Sync to flush writes to stable storage.
  out.Sync()
}
func SendRequest (method string, url string, payload io.Reader) (data []byte) {
  client := &http.Client {
    }
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    fmt.Println(err)
  //  return
  }
//  fmt.Printf("token qp1: %v \n", config["ES_KEY"])
  req.Header.Add("token", config["ES_KEY"])
  req.Header.Add("Content-Type", "application/json")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
  //  return
  }
  defer res.Body.Close()
//  fmt.Println("body1: ", res.Body)

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
  //  return
  }
  return body
}

func extract_part1(body []byte) (string, int) {

  var results map[string]interface{}
  json.Unmarshal([]byte(body), &results)

  // Extract scroll_id
  scroll_id := fmt.Sprint(results["_scroll_id"])
  fmt.Printf("Found scroll_id: %v \n", scroll_id)

  // Extract number of logs
  hits_value := fmt.Sprint(results["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"])
  fmt.Printf("Found hits (number of logs): %v \n", hits_value)

  hits, err := strconv.Atoi(hits_value)
  if err != nil {
    fmt.Println(hits)
  }

  return scroll_id, hits
}

// Main function
func main() {

  var es_url string
  
  // Check if ES KEY is correct
  if config["ES_KEY"] == "11111111-1111-1111-1111-111111111111" {
    fmt.Println("The default key found. Update the ES KEY parameter in the config file.")
    os.Exit(1)
  } else {
    key_len := len(config["ES_KEY"])
    first4 := config["ES_KEY"][0:3]
    last4  := config["ES_KEY"][key_len-3:]
    if key_len != 36 {
      fmt.Printf("Incorrect ES KEY format. Expected 36 caracters and there is %d characters\n", key_len)
      os.Exit(1)
    } else {
      fmt.Printf("Your ES KEY is %v****-****-****-****-********%v characters\n", first4, last4)
    }
  }

  // Setting the endpoint URL
  switch config["CLUSTER"] {
    case "EU":
      es_url = "coralogix-esapi.coralogix.com"
    case "IN":
      es_url = "es-api.app.coralogix.in"
    case "US":
      es_url = "esapi.coralogix.us"
    default:
      es_url = "coralogix-esapi.coralogix.com"
  }
  fmt.Printf("Your endpoint URL is %v\n", es_url)

  url := fmt.Sprintf("https://%v:9443/*/_search?scroll=5m", es_url)
  method := "POST"

  // Load the query from a file
  loaded_query, size := loadQuery()

  // Create io.Reader from loaded query
  payload := strings.NewReader(loaded_query)

  // Create a directory for results
  path := "query_result_files"
  firt_response := fmt.Sprintf("%v/results_0.out", path)
  if _, err := os.Stat(path); os.IsNotExist(err) {
    err := os.Mkdir(path, 0755)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  }
  err = RemoveContents(path)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
  fmt.Printf("Results will be saved at %v directory\n", path)

  // Send the first request and save body of the response
  response_body := SendRequest (method, url, payload)

  // Save the first response
  saveToFile(response_body, firt_response)

  // Read the scroll id to create next requests
  scroll_id, hits := extract_part1(response_body)

  // Setting the endpoint URL for other logs
  url2 := fmt.Sprintf("https://%v:9443/_search/scroll", es_url)
  fmt.Printf("Your endpoint URL for the rest of scroll requests is %v\n", url2)
  
  // body of next requests


  // Check if there is a need to run next requests
  if hits > size {
    // Check how many more requests are needed
    count_req := int(math.Ceil(float64(hits) / float64(size)))
    fmt.Println("Number of requests: ", count_req)
    for i := 1; i < count_req; i++ {
    //for i := 1; i < 2; i++ {  
      fmt.Printf("Request: %d of %d\n", i, count_req)
      msg := `{"scroll": "5m","scroll_id": "%v"}`
      body := fmt.Sprintf(msg, scroll_id)
      second_payload := strings.NewReader(body)
      response_body := SendRequest (method, url2, second_payload)
      
      saveToFile(response_body, fmt.Sprintf("%v/results_%d.out", path, i))
    }
  }
  fmt.Println("Completed!")
}
