package goosurface

/*
#include "gswrap.c"
*/
import "C"

var garbage interface{};

/*
 *  type definitions
 */
  type Event struct {
    _Description string;
    _Surface *Surface;
  }


  type Message struct {
    _Description string;
    _Surface *Surface;
  }


  type Size struct {
    Width float;
    Height float;
  }


  type SurfaceDelegate interface {
    Draw(*Surface);
  }
  type surfaceDelegateInit interface { Initialize(*Surface); }
  type surfaceDelegateClosed interface { Closed(*Surface); }
  type surfaceDelegateMouseMotion interface { MouseMoved(*Surface, float, float); }
  type surfaceDelegateButtonDown interface { ButtonDown(*Surface, int); }
  type surfaceDelegateButtonUp interface { ButtonUp(*Surface, int); }


  type Surface struct {
    _ID int;
    _Delegate SurfaceDelegate;
    _Pointer *C.gcsurface;
    _Size Size;
  }



/*
 *  global constants
 */
  var chin, chout chan interface {};
  var chev chan bool;
  var surfacemap map[int]*Surface;



/*
 *  Functions for communication between goroutines, etc
 */
func Initialize() {

  // create our channels
  chin  = make(chan interface {});
  chout = make(chan interface {});
  chev  = make(chan bool);
  surfacemap = make(map[int]*Surface);

  go guid(chin, chout, chev); // gui daemon
  <-chev; // wait on signal from guid to let us know gtk is initialized.
}


func guid(chout chan interface {}, chin chan interface {}, chev chan bool) {
  C.gcinit();
  chev <- true;

  go eventd(chev);  // gui event daemon (for guid)
  go inputd(chout); // input event daemon (for us)

  for {
    select {                       // we either
      case <- chev: C.gciterate(); // wait for a gtk event to process

      case evi := <-chin:          // or wait for commands from the backend
        ev, ok := evi.(*Event);
        if ok {
          print("command into gui: ", ((*Event)(ev))._Description, "\n");
        } else {
          print("something went awry with the gui msg");
        }

    }
  }
}


func inputd(chin chan interface {}) {
  for {
    gcev := C.gcget(); // get c event

    ev := new(Event);
    ev._Description = string(gcev._type);
    ev._Surface = surfacemap[int(gcev._id)];
    C.gcfree(gcev); // free c event

    chin <- ev; // dispatch event to main thread
  }
}


func eventd(ch chan bool) {
  for {
    C.gccheckev();
    ch <- true;
  }
}



/*
 * Surface structure
 */
func CreateSurface(d SurfaceDelegate) *Surface {
  s := new(Surface);
  s._Pointer = C.gcsurfacecreate();
  s._Delegate = d;
  s._ID = int(s._Pointer._id);
  surfacemap[(s._ID)] = s;

  // if the delegate has an Initialize method, call it
  del, okinit := d.(surfaceDelegateInit);
  if okinit {
    del.Initialize(s);
  }

  // if the delegate needs mouse events, enable them in gtk
  delmm, okmm := d.(surfaceDelegateMouseMotion);
  if okmm {
    garbage = delmm.(interface {});
    C.gcsurfaceenablemm(s._Pointer);
  }

  return s;
}


func (self *Surface) Begin() {
  C.gcbegincontext(self._Pointer);
}


func (self *Surface) End() {
  C.gcbegincontext(self._Pointer);
}


func (self *Surface) SetSize(w int, h int) {
  C.gcsurfacesetsize(self._Pointer, C.int(w), C.int(h));
}


func (self *Surface) Show() {
  C.gcsurfaceshow(self._Pointer);
}


func (self *Surface) SetColor(r float, g float, b float, a float) {
  C.gcsetcolor(self._Pointer, C.float(r), C.float(g), C.float(b), C.float(a));
}


func (self *Surface) Clear(r float, g float, b float, a float) {
  C.gcclear(self._Pointer, C.float(r), C.float(g), C.float(b), C.float(a));
}


func (self *Surface) MoveTo(x float, y float) {
  C.gcmoveto(self._Pointer, C.float(x), C.float(y));
}


func (self *Surface) LineTo(x float, y float) {
  C.gclineto(self._Pointer, C.float(x), C.float(y));
}


func (self *Surface) Rectangle(x float, y float, w float, h float) {
  C.gcrectangle(self._Pointer, C.float(x), C.float(y), C.float(w), C.float(h));
}


func (self *Surface) Stroke() {
  C.gcstroke(self._Pointer);
}


func (self *Surface) Fill() {
  C.gcfill(self._Pointer);
}


func (self *Surface) Paint() {
  C.gcfill(self._Pointer);
}


func (self *Surface) SetFontSize(size int) {
  C.gcsetfontsize(self._Pointer, C.int(size));
}


func (self *Surface) ShowText(text string) {
  C.gcshowtext(self._Pointer, C.CString(text));
}


func (self *Surface) UpdateSize() {
  size := C.gcsurfacegetsize(self._Pointer);
  self._Size.Width = float(size.w);
  self._Size.Height = float(size.h);
}



/*
 *  Main loop
 */
func Begin() {
  for {
    evi := <-chin;
    ev, ok := evi.(*Event);
    var s *Surface = ev._Surface;

    if ok {
      if ev._Description == "e" {
        s.UpdateSize();
        s._Delegate.Draw(s);
      } else if ev._Description == "m" {
        delm, okdm := s._Delegate.(surfaceDelegateMouseMotion);
        if okdm {
          mx := float(C.gcmousex());
          my := float(C.gcmousey());
          delm.MouseMoved(s, mx, my);
        }
      } else if ev._Description == "x" {
        surfacemap[s._ID] = nil;
        delx, okx := s._Delegate.(surfaceDelegateClosed);
        if okx { delx.Closed(s); }
      } else if ev._Description == "1" {
        del1, ok1 := s._Delegate.(surfaceDelegateButtonDown);
        if ok1 { del1.ButtonDown(s, 1); }
      } else if ev._Description == "2" {
        del1, ok1 := s._Delegate.(surfaceDelegateButtonDown);
        if ok1 { del1.ButtonDown(s, 2); }
      } else if ev._Description == "3" {
        del1, ok1 := s._Delegate.(surfaceDelegateButtonDown);
        if ok1 { del1.ButtonDown(s, 3); }
      } else if ev._Description == "!" {
        del11, ok11 := s._Delegate.(surfaceDelegateButtonUp);
        if ok11 { del11.ButtonUp(s, 1); }
      } else if ev._Description == "@" {
        del11, ok11 := s._Delegate.(surfaceDelegateButtonUp);
        if ok11 { del11.ButtonUp(s, 2); }
      } else if ev._Description == "#" {
        del11, ok11 := s._Delegate.(surfaceDelegateButtonUp);
        if ok11 { del11.ButtonUp(s, 3); }
      }
    }
  }
}


/*
 *  Misc.
 */
func MessageBox(s *Surface, msg string, title string) {
}



