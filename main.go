package main

import (
	"bufio"
	"flag"
	"github.com/Sirupsen/logrus"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type crawlerFunc func()

var (
	isDebug    = flag.Bool("d", false, "run in debug mode")
	crawlerArr = []crawlerFunc{
		// Register your crawler here
		CrawlBundesanzeiger,
	}

	outputFile = "data.csv"
)

func main() {
	flag.Parse()
	if *isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("run in debug mode")
	}

	reader := bufio.NewReader(os.Stdin)

	logrus.Println("Select the crawler you want to run:")
	for i, c := range crawlerArr {
		logrus.Println(i, ":", runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name())
	}
	numStr, _ := reader.ReadString('\n')
	numStr = strings.Trim(numStr, "\n")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		logrus.Fatalln("please enter a number")
	} else if num >= len(crawlerArr) || num < 0 {
		logrus.Fatalln("please choose the correct number")
	}
	f := crawlerArr[num]

	logrus.Println("Enter the output filename:")
	outputFile, _ = reader.ReadString('\n')

	f()
}
