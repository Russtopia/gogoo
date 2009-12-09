package main

import "goosurface"
import "fmt"


type GooDelegate struct {}

// this method is required by surface delegates!
func (self *GooDelegate) Draw(s *goosurface.Surface) {
  s.Begin();
    s.Clear(0.0, 0.0, 0.0, 0.9);
  s.End();
}

// if we need to do any kind of setup when a surface associated
// with this delegate is created, add the Initialize method
// definition.
func (self *GooDelegate) Initialize(s *goosurface.Surface) {
  print("GooDelegate: a surface was just created!\n");
}

// if we want to know about mouse motion on our surface(s) 
// associated with a delegate, we add this method.
// Simply by putting this here, the goosurface package
// knows to let gtk2 know that we want mouse motion events.
func (self *GooDelegate) MouseMoved(s *goosurface.Surface, x float, y float) {
  s.Begin();
    s.Clear(0, 0, 0, 0.9);
    s.SetColor(1, 1, 1, 0.8);
    s.MoveTo(40, 30);
    s.SetFontSize(16);
    str := fmt.Sprintf("(%6.2f, %6.2f)", x, y);
    s.ShowText(str);
  s.End();
}

// if we need to do anything when the window closes,
// add this method definition and it will be called.
func (self *GooDelegate) Closed(s *goosurface.Surface) {
  print("GooDelegate: one of them damn winders closed\n");
}

func (self *GooDelegate) ButtonDown(s *goosurface.Surface, b int) {
  fmt.Printf("GooDelegate: button %d down\n", b);
}

func (self *GooDelegate) ButtonUp(s *goosurface.Surface, b int) {
  fmt.Printf("GooDelegate: button %d up\n", b);
  goosurface.MessageBox(s, "asdf", "asdf");
}



// main function
func main() {
  goosurface.Initialize();
  d := new(GooDelegate);

  s1 := goosurface.CreateSurface(d);
  s1.SetSize(420, 240);
  s1.Show();

  s2 := goosurface.CreateSurface(d);
  s2.SetSize(240, 420);
  s2.Show();

  goosurface.Begin();
}




