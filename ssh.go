/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : ssh.go

* Purpose :

* Creation Date : 11-18-2013

* Last Modified : Sun 29 Dec 2013 09:18:47 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/
package mssh

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
	"sync"
	"time"
)

func strip(v string) string {
	return strings.TrimSpace(strings.Trim(v, "\n"))
}

type keychain struct {
	keys []ssh.Signer
}

func (k *keychain) Key(i int) (ssh.PublicKey, error) {
	if i < 0 || i >= len(k.keys) {
		return nil, nil
	}
	return k.keys[i].PublicKey(), nil
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	return k.keys[i].Sign(rand, data)
}

func (k *keychain) add(key ssh.Signer) {
	k.keys = append(k.keys, key)
}

func (k *keychain) loadPEM(file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return err
	}
	k.add(key)
	return nil
}

type Conf struct {
	User       string
	Key        string
	Command    string
	Sudo       bool
	Workers    int
	Background bool
	Debug      bool
}

type Result struct {
	Host string
	Res  string
}

func filterHosts(hosts []string) []string {
	var res []string
	for _, host := range hosts {
		var conn string
		token := strings.Split(host, ":")
		if len(token) == 1 {
			conn = host + ":22"
		} else {
			conn = host
		}
		res = append(res, conn)
	}
	return res
}

func Ssh(hosts []string, conf Conf) []Result {
	k := new(keychain)
	// Add path to id_rsa file
	err := k.loadPEM(conf.Key)

	if err != nil {
		panic("Cannot load key: " + err.Error())
	}

	config := &ssh.ClientConfig{
		User: conf.User,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(k),
		},
	}

	var command string
	if conf.Sudo {
		command = fmt.Sprintf("/usr/bin/sudo bash <<CMD\nexport PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin\n%s\nCMD", conf.Command)
	} else {
		command = fmt.Sprintf("bash <<CMD\nexport PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin\n%s\nCMD", conf.Command)
	}

	if conf.Background {
		command = fmt.Sprintf("/usr/bin/nohup bash -c \\\n\"%s\" `</dev/null` >nohup.out 2>&1 &", command)
	}

	if conf.Debug {
		color.Printf("@{b}%s\n", command)
	}

	var wg sync.WaitGroup
	queue := make(chan Result)
	count := new(int)
	var results []Result
	conns := filterHosts(hosts)
	if conf.Workers == 0 {
		conf.Workers = 24
	}

	for _, conn := range conns {
		wg.Add(1)
		*count++
		if conf.Debug {
			color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
		}
		for *count >= conf.Workers {
			time.Sleep(10 * time.Millisecond)
		}
		go func(h string) {
			defer wg.Done()
			var r Result

			r.Host = h
			client, err := ssh.Dial("tcp", h, config)
			if err != nil {
				color.Printf("@{!r}%s: Failed to connect: %s\n", h, err.Error())
				*count--
				if conf.Debug {
					color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
				}
				return
			}

			session, err := client.NewSession()
			if err != nil {
				color.Printf("@{!r}%s: Failed to create session: %s\n", h, err.Error())
				*count--
				if conf.Debug {
					color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
				}
				return
			}
			defer session.Close()

			/* not working */
			/* session.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin") */

			var b bytes.Buffer
			var e bytes.Buffer
			session.Stdout = &b
			session.Stderr = &e
			if err := session.Run(command); err != nil {
				color.Printf("@{!r}%s: Failed to run: %s\n", h, err.Error())
				color.Printf("@{!r}%s\n", strip(e.String()))
				*count--
				if conf.Debug {
					color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
				}
				return
			}
			if !conf.Background {
				r.Res = strip(b.String())
			} else {
				r.Res = "command success sent and out put in remote server's ~/nohup.out"
			}
			color.Printf("@{!g}%s\n", r.Host)
			fmt.Println(r.Res)

			*count--
			if conf.Debug {
				color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
			}

			runtime.Gosched()
			queue <- r
		}(conn)
	}
	go func() {
		defer wg.Done()
		for r := range queue {
			results = append(results, r)
		}
	}()
	wg.Wait()
	return results
}
