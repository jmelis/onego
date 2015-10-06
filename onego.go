package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OpenNebula/goca"
	"github.com/codegangsta/cli"
)

func exitError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	if _, err := goca.SystemVersion(); err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "onego"
	app.Usage = "OpenNebula Utility Belt for CLI ninjas"
	app.Author = "Jaime Melis <jmelis@opennebula.org>"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "ip",
			Usage:  "Get ip of a VM",
			Action: cmdIp,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "id",
					Value: -1,
					Usage: "Id of the VM. If specified, --name will be ignored.",
				},
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "Name of the VM",
				},
				cli.BoolFlag{
					Name:  "all",
					Usage: "Get all the IPs of the VM.",
				},
			},
		},
	}

	app.Run(os.Args)
}

func cmdIp(c *cli.Context) {
	var (
		vm  *goca.VM
		err error
	)

	id := c.Int("id")
	name := c.String("name")

	if id == -1 && name == "" {
		exitError("Please provide --id or --name.")
	}

	if id != -1 {
		vm = goca.NewVM(uint(id))
	} else if name != "" {
		vm, err = goca.NewVMFromName(name)
	}

	if err != nil {
		log.Fatal(err)
	}

	if err := vm.Info(); err != nil {
		log.Fatal(err)
	}

	if c.Bool("all") {
		iter := vm.XPathIter("/VM/TEMPLATE/NIC")
		for iter.Next() {
			node := iter.Node()
			ip, ok := node.XPath("IP")
			if ok == false {
				exitError("Unable to find IP.")
			}

			fmt.Println(ip)
		}
	} else {
		ip, ok := vm.XPath("/VM/TEMPLATE/NIC/IP")
		if ok == false {
			exitError("Unable to find IP.")
		}

		fmt.Println(ip)
	}
}
