# MVP
+ one sonar (can see)
+ one navigator (can move)
+ 2d map
+ network. one server, two clients

## sample 
```
client req: /map
server resp: ### 
             .S#
			 #..

client req: /move/south

client req: /map
server resp: ..#
			 #S.
			 ###


User:
Client 1:
###
.S#
#..

..#
#S.
###

client 2:
move WNES > S
move WNES >
```

Next step:
Make it into the web.

Display of basic radar/sonar/map
Easy navigation



- MVP console
- MVP web
- WEB real time
  - sonar recieves updates every seconds
  - ship moves constantly
