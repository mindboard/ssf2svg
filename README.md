
# ssf2svg

This is a command line tool for converting SSF to SVG.

![Coelacanth SVG](https://github.com/mindboard/ssf2svg/blob/main/examples/coelacanth.svg)


## about SSF

SSF is Small Sketch native format.  
More details about Small Sketch :  
https://play.google.com/store/apps/details?id=com.mindboardapps.app.smallsketch


## build

Just do _go build_.  
Created __ssf2svg__ is executable.

This code was build and test go version go1.14.7 linux/amd64.


## usage

`cat examples/iphone-se.ssf | ./ssf2svg > iphone-se.svg`

Use an examples/style.json in order to customize output style.

`cat examples/iphone-se.ssf | ./ssf2svg examples/styles.json > iphone-se.svg`


There are some ssf and json examples in examples dir.



## License

See the LICENSE.txt file for license rights and limitations (MIT).

