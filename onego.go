package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/OpenNebula/goca"
	"github.com/codegangsta/cli"
)

func exitError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func checkFlagsIncompatible(c *cli.Context, flags ...string) {
	count := 0
	fflags := make([]string, len(flags))
	for i, f := range flags {
		fflags[i] = "--" + f

		if c.IsSet(f) {
			count++
		}
	}

	if count > 1 {
		msg := "Specify only one of the following flags: " + strings.Join(fflags, ", ")
		exitError(msg)
	}
}

func checkFlagsMust(c *cli.Context, flags ...string) {
	count := 0
	fflags := make([]string, len(flags))
	for i, f := range flags {
		fflags[i] = "--" + f

		if c.IsSet(f) {
			count++
		}
	}

	if count < 1 {
		msg := "Specify one of the following flags: " + strings.Join(fflags, ", ")
		exitError(msg)
	}

}

func main() {
	app := cli.NewApp()
	app.Name = "onego"
	app.Usage = "OpenNebula Utility Belt for CLI ninjas"
	app.Author = "Jaime Melis <jmelis@opennebula.org>"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "ip",
			Usage:  "Get IP of a VM",
			Action: cmdIp,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "id",
					Value: -1,
					Usage: "Id of the VM.",
				},
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "Name of the VM.",
				},
				cli.IntFlag{
					Name:  "nic_id",
					Value: -1,
					Usage: "Get the IP of this NIC.",
				},
				cli.IntFlag{
					Name:  "network_id",
					Value: -1,
					Usage: "Get the IP of this Network ID.",
				},
				cli.StringFlag{
					Name:  "network",
					Value: "",
					Usage: "Get the IP of this Network.",
				},
				cli.BoolFlag{
					Name:  "all",
					Usage: "Get all the IPs instead of just the first one.",
				},
			},
		},
	}

	app.Run(os.Args)
}

func cmdIp(c *cli.Context) {
	checkFlagsMust(c, "id", "name")
	checkFlagsIncompatible(c, "id", "name")
	checkFlagsIncompatible(c, "nic_id", "network", "network_id")

	var (
		vm  *goca.VM
		err error
	)

	if c.IsSet("id") {
		vm = goca.NewVM(uint(c.Int("id")))
	} else {
		if vm, err = goca.NewVMFromName(c.String("name")); err != nil {
			log.Fatal(err)
		}
	}

	if err = vm.Info(); err != nil {
		log.Fatal(err)
	}

	xpath := "/VM/TEMPLATE/NIC"

	if c.IsSet("nic_id") {
		xpath += fmt.Sprintf("[NIC_ID='%d']", c.Int("nic_id"))
	} else if c.IsSet("network_id") {
		xpath += fmt.Sprintf("[NETWORK_ID='%d']", c.Int("network_id"))
	} else if c.IsSet("network") {
		xpath += fmt.Sprintf("[NETWORK='%s']", c.String("network"))
	}

	if c.Bool("all") {
		iter := vm.XPathIter(xpath)
		for iter.Next() {
			node := iter.Node()
			ip, ok := node.XPath("IP")
			if ok == false {
				exitError("Unable to find IP.")
			}

			fmt.Println(ip)
		}
	} else {
		xpath += "/IP"
		ip, ok := vm.XPath(xpath)
		if ok == false {
			exitError("Unable to find IP.")
		}

		fmt.Println(ip)
	}
}
