package main

import (
	"fmt"
	"flag"
	"github.com/garyburd/redigo/redis"
	"os"
)

var Host string
var Port string
var Auth string
var DB string

func main() {
	flag.StringVar(&Host, "h","127.0.0.1", "redis host")
	flag.StringVar(&Host, "host","127.0.0.1", "redis host")
	flag.StringVar(&Port, "port","6379", "redis port")
	flag.StringVar(&Port, "p","6379", "redis port")
	flag.StringVar(&Auth, "auth","", "redis auth")
	flag.StringVar(&DB, "db","", "redis db")
	flag.Parse()
	fmt.Printf("redis --> host:%s port:%s\n", Host, Port)
	c , err := redis.Dial("tcp", Host+":"+Port);
	defer c.Close()
	if err != nil {
		fmt.Println("connect redis error:", err)
		os.Exit(0);
	}
	if Auth != "" {
		if _, err = c.Do("AUTH", Auth); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
	if DB != "" {
		if _, err = c.Do("SELECT", DB); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
	if _, err = c.Do("PING"); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	for {
		var commond string
		var arg string
		fmt.Printf("%s:%s[%s]>",Host,Port,DB)
		_, err := fmt.Scanf("%s %s", &commond, &arg)
		if err !=nil {
			fmt.Println(err)
			continue
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch string(commond) {
		case "del":
			if arg =="*" {
				fmt.Println("delete * not allow!!!!")
				continue
			}
			fmt.Println("delete keys:", arg)
			res , err := c.Do("keys", arg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			switch res := res.(type) {
			case []interface {}:
				num := 0
				for _,v:= range res {
					switch v := v.(type) {
					case []byte:
						c.Send("DEL", string(v))
						num++
					}
				}
				if err = c.Flush();err != nil {
					fmt.Println(err)
					continue
				} else {
					fmt.Println("delete success: delete num:", num)
				}
			}
		case "keys":
			res, err := c.Do("keys", arg)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			switch  res := res.(type) {
			case []interface {}:
				num := 0
				for _,v:= range res {
					switch v:= v.(type) {
					case []byte:
						num++
						fmt.Println(string(v))
					}
				}
			}
		case "exit":
			os.Exit(0)
		}
	}
}
