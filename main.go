package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

type StateHolder interface {
	State() bool
}

type Active interface {
	Update(tick int)
}
type Actuator interface {
	Action()
}

type Entity interface {
	Id() int
	Description() string
}
type World struct {
	rooms []Room
}
type Room struct {
	Name     string
	Entities []Entity
}
type Item struct {
	id          int
	description string
}

type BinarySignal interface {
	State() bool
}
type Switch struct {
	item  Item
	state bool
}

func (s *Switch) Id() int {
	return s.item.id
}
func (s *Switch) Description() string {
	return s.item.description
}

func (s *Switch) State() bool {
	return s.state
}

func (s *Switch) Action() {
	s.state = !s.state
}

func (i *Item) Description() string {
	return i.description
}
func (i *Item) Id() int {
	return i.id
}

type Led struct {
	item  Item
	state bool
	input *BinarySignal
}

func (s *Led) Id() int {
	return s.item.id
}
func (s *Led) Description() string {
	return s.item.description
}

func (s *Led) State() bool {
	return s.state
}
func (s *Led) Update(tick int) {
	if tick%5 == 0 {
		fmt.Printf("UPDATE LED\n")
		s.state = !s.state
	}
}

func main() {
	w := World{
		rooms: []Room{
			Room{Name: "Machine Room",
				Entities: []Entity{
					&Item{1, "Generator"},
					&Item{2, "Turbine"},
					&Item{3, "Button"},
					&Switch{state: true, item: Item{4, "Switch number one"}},
					&Led{state: false, item: Item{5, "FIRST LED"}},
				},
			},
			Room{Name: "Command Room",
				Entities: []Entity{
					&Item{6, "Terminal"},
					&Switch{state: true, item: Item{7, "Switch number TWO"}},
					&Led{state: false, item: Item{8, "SECOND LED"}},
				},
			},
		},
	}
	tick := 0
	c := make(chan int)
	go getId(c)
	for tick < 300 {
		worldLoop(w, tick)
		tick++
		time.Sleep(2 * time.Second)
		select {
		case cmd := <-c:
			fmt.Printf("#################################################### %v", cmd)
			go getId(c)
			if cmd == 0 {
				os.Exit(2)
			}
		default:
		}
	}

}

func worldLoop(w World, t int) {
	fmt.Printf("==== TICK %d ====\n", t)
	for _, r := range w.rooms {
		fmt.Printf("==== Room %s ====\n", r.Name)
		for _, e := range r.Entities {
			fmt.Printf("id %d is %s and is type %v\n", e.Id(), e.Description(), reflect.TypeOf(e))
		}
		for _, e := range r.Entities {
			switch e.(type) {
			case Actuator:
				fmt.Printf("%s has state: %v\n", e.Description(), e.(*Switch).State())
				e.(Actuator).Action()
				fmt.Printf("%s has state: %v\n", e.Description(), e.(*Switch).State())
			case Active:
				fmt.Printf("%s has state: %v\n", e.Description(), e.(StateHolder).State())
				e.(Active).Update(t)
				fmt.Printf("%s has state: %v\n", e.Description(), e.(StateHolder).State())
			}

		}
	}

}

func getId(c chan int) {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("Invalid number")
		}
		return err
	}
	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	i, _ := strconv.Atoi(result)
	c <- i
}
