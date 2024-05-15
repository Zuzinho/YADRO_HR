package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"yadro/pkg/context"
)

func main() {
	args := os.Args

	var path string
	if len(args) == 2 {
		path = args[1]
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	tableCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}

	scanner.Scan()
	durs := strings.Split(scanner.Text(), " ")
	start, err := time.Parse("15:04", durs[0])
	if err != nil {
		log.Fatal(err)
	}

	end, err := time.Parse("15:04", durs[1])
	if err != nil {
		log.Fatal(err)
	}

	scanner.Scan()
	price, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return
	}

	ctx := context.NewContext(tableCount, price, start, end)
	handler := context.NewHandler()

	fmt.Println(durs[0])
	for scanner.Scan() {
		eventStr := scanner.Text()
		fmt.Println(eventStr)

		eventSplit := strings.Split(eventStr, " ")
		if len(eventSplit) < 3 {
			log.Fatal("no enough args")
		}

		t, err := time.Parse("15:04", eventSplit[0])
		if err != nil {
			log.Fatal(err)
		}

		e, err := handler.Handle(context.NewEvent(eventSplit[1], t, eventSplit[2:]...), ctx)
		if err != nil {
			log.Fatal(err)
		}

		if e != nil {
			fmt.Println(e.String())
		}
	}

	events, err := handler.Close(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range events {
		fmt.Println(e.String())
	}

	fmt.Println(durs[1])

	results := ctx.TablesMoney()
	fmt.Println(strings.Join(results, "\n"))
}
