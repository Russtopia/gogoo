#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <cairo/cairo.h>
#include <cairo/cairo-xlib.h>
#include <gtk/gtk.h>
#include <glib.h>


static int gcpipefd[2];
static int gcsurface_id = 0;
static float gcmousepos[2];


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
  int  _id;
  void *surface;
  void *context;
  void *window;
} gcsurface;

typedef struct {
  int w;
  int h;
} gcsize;

typedef struct {
  int x;
  int y;
} gcpoint;

typedef struct {
  gcsize size;
  gcpoint origin;
} gcrect;


//
// Callbacks
//

// called when gtk+ window closes
gboolean gcclose(GtkWidget *widget, gpointer udata) {
  gcevent_send('x', *((int*)udata));
  free((int*)udata);
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


void gcbuttonpress(GtkWidget *widget, GdkEventButton *ev, gpointer udata) {
  if(ev->button == 1) gcevent_send('1', *((int*)udata));
  if(ev->button == 2) gcevent_send('2', *((int*)udata));
  if(ev->button == 3) gcevent_send('3', *((int*)udata));
}



void gcbuttonrelease(GtkWidget *widget, GdkEventButton *ev, gpointer udata) {
  if(ev->button == 1) gcevent_send('!', *((int*)udata));
  if(ev->button == 2) gcevent_send('@', *((int*)udata));
  if(ev->button == 3) gcevent_send('#', *((int*)udata));
  if(ev->button == 4) gcevent_send('3', *((int*)udata));
}



// called when a new area of the gtk+ window becomes exposed
gboolean gcexpose(GtkWidget *widget, GdkEventExpose *ev, gpointer udata) {
  //gint w, h;
  //gtk_window_get_size(GTK_WINDOW(widget), &w, &h);
  gcevent_send('e', *((int*)udata));

  return FALSE;
}


void gcbegincontext(gcsurface *s) {
  gdk_threads_enter();
  cairo_t *context = NULL;
  context = gdk_cairo_create(((GtkWidget*)s->window)->window);
  if(!context) { printf("cairo fail\n"); return; }
  cairo_set_operator(context, CAIRO_OPERATOR_SOURCE);
  s->context = (void*)context;
  gdk_threads_leave();
}


void gcendcontext(gcsurface *s) {
  gdk_threads_enter();
  cairo_destroy((cairo_t*)s->context);
  gdk_threads_leave();
}


void gcsetcolor(gcsurface *s, float r, float g, float b, float a) {
  gdk_threads_enter();
  cairo_set_source_rgba((cairo_t*)s->context, r, g, b, a);
  gdk_threads_leave();
}


void gcclear(gcsurface *s, float r, float g, float b, float a) {
  gdk_threads_enter();
  cairo_set_source_rgba((cairo_t*)s->context, r, g, b, a);
  cairo_paint((cairo_t*)s->context);
  gdk_threads_leave();
}


void gcmoveto(gcsurface *s, float x, float y) {
  gdk_threads_enter();
  cairo_move_to((cairo_t*)s->context, x, y);
  gdk_threads_leave();
}


void gclineto(gcsurface *s, float x, float y) {
  gdk_threads_enter();
  cairo_line_to((cairo_t*)s->context, x, y);
  gdk_threads_leave();
}


void gcstroke(gcsurface *s) {
  gdk_threads_enter();
  cairo_stroke((cairo_t*)s->context);
  gdk_threads_leave();
}


void gcsetfontsize(gcsurface *s, int size) {
  gdk_threads_enter();
  cairo_set_font_size((cairo_t*)s->context, size);
  gdk_threads_leave();
}


void gcshowtext(gcsurface *s, char *text) {
  gdk_threads_enter();
  cairo_select_font_face((cairo_t*)s->context, "Sans", CAIRO_FONT_SLANT_NORMAL, CAIRO_FONT_WEIGHT_NORMAL);
  cairo_show_text((cairo_t*)s->context, text);
  gdk_threads_leave();
}


void gcrectangle(gcsurface *s, float x, float y, float w, float h) {
  gdk_threads_enter();
  cairo_rectangle((cairo_t*)s->context, x, y, w, h);
  gdk_threads_leave();
}


void gcfill(gcsurface *s) {
  gdk_threads_enter();
  cairo_fill((cairo_t*)s->context);
  gdk_threads_leave();
}


float gcmousex(void) {
  return gcmousepos[0];
}



float gcmousey(void) {
  return gcmousepos[1];
}


gboolean gcmotion(GtkWidget *widget, GdkEventMotion *ev, gpointer udata) {
  gcevent_send('m', *((int*)udata));
  gcmousepos[0] = ev->x;
  gcmousepos[1] = ev->y;
  return FALSE;
}


void gcinit(void) {
  g_thread_init(NULL);
  gdk_threads_init();
  gtk_init(NULL, NULL);
  pipe(gcpipefd);
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


void gcmsgbox(gcsurface *s, char *msg, char *title) {
  GtkWidget *dialog = gtk_dialog_new_with_buttons(title, s->window, 0, NULL);
  GtkWidget *label = gtk_label_new(msg);
  gtk_container_add(GTK_CONTAINER(GTK_DIALOG(dialog)->vbox), label);
  gtk_widget_show_all(dialog);
  gtk_dialog_run(GTK_DIALOG(dialog));
}



gcsurface* gcsurfacecreate(void) {
  gdk_threads_enter();
  gcsurface *s = malloc(sizeof(gcsurface));
  s->window = (void*)gtk_window_new(GTK_WINDOW_TOPLEVEL);
  gtk_window_set_title(GTK_WINDOW(s->window), "Go Canvas");
  gcreset(s->window, NULL, NULL);

  gtk_widget_set_app_paintable((GtkWidget*)s->window, TRUE);

  int *sid = malloc(sizeof(int));
  (*sid) = gcsurface_id++;
  s->_id = *sid;
  g_signal_connect(G_OBJECT(s->window), "destroy", G_CALLBACK(gcclose), sid);
  g_signal_connect(G_OBJECT(s->window), "expose-event", G_CALLBACK(gcexpose), sid);
  g_signal_connect(G_OBJECT(s->window), "configure-event", G_CALLBACK(gcexpose), sid);
  g_signal_connect(G_OBJECT(s->window), "screen-changed", G_CALLBACK(gcreset), sid);
  g_signal_connect(G_OBJECT(s->window), "button-press-event", G_CALLBACK(gcbuttonpress), sid);
  g_signal_connect(G_OBJECT(s->window), "button-release-event", G_CALLBACK(gcbuttonrelease), sid);
  gtk_widget_set_events((GtkWidget*)s->window, gtk_widget_get_events((GtkWidget*)s->window) | GDK_BUTTON_PRESS_MASK | GDK_BUTTON_RELEASE_MASK);

  gdk_threads_leave();
  return s;
}

void gcsurfaceshow(gcsurface *s) {
  gtk_widget_show((GtkWidget*)s->window);
}

void gcsurfaceenablemm(gcsurface *s) {
  gdk_threads_enter();
  gtk_widget_set_events((GtkWidget*)s->window, gtk_widget_get_events((GtkWidget*)s->window) | GDK_POINTER_MOTION_MASK | GDK_POINTER_MOTION_HINT_MASK);
  g_signal_connect(G_OBJECT(s->window), "motion-notify-event", G_CALLBACK(gcmotion), &(s->_id));
  gdk_threads_leave();
}


void gcsurfacesetsize(gcsurface *s, int width, int height) {
  gdk_threads_enter();
  gtk_window_set_default_size((GtkWindow*)s->window, width, height);
  gdk_threads_leave();
}

gcsize* gcsurfacegetsize(gcsurface *s) {
  gcsize *size = malloc(sizeof(gcsize));
  gdk_threads_enter();
  gtk_window_get_size((GtkWindow*)s->window, &(size->w), &(size->h));
  gdk_threads_leave();
  return size;
}


