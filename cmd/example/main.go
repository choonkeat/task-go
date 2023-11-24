package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/choonkeat/task-go"
)

type emailAddress struct {
	user string
	host string
}

// convert a string to an emailAddress
func toEmailAddress(s string) (emailAddress, error) {
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return emailAddress{}, fmt.Errorf("invalid email address: %s", s)
	}
	return emailAddress{parts[0], parts[1]}, nil
}

func ipaddressOf(host string) (*net.IPAddr, error) {
	return net.ResolveIPAddr("ip", host)
}

func main() {
	for _, str := range os.Args[1:] {
		email, err := toEmailAddress(str)
		if err != nil {
			log.Printf("[idiom]\terr: %s", err)
			continue
		}

		ip, err := ipaddressOf(email.host)
		if err != nil {
			log.Printf("[ipaddressOf] err: %s", err)
			continue
		}
		log.Printf("[idiom]\t%s: %s\n", email.host, ip)
	}

	// do the same thing with task

	for _, str := range os.Args[1:] {
		t1 := task.Wrap(toEmailAddress(str))
		t2 := task.Map(t1, func(email emailAddress) string { return email.host })
		t3 := task.AndThen(t2, func(host string) *task.Task[*net.IPAddr] {
			return task.Wrap(ipaddressOf(host))
		})

		t3.Unwrap(task.Scenarios[*net.IPAddr]{
			Ok: func(ip *net.IPAddr) {
				log.Printf("[unwrap]\t%s: %s\n", str, ip)
			},
			Err: func(err error) {
				log.Printf("[unwrap]\terr: %s", err)
			},
		})
	}
}
