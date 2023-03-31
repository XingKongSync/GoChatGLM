package main

import (
	"flag"
	"fmt"
	"mychatglm/service"
)

var port = flag.Int("port", 3001, "Port of http server")
var wwwPath = flag.String("www", "./www", "Static files path")

func main() {
	flag.Parse()

	fmt.Printf("Server is running on port: %d\n", *port)
	fmt.Printf("Static files path: %s\n", *wwwPath)

	service.NewServer(wwwPath).Run(fmt.Sprintf(":%d", *port))
}
