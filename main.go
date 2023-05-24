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
)

type Response struct {
	Result []struct {
		Hostname string `json:"hostname"`
		SSL struct {
			Cert []struct {
				Expires string `json:"expires_on"`
				Issued string `json:"issued_on"`
			} `json:"certificates"`
		} `json:"ssl"`
	} `json:"result"`
}

func PrettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
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

	//sb := string(body)
	//log.Printf(sb)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	//fmt.Println(PrettyPrint(result.Result))
	for _, rec := range result.Result {
		for _, cer := range rec.SSL.Cert {
			fmt.Print(rec.Hostname)
			fmt.Println(rec.Hostname, cer.Issued, cer.Expires)
		}
	}

	outputFile, err := os.Create("certs.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"Hostname", "Issued Date", "Expiration Date"}
	if err := writer.Write(header); err != nil {
		log.Fatalln(err)
	}

	for _, rec := range result.Result {
		var csvRow []string
		csvRow = append(csvRow, rec.Hostname)
		for _, cer := range rec.SSL.Cert {
			csvRow = append(csvRow, cer.Issued, cer.Expires)
		}
		if err := writer.Write(csvRow); err != nil {
			log.Fatalln(err)
		}
	}
}