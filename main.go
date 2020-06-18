package main

import (
	"fmt"

	"github.com/figassis/hnfaves/pkg/api"

	"github.com/figassis/hnfaves/pkg/utl/config"
)

func main() {

	cfg, err := config.Load()
	checkErr(err)
	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
}
