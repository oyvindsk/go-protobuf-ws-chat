package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/oyvindsk/go-protobuf-ws-chat/ab"
)

var fname = "test.pb"

func main() {
	book := &ab.AddressBook{}

	book.People = append(book.People, &ab.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*ab.Person_PhoneNumber{
			{Number: "555-4321", Type: ab.Person_HOME},
		},
	})
	// Write the new address book back to disk.
	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}

	// Read the existing address book.
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	book2 := &ab.AddressBook{}
	if err := proto.Unmarshal(in, book2); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	fmt.Println("book2:", book2)
	for _, p := range book2.People {
		fmt.Println("name:", p.Name)
	}
}
