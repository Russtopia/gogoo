package main

import "goosurface"
import "strconv"


type GooDelegate struct {
}

func (self *GooDelegate) Draw(s *goosurface.Surface) {
  s.Begin();
    s.Clear(0.5, 0.5, 1.0, 0.5);
  s.End();
}

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

func (self *GooDelegate) Closed(s *goosurface.Surface) {
  print("one of them damn winders closed\n");
}



var s1, s2 *goosurface.Surface;

func main() {
  goosurface.Initialize();
    m := new(GooDelegate);
    s1 = goosurface.CreateSurface(m);
    s2 = goosurface.CreateSurface(m);
    goosurface.Begin();
}
