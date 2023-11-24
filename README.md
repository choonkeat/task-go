# Task library for Go

I know it's _the_ Go idiom, but let's say you don't want to deal with `data, err` there and then.
Let's say, you want to be able to compose failble functions into a unit of work.

# Usage

With the context of a type and some functions that could fail with error

```go
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
```

We would like to do a series of transformations

```go
task.Wrap(toEmailAddress(str)).
    Map(func(email emailAddress) string { return email.host }).
    AndThen(func(host string) *task.Task[*net.IPAddr] { return task.Wrap(ipaddressOf(host)) }).
    Unwrap(task.Scenarios[*net.IPAddr]{
        Ok: func(ip *net.IPAddr) {
            log.Printf("[unwrap]\t%s: %s\n", str, ip)
        },
        Err: func(err error) {
            log.Printf("[unwrap]\terr: %s", err)
        },
    })
```

> ```
> main.go: method must have no type parameters
> ```

So we cannot have nice things.

