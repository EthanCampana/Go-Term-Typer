package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GenerateWords() []string {
	r_url := "https://random-word-api.herokuapp.com/word?number=600"
	res, err := http.Get(r_url)
	if err != nil {
		fmt.Printf("err = %s \n Exiting Application", err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("err = %s \n Exiting Application", err)
		os.Exit(1)
	}
	sbody := string(body)
	sbody = strings.Replace(sbody, "[", "", -1)
	sbody = strings.Replace(sbody, "]", "", -1)
	sbody = strings.Replace(sbody, "\"", "", -1)
	sbody = strings.Replace(sbody, ",", ", ,", -1)

	return strings.Split(sbody, ",")
}
