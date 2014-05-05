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

// read a list of hosts to run command on
func GetHostsByFile(file string) []string {
	var hosts []string
	fr, _ := os.Open(file)
	defer fr.Close()
	r := bufio.NewReader(fr)
	for {
		readline, _ := r.ReadString('\n')
		// ignore 'comment' lines
		if !strings.Contains(readline, "#") {
			// add host to list
			if readline != "" {
				hosts = append(hosts, strings.Trim(readline, "\n"))
			}
		}
		// check for end of file
		if readline == "" {
			break
		}
	}
	return hosts
}
