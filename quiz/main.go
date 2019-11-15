package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	// Open the file
	csvfile, err := os.Open("problems.csv")
	trueCount := 0;
	falseCount := 0;
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Question : ", record[0]+"\n")
			text, _ := reader.ReadString('\n')

			if record[1] == strings.TrimRight(text, "\n") {
				trueCount++
			} else {
				falseCount++;
			}

		}

		wg.Done()
	}()

	go func() {
		time.Sleep(13 * time.Second)
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("%s%d%s","You know " ,trueCount ,"question on 13 question" )

}
