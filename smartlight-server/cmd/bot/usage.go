package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"internal/usage"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const usagesFile = "usages.txt"
const timeFormat = time.RFC3339

func toUsage(str string) usage.Usage {
	var u usage.Usage

	split := strings.Split(str, "|")
	t, _ := time.Parse(timeFormat, split[1])

	u.LampID = split[0]
	u.Time = t
	u.TurnOn, _ = strconv.ParseBool(split[2])

	return u
}

func retrieveStoredUsage() ([]usage.Usage, error) {
	f, err := os.Open(usagesFile)
	if err != nil {
		log.Printf("failed to retrieve stored usage: %v", err)
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var usages []usage.Usage
	for scanner.Scan() {
		usages = append(usages, toUsage(scanner.Text()))
	}

	return usages, nil
}

func sendUsages() {
	resp, err := http.Get(fmt.Sprintf("http://%s:10000/ping", serverIP))
	if err != nil {
		log.Printf("failed to ping server: %v", err)
		return
	}
	defer resp.Body.Close()

	usages, err := retrieveStoredUsage()
	if err != nil || usages == nil {
		return
	}
	os.Remove(usagesFile)

	for _, u := range usages {
		sendUsage(u)
	}
}

func sendUsage(u usage.Usage) {
	log.Printf("sending %+v", u)
	body, err := json.Marshal(u)
	if err != nil {
		log.Printf("failed to marshal usage: %v", err)
		storeUsage(u)
		return
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:10000/usage/add", serverIP), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("failed to send usage: %v", err)
		storeUsage(u)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Print("failed to insert usage")
		storeUsage(u)
	}
}

func storeUsage(u usage.Usage) error {
	f, err := os.OpenFile(usagesFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("failed to open usage file: %v", err)
		return err
	}
	defer f.Close()

	str := fmt.Sprintf("%s|%s|%t\n", u.LampID, u.Time.Format(timeFormat), u.TurnOn)
	_, err = f.Write([]byte(str))
	if err != nil {
		log.Printf("failed to write usage: %v", err)
		return err
	}

	return nil
}
