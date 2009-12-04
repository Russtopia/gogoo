package gocairo

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <cairo/cairo.h>
#include <cairo/cairo-xlib.h>
#include <gtk/gtk.h>
#include <glib.h>


static int gcpipefd[2];
static int gcsurface_id = 0;


// structure for passing an event to Go
typedef struct {
  char _type;
  int  _id;
} gcevent;


void gcevent_send(char c, int id) {
  gcevent *gcev = malloc(sizeof(gcevent));
  gcev->_type = c;
  gcev->_id = id;
  write(gcpipefd[1], gcev, sizeof(gcevent));
  free(gcev);
}


// structure for containing all neccesary
// pointers to a surface's real data
typedef struct {
  void *surface;
  void *context;
  void *window;
} gcsurface;


//
// Callbacks
//

// called when gtk+ window closes
gboolean gcclose(GtkWidget *widget, gpointer udata) {
  gcevent_send('x', *((int*)udata));
  return FALSE;
}


// called when gtk+ window needs to reset it's colormap
// usually only called right when the window is created.
// also called if a window changes X11 screens
void gcreset(GtkWidget *widget, GdkScreen *oldscreen, gpointer userData) {
  GdkScreen* screen = gtk_widget_get_screen (widget);
  GdkColormap* colormap = gdk_screen_get_rgba_colormap (screen);
  if(!colormap) { printf("rgbafail\n"); colormap = gdk_screen_get_rgb_colormap (screen); }
  gtk_widget_set_colormap(widget, colormap);
}


// called when a new area of the gtk+ window becomes exposed
gboolean gcexpose(GtkWidget *widget, GdkEventExpose *ev, gpointer udata) {
  //gint w, h;
  //gtk_window_get_size(GTK_WINDOW(widget), &w, &h);

  gcevent_send('e', *((int*)udata));

  return FALSE;

  //cairo_t *context = NULL;
  //context = gdk_cairo_create(widget->window);
  //if(!context) { printf("cairofail\n"); return FALSE; }
  //cairo_set_source_rgba(context, 1.0f, 1.0f, 1.0f, 0.7f);
  //cairo_set_operator(context, CAIRO_OPERATOR_SOURCE);
  //cairo_paint(context);
  //cairo_destroy(context);
}


gboolean gcmotion(GtkWidget *widget, GdkEventMotion *ev, gpointer udata) {
  gcevent_send('m', *((int*)udata));
  return FALSE;
}


int gcinit(void) {
  g_thread_init(NULL);
  gdk_threads_init();
  gtk_init(NULL, NULL);
  pipe(gcpipefd);
  return gcpipefd[1];
}


void gciterate(void) {
  gdk_threads_enter();
  while(gtk_events_pending()) gtk_main_iteration();
  gdk_threads_leave();
}


void gccheckev(void) {
  gboolean i = 0;
  do {
    i = g_main_context_pending(NULL);
    usleep(9000);
  } while(!i);
}


gcevent* gcget(void) {
  gcevent *gcev = malloc(sizeof(gcevent));
  read(gcpipefd[0], gcev, sizeof(gcevent));
  return gcev;
}


void gcfree(gcevent *gcev) {
  free((gcevent*)gcev);
}



void* gcnew_surface(void) {
  gcsurface *s = malloc(sizeof(gcsurface));
  s->window = (void*)gtk_window_new(GTK_WINDOW_TOPLEVEL);
  gtk_window_set_title(GTK_WINDOW(s->window), "Go Canvas");
  gcreset(s->window, NULL, NULL);

  gtk_widget_set_app_paintable((GtkWidget*)s->window, TRUE);
  gtk_widget_set_events((GtkWidget*)s->window, gtk_widget_get_events((GtkWidget*)s->window) | GDK_POINTER_MOTION_MASK | GDK_POINTER_MOTION_HINT_MASK);

  int *sid = malloc(sizeof(int));
  (*sid) = gcsurface_id++;
  g_signal_connect(G_OBJECT(s->window), "destroy", G_CALLBACK(gcclose), sid);
  g_signal_connect(G_OBJECT(s->window), "expose-event", G_CALLBACK(gcexpose), sid);
  g_signal_connect(G_OBJECT(s->window), "motion-notify-event", G_CALLBACK(gcmotion), sid);
  g_signal_connect(G_OBJECT(s->window), "screen-changed", G_CALLBACK(gcreset), sid);

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
  gcpipefd = int(C.gcinit());
  gcoutchan = chout;
  gcinchan = chin;

  go eventd(chev);  // gui event daemon (for guid)
  go inputd(chout); // input event daemon (for us)

  s := CreateSurface();
  s.doNada();

  for {
    select {                  // we either
      case <- chev:           // wait for a gtk event to process
        C.gciterate();
      case ev := <- chin:     // or wait for commands from the backend
        print("command into gui: ", ev.Description, "\n");
    }
  }
}


func inputd(chin chan *Event) {
  for {
    gcev := C.gcget();
    
    // todo: process message
    ev := new(Event);
    ev.Description = string(gcev._type);
    ev.SurfaceID = int(gcev._id);
    C.gcfree(gcev);

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
  p := cpointer(C.gcnew_surface());

  s.setPointer(p);

  return s;
}
//


