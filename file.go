/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : file.go

* Purpose :

* Creation Date : 11-25-2013

* Last Modified : Sun 29 Dec 2013 09:09:26 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package mssh

import (
	"bufio"
	"os"
	"strings"
)

func GetHostsByFile(file string) []string {
	var hosts []string
	fr, _ := os.Open(file)
	defer fr.Close()
	r := bufio.NewReader(fr)
	for {
		readline, _ := r.ReadString('\n')
		if !strings.Contains(readline, "#") {
			if readline != "" {
				hosts = append(hosts, strings.Trim(readline, "\n"))
			}
		}
		if readline == "" {
			break
		}
	}
	return hosts
}
