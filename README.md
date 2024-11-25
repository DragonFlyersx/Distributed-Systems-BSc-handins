# Distributed active replication of auction house implementation
This is a project for our handin 5 in Distributed Systems.
We had to implement a System that works as an auction house and uses replication
We have chosen to use active replication meaning we have a client, a frontend, and 3 servers/nodes running. This code can sustain a node failure or absence of a node.

## How to run
1. Open the `Server.go` file and specify the ip-address/ port of the server

2. Run the `Server.go` file by typing `go run .\AuctionServer\Server.go`

a. (You will need to do this 3 times for the 3 servers.)

3. Open the `frontend.go` and make sure the ip - address are correct with the servers.

4. Run the `Client.go` file by typing `go run .\AuctionClient\Client.go`

a. Now you can use the command `Bid` specified with an amount after, this will start the auction. You can also use the `Result` to see the status of the auction 

## Authors
- [@DragonFlyersx](https://github.com/DragonFlyersx)
- [@Nikolaj787](https://github.com/Nikolaj787)
- [@RasmusAChr](https://github.com/RasmusAChr)
- [@niko391a](https://github.com/niko391a)
![Logo](https://github.com/user-attachments/assets/b6a7a068-c6be-4760-ac28-969adc7b6371)
