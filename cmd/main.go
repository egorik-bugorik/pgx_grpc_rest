package main

import "flag"

var (
	httpAddress = flag.String("http", "localhost:3000", "http addres of server is listen to... ")
	grpcAddress = flag.String("grps", "localhost:8082", "grps addres of server is listen to... ")
)

func main() {
	flag.Parse()

}
