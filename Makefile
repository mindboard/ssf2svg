
srcgo:=main.go matrix.go svg.go model.go style.go

svg : ssf2svg
	cat examples/iphone-se.ssf | ./ssf2svg > 1.svg
	cat examples/iphone-se.ssf | ./ssf2svg examples/style.json > 2.svg
	cat examples/coelacanth.ssf | ./ssf2svg > 3.svg
	cat examples/coelacanth.ssf | ./ssf2svg examples/style.json > 4.svg

ssf2svg: $(srcgo)
	go build

clean:
	$(RM) *.svg
	$(RM) ssf2svg
