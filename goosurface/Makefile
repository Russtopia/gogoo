include $(GOROOT)/src/Make.$(GOARCH)

TARG=goosurface
CGOFILES=goosurface.go
CGO_CFLAGS=`pkg-config gtk+-2.0 --cflags`
CGO_LDFLAGS=`pkg-config gtk+-2.0 cairo gthread-2.0 --libs`

include $(GOROOT)/src/Make.pkg

%: install %.go
	$(GC) $*.go
	$(LD) -o $@ $*.$O
