# minialertAisprid
Minialert is a minimalistic chalenge to send messages and receive alerts.

To build, launch "make" in the repository. If errors occures, it may be causes by dependancies. So install them, and relaunch.

To test launch in a terminal launch "./minialertAisprid -p 8080"
And in the other something like :
"
	$ ./minialertAisprid send battery 3
 
	$ ./minialertAisprid send battery 15
 
	$ ./minialertAisprid send battery 45
 
	$ ./minialertAisprid send battery 90
 
	$ ./minialertAisprid send battery 99
 
	$ ./minialertAisprid send battery 100
 
	$ ./minialertAisprid get alerts
"
