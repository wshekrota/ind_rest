package main

import "fmt"
import "log"
import "net/http"
import "encoding/json"
import "io/ioutil"
import "strings"
import "strconv"
import "os"

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to Walts rest endpoints!")
}

func handleRequests() {

// call ind functions for this endpoint
    http.HandleFunc("/distribution", aggregate)

// wait on port 10000
    log.Fatal(http.ListenAndServe(":10000", nil))
}

// Endpoint of an endpoint  
// Set hostname of elasticsearch instance to use globally using  
// RUN_HOST environment variable 
// Elasticstack endpoint _cat/indices produces an index report.
// Will accrue all indices beginning with 'log-'. You will get 3 
// integers for each distinct log type. (excluding date) This should facilitate 
// setting up a progress bar and/or talking to something like checkmk.
// The resultant hask keys can then articulate where the problem 
// might be. The 3 integers represent totals [red, yellow, green].
// What is progress?
// x x 0 progress to yellow
// 0 x x progress to green
// 0 x 0 new status I discovered stuck yellow when a node issue
// 0 0 x fully green
//
// The schema for the output from that endpoint is..
// health		string
// status		string
// index		string
// uuid			string
// pri			int
// rep			int
// docs.count		int
// docs.deleted		int
// store.size		string
// pri.store.size	string
//
// The real question here is to guage progress on doc count or byte count.
// Currently this is coded with doc count. If you were to do byte would 
// require conversion of that field.
// Modified return value to serialized map ie json

func updateMap(v *map[string][3]int, health string, key string, size int) {

      	hash := *v

	// health values {red,yellow,green}
	// increment the appropriate one or set if new key
        //
	health_map := map[string]int {
		"red": 0,
		"yellow": 1,
		"green": 2,
	}

	// does key exist
	_, ok2 := hash[key]

	if ok2 {
	// update existing key array
		a := hash[key] 
		a[health_map[health]] += size
		hash[key] = a
	} else {

	// initialize new key
                hash[key] = [3]int{0,0,0}
		a := hash[key]
		a[health_map[health]] = size
		hash[key] = a

        }
}

func aggregate(w http.ResponseWriter, r *http.Request) {

        // http handling
	// call indices endpoint
	//

        var host = os.Getenv("RUN_HOST")
	var url = fmt.Sprintf("http://%s:9200/_cat/indices",host)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("*error", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(3)
	}

	fmt.Println("rc=", resp.StatusCode)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error", err)
	}

	// all lines convert
        s := string(body)

	// initialize our map
	ind_map := make(map[string][3]int)

	// Process the http response line by line
	// each line
	//
        for _, line := range strings.Split(strings.TrimSuffix(s,"\n"), "\n") {
		// skip line with no health (closed?)
		if(strings.HasPrefix(line," ")) {
			continue
		}
		// words by field
		words := strings.Fields(line)
		// either create hash key or update it key => [red,yellow,green]
		// call update with map,health,key,size
		// which int increments based on health for each unique index
		key := strings.TrimRight(strings.Join(words[2:3]," "),"0123456789-.")
		if(! strings.HasPrefix(key,"log-")) {
			continue
		}
		// Use single word slices
		num, _ := strconv.Atoi(strings.Join(words[6:7]," "))
		updateMap(&ind_map,strings.Join(words[0:1]," "),key,num)
	}

	// now return the map
    serial, err := json.Marshal(ind_map)
    fmt.Fprintf(w,string(serial))

}

func main() {
    fmt.Println("API v1.0 ")
    handleRequests()
}
