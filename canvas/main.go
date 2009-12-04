package main

import "fmt"
import "gocairo"


func main() {
  chin, chout := gocairo.Initialize();

  for {
    evint := <- chin;
    fmt.Println("event ", evint.Description, " for surface ", evint.SurfaceID);
  }

  ev := new(gocairo.Event);
  chout <- ev;

}
