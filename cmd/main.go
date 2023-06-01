package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type OrganDohodov struct {
	Name        string
	Beneficiars []*Beneficiar
	Size        int
}
type Beneficiar struct {
	Code int
	BIN  int
	Name string
}

func main() {
	resp, err := http.DefaultClient.Get("https://egov.kz/cms/ru/articles/taxes_bin")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sum := 0
	bufalloReader := bufio.NewScanner(resp.Body)
	binGosDoohodov := []*OrganDohodov{}
	leng := -1
	queue := 0
	for bufalloReader.Scan() {
		run := bufalloReader.Text()
		func() {
			a, b, _ := strings.Cut(run, "&nbsp;")
			run = a + b
		}()
		fmt.Println(run)
		if strings.Contains(run, "<td></td>") {
			continue
		}
		if strings.Contains(run, "td") && leng >= 0 {
			switch queue {
			case 0:
				queue++
				binGosDoohodov[leng].Beneficiars = append(binGosDoohodov[leng].Beneficiars, &Beneficiar{})
				binGosDoohodov[leng].Beneficiars[binGosDoohodov[leng].Size].Code, err = strconv.Atoi(run[7 : len(run)-5])
				if err != nil {
					fmt.Println(err.Error())
				}
				binGosDoohodov[leng].Size++
				fmt.Println("blyaha")
			case 1:
				binGosDoohodov[leng].Beneficiars[binGosDoohodov[leng].Size-1].BIN, err = strconv.Atoi(run[7 : len(run)-5])
				if err != nil {
					fmt.Println(err.Error())
				}
				queue++
			case 2:
				binGosDoohodov[leng].Beneficiars[binGosDoohodov[leng].Size-1].Name = run[7 : len(run)-5]
				queue = 0
			}
		}
		if strings.Contains(run, "slidedown-title toggle") {
			sum++
			// fmt.Println(run[35 : len(run)-5])
			binGosDoohodov = append(binGosDoohodov, &OrganDohodov{Name: run[35 : len(run)-5]})
			leng++
			counter := 0
			for bufalloReader.Scan() && counter < 2 {
				if run = bufalloReader.Text(); strings.Contains(run, "tr>") {
					fmt.Println(run)
					counter++
				} else {
					fmt.Println(run)
				}

			}
			// fmt.Println("totottootot")
		}
	}
	// for run, _, err := bufalloReader.ReadLine(); err == nil; {
	// 	sum++
	// 	fmt.Print(string(run))
	// }
	fmt.Println(sum)
	str, _ := json.Marshal(binGosDoohodov)
	f, err := os.Create("./AllBins.js")
	check(err)
	defer f.Close()
	n2, err := f.Write(str)
	check(err)
	fmt.Println(string(str), n2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
