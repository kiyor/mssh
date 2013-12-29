# This is multissh golang repo

  This is useful if you work on bunch of server ( in my case like 200 ) and not able to install some manager system like ansiable. ( I have puppet, but puppet is so slow )
  
  Usually if we need quick run command in multiable server, we might do this:
  
  ```bash
  for host in `cat hostsfile`; do
    echo $host
    ssh $host "/usr/bin/sudo ifconfig|grep inet"
  done
  ```
  If you have less then 10 servers, it's fine, but if you have more, more and more in the future...
  
  This is very light and make it super fast.
  
  The only thing you need to do is like that:
  
  in your hosts file
  ```
  hostname1.com
  hostname2.net:20
  ```
  It's able to detect if you have different port running sshd, default is 22 if you don't add :port
  
  ```go
  package main

  import (
	  "flag"
	  "github.com/kiyor/mssh"
  )

  var (
	  file *string = flag.String("hostsfile", "", "use hostsfile")
	  cmd *string = flag.String("cmd", "hostname", "execute command")
  )

  func main() {
	  flag.Parse()
	  var hosts []string

	  hosts = mssh.GetHostsByFile(*file)

	  var conf mssh.Conf
	  conf.User = "kiyor"
	  conf.Key = "/home/kiyor/.ssh/id_rsa"
	  conf.Command = *cmd
	  mssh.Ssh(hosts, conf)
	  }
  ```
  ```bash
  ./yourbinary -hostsfile host -cmd "ifconfig|grep inet"
  ```
  That's all, it will send 24 request at same time.
  It's lite, it's quick, you don't need full path, and support nohup.out
  
  some option
  
  ```
  conf.Sudo = false //option default=false                                           
  conf.Debug = false //option default=false
  conf.Workers = 24 //option default=24
  conf.Background = false //option default=false
  ```

  
