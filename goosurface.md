# About #

this package implements a drawable gui surface using cgo to create a wrapper for cairo and gtk2. A user defined interface is passed to the package, where the package will execute interface methods in response to gui events. What's particularly nice about the callbacks is that they're presence means they will get called. No other setup is neccesary (a particularly nice feature of Go's interfaces).

Here's a screenshot of the example 'main' file in the source repository. It creates two surfaces, each with the same delegate (meaning, they will behave identically in this case). Mouse motion is tracked by an interface function and drawn to the surface using simple commands.

![http://lh5.ggpht.com/_asb2EOPVV_8/Sx7UqFRadRI/AAAAAAAAADc/eYcH_eA-JiA/s800/Screenshot.png](http://lh5.ggpht.com/_asb2EOPVV_8/Sx7UqFRadRI/AAAAAAAAADc/eYcH_eA-JiA/s800/Screenshot.png)