package main

import(
	"io/ioutil"
	"log"
	"fmt"
	"encoding/json"
	"net/http"
	"os"
	"encoding/csv"
	"github.com/joho/godotenv"
	"time"
)

type Response struct {
	Result []struct {
		Hostname string `json:"hostname"`
		SSL struct {
			Cert []struct {
				Expires string `json:"expires_on"`
				Issued string `json:"issued_on"`
				ID string `json:"id"`
			} `json:"certificates"`
		} `json:"ssl"`
	} `json:"result"`
}

func PrettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}

func PrintJSON(str string) {
	const layout = "01-02-2006"
	t := time.Now()
	file := "data/data-" + t.Format(layout) + ".txt"

	f, err := os.Create(file)
    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    _, err2 := f.WriteString(str)
    if err2 != nil {
        log.Fatal(err2)
    }
}

func PrintCSV(result Response) {
	const layout = "01-02-2006"
	t := time.Now()
	file := "certs/cert-" + t.Format(layout) + ".csv"

	outputFile, err := os.Create(file)
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"Hostname", "Issued Date", "Expiration Date", "ID"}
	if err := writer.Write(header); err != nil {
		log.Fatalln(err)
	}

	for _, rec := range result.Result {
		var csvRow []string
		csvRow = append(csvRow, rec.Hostname)
		for _, cer := range rec.SSL.Cert {
			csvRow = append(csvRow, cer.Issued, cer.Expires, cer.ID)
		}
		if err := writer.Write(csvRow); err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}

	client := &http.Client{}

	resp, err := http.NewRequest("GET", os.Getenv("API_URL"), nil)
	if err != nil {
		log.Fatalln(err)
	}

	bearer := "Bearer " + os.Getenv("BEARER_TOKEN")
	resp.Header.Add("Authorization", bearer)
	res, err := client.Do(resp)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)

	PrintJSON(sb)
	
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	//fmt.Println(PrettyPrint(result.Result))
	// for _, rec := range result.Result {
	// 	for _, cer := range rec.SSL.Cert {
	// 		fmt.Print(rec.Hostname)
	// 		fmt.Println(rec.Hostname, cer.Issued, cer.Expires)
	// 	}
	// }

	PrintCSV(result)
}