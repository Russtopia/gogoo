package main

import "fmt"
import "gocairo"


func main() {
  chin, chout := gocairo.Initialize();

  for {
    evint := <- chin;
    fmt.Printf("Event: %s\n", evint.Description);
  }

  ev := new(gocairo.Event);
  chout <- ev;

}
