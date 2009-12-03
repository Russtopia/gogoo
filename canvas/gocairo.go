package gocairo

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <cairo/cairo.h>
#include <cairo/cairo-xlib.h>
#include <gtk/gtk.h>
#include <glib.h>


static int gocairo_pipefd[2];
static int gocairo_surface_id = 0;


typedef struct {
  void *surface;
  void *context;
  void *window;
} gocairo_surface;


void gocairo_resetsurface(GtkWidget *widget, GdkScreen *oldscreen, gpointer userData) {
  GdkScreen* screen = gtk_widget_get_screen (widget);
  GdkColormap* colormap = gdk_screen_get_rgba_colormap (screen);
  if(!colormap) { printf("rgbafail\n"); colormap = gdk_screen_get_rgb_colormap (screen); }
  gtk_widget_set_colormap(widget, colormap);
}


gboolean gocairo_expose(GtkWidget *widget, GdkEventExpose *ev, gpointer udata) {
  gint w, h;
  gtk_window_get_size(GTK_WINDOW(widget), &w, &h);

  //int *sid = (int*)udata;
  char data = 'e';
  write(gocairo_pipefd[1], &data, 1);

  return FALSE;

  //cairo_t *context = NULL;
  //context = gdk_cairo_create(widget->window);
  //if(!context) { printf("cairofail\n"); return FALSE; }
  //cairo_set_source_rgba(context, 1.0f, 1.0f, 1.0f, 0.7f);
  //cairo_set_operator(context, CAIRO_OPERATOR_SOURCE);
  //cairo_paint(context);
  //cairo_destroy(context);
}


gboolean gocairo_motion(GtkWidget *widget, GdkEventMotion *ev, gpointer udata) {
  char data = 'm';
  write(gocairo_pipefd[1], &data, 1);
  return FALSE;
}


int gocairo_init(void) {
  g_thread_init(NULL);
  gdk_threads_init();
  gtk_init(NULL, NULL);
  pipe(gocairo_pipefd);
  return gocairo_pipefd[1];
}


void gocairo_iterate(void) {
  gdk_threads_enter();
  while(gtk_events_pending()) gtk_main_iteration();
  gdk_threads_leave();
}


void gocairo_checkev(void) {
  gboolean i = 0;
  do {
    i = g_main_context_pending(NULL);
    usleep(9000);
  } while(!i);
}


char gocairo_get(void) {
  char c;

  read(gocairo_pipefd[0], &c, 1);

  return c;
}




void* gocairo_new_surface(void) {
  gocairo_surface *s = malloc(sizeof(gocairo_surface));
  s->window = (void*)gtk_window_new(GTK_WINDOW_TOPLEVEL);
  gtk_window_set_title(GTK_WINDOW(s->window), "Go Canvas");
  gocairo_resetsurface(s->window, NULL, NULL);

  gtk_widget_set_app_paintable((GtkWidget*)s->window, TRUE);
  gtk_widget_set_events((GtkWidget*)s->window, gtk_widget_get_events((GtkWidget*)s->window) | GDK_POINTER_MOTION_MASK | GDK_POINTER_MOTION_HINT_MASK);

  int *sid = malloc(sizeof(int));
  (*sid) = gocairo_surface_id++;
  g_signal_connect(G_OBJECT(s->window), "expose-event", G_CALLBACK(gocairo_expose), sid);
  g_signal_connect(G_OBJECT(s->window), "motion-notify-event", G_CALLBACK(gocairo_motion), sid);
  g_signal_connect(G_OBJECT(s->window), "screen-changed", G_CALLBACK(gocairo_resetsurface), sid);

  gtk_widget_show_all(GTK_WIDGET(s->window));

  return (void*)s;
}


*/
import "C"
import "unsafe"
//import "fmt"


// Type definitions
type Event struct {
  Description string;
  SurfaceID int;
}


// Global variables
var gcpipefd int;
var gcoutchan chan *Event;
var gcinchan chan *Event;


/*
 *  Functions for communication between goroutines, etc
 *
 *
 *
 */
func Initialize() (chan *Event, chan *Event) {

  // create our channels
  chin  := make(chan *Event);
  chout := make(chan *Event);
  chev  := make(chan bool);

  go guid(chin, chout, chev); // gui daemon

  return chin, chout;
}


func guid(chout chan *Event, chin chan *Event, chev chan bool) {
  gcpipefd = int(C.gocairo_init());
  gcoutchan = chout;
  gcinchan = chin;

  go eventd(chev);  // gui event daemon (for guid)
  go inputd(chout); // input event daemon (for us)

  s := CreateSurface();
  s.doNada();

  for {
    select {                  // we either
      case <- chev:           // wait for a gtk event to process
        C.gocairo_iterate();
      case ev := <- chin:     // or wait for commands from the backend
        print("command into gui: ", ev.Description, "\n");
    }
  }
}


func inputd(chin chan *Event) {
  for {
    evtype := string(C.gocairo_get());
    
    // todo: process message
    ev := new(Event);
    ev.Description = evtype;

    chin <- ev; // dispatch event to main thread
  }
}


func eventd(ch chan bool) {
  for {
    C.gocairo_checkev();
    ch <- true;
  }
}



/*
 * Surface structure
 *
 *
 *
 */
type cpointer unsafe.Pointer
type Surface struct {
  pointer cpointer;
}

func (self *Surface) setPointer(p cpointer) {
  self.pointer = p;
}

func (self *Surface) doNada() {
}

func CreateSurface() *Surface {
  s := new(Surface);
  p := cpointer(C.gocairo_new_surface());

  s.setPointer(p);

  return s;
}
//


