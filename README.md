# ssf2svg

This is a command line tool for converting SSF to SVG.

![Coelacanth SVG](https://github.com/mindboard/ssf2svg/blob/main/examples/coelacanth.svg)


## About SSF

SSF is a Small Sketch native format.  
About Small Sketch, please see this page: https://play.google.com/store/apps/details?id=com.mindboardapps.app.smallsketch


## Build

In order to create a ssf2svg command:

```
go build
```

## Usage

```
cat examples/iphone-se.ssf | ./ssf2svg > iphone-se.svg
```

Use an examples/style.json in order to customize output style.

```
cat examples/iphone-se.ssf | ./ssf2svg examples/style.json > iphone-se.svg
```

There are some examples ssf and json in examples dir.

