# Distributed Mutual Exclusion in a Token Ring implementation
This is a project for our handin 4 in Distributed Systems.
The application implements the Token Ring algorithm to determine which client has access to a critical section.

## How to run
1. Open the `Node.go` file and specify the ip-address to the node you want to link up to.
a. (make sure that the person connecting to you have your ip-address)
2. Run the `Node.go` file by typing `go run .\Node\Node.go`
3. To initialize the token getting passed around in the ring, you will have to type `A`.
a. Now the token is being passed around in the circle of clients.
4. To access the critical section, you can type `A`, and it will begin to pass around your token to compare and finally give you access.

## Authors
- [@DragonFlyersx](https://github.com/DragonFlyersx)
- [@Nikolaj787](https://github.com/Nikolaj787)
- [@RasmusAChr](https://github.com/RasmusAChr)
- [@niko391a](https://github.com/niko391a)
![Logo](https://github.com/user-attachments/assets/b6a7a068-c6be-4760-ac28-969adc7b6371)

