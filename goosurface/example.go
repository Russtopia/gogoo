package main

import "goosurface"
import "strconv"


type GooDelegate struct {
}

// this method is required by surface delegates!
func (self *GooDelegate) Draw(s *goosurface.Surface) {
  s.Begin();
    s.Clear(0.5, 0.5, 1.0, 0.5);
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
    s.Clear(0.5, 0.5, 1.0, 0.5);
    s.SetColor(0, 0, 0, 1);
    s.MoveTo(40, 30);
    s.SetFontSize(20);
    str := "(";
    str += strconv.Ftoa(x, 'f', -1);
    str += ", ";
    str += strconv.Ftoa(y, 'f', -1);
    str += ")";
    s.ShowText(str);
  s.End();
}

// if we need to do anything when the window closes,
// add this method definition and it will be called.
func (self *GooDelegate) Closed(s *goosurface.Surface) {
  print("GooDelegate: one of them damn winders closed\n");
}



var s1, s2, s3 *goosurface.Surface;

func main() {
  goosurface.Initialize();
    m := new(GooDelegate);
    s1 = goosurface.CreateSurface(m);
    s2 = goosurface.CreateSurface(m);
    s3 = goosurface.CreateSurface(m);
    goosurface.Begin();
}
