package main

import (
	"encoding/json"
	"fmt"
	"github.com/goprotobufexample/person"
	"google.golang.org/protobuf/proto"
)

func main() {
	person1 := person.Person{}
	person1.PersonName = "feng"
	person1.PersonAge = 18
	person1.PersonAddress = "China"
	person1.PersonId = 1

	bytes, _ := json.Marshal(person1)
	fmt.Println("bytes =", bytes, " \n", "person1 : ", string(bytes))

	out, err := proto.Marshal(&person1)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("out = ", out)

	person2 := person.Person{}
	err = proto.Unmarshal(out, &person2)

	if err != nil {
		fmt.Println(err)
	}
	bytes, _ = json.Marshal(person2)
	fmt.Println("person2", string(bytes))
}
