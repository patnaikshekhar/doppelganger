package main

import "doppelganger/cmd"

func main() {
	cmd.Execute()
}

/*

dg ls --> show all the services available locally on the machine
dg fwd service-name
dg fwd --namespaces [NAME]...[NAME2]
dg fwd --all
*/
