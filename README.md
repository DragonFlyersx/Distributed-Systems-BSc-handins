# Questions

## a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?
 The packages are integer numbers that we use, int the channels when sending and recieving from the client to server and server to client

## b) Does your implementation use threads or processes? Why is it not realistic to use threads?
Our implementation uses threads as we wanted to familiarize ourselves more with it. 
The reason why it is not realistic to use threads is because it is not realistic in terms of having more clients, as race conditions could occur if multiple clients. 
It is very difficult and complicated to set it up so that the server constantly can receive multiple requests and also send,
the correct information back to the correct clients. Which makes it a very complicated and probably more stress inducing process than it needs to be,
with other types of implementation. And at the same time it is not a real network

## c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?


## d) In case messages can be delayed or lost, how does your implementation handle message loss?
The program does not handle the case if a message gets delayed or lost. Any message loss or delay
would lead to a deadlock 

## e) Why is the 3-way handshake important?
The 3-way handskake is very important both for security but also for the check of not loosing information
It helps create a connection to the server and if one package is missing the connection fails.
This will ensure that we have a secure and reliable connection to the server.
 Furthermore it makes the platform more secure since the sequence is random and is incremented from the server.
 This ensures that a bad actor cant intercept.
