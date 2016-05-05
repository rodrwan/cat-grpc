# Categorizer with gRPC

### How to

To build proto
```sh
$ make proto
```
This will generate a new directory call `categoryapi` with the new `pb` inteface.

Probably you should need install `lucky`
```sh
$ go get github.com/rodrwan/lucky
```
After all of this, you may run:
```sh
$ sh run.sh
```
This will up a new server instance. To test the server you need to run the following command on another terminal.
```sh
$ go run clint/client.go -description=Spotify
```


Cheers

R.
