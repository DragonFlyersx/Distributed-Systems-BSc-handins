# Questions

## a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?
 The packages are integer numbers that we use. We use channels when sending and recieving from the client to server and server to client

## b) Does your implementation use threads or processes? Why is it not realistic to use threads?
The reason why it is not realistic to use threads is because having more clients could introduce race conditions if multiple clients attempting a handshake are not properly synchronized whilst using channels or mutex. 
It is therefore very difficult and complicated to set it up so that the server constantly can receive multiple requests and also send,
the correct information back to the correct clients. This makes it a very complicated to achieve the precise sequencing that 3-way handhsakes need. 
Because it is not a real network, and only a "simulated" network for learning purposes, we therefore chose to use threads in our implementation as we wanted to familiarize ourselves more with it. 

## c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
If the network changes the order of the messages, we could implement track sequence numbers. Meaning that each message includes a sequence number. 
We would then have to on the receiver end, check whether or not the sequence number is what we expected. If it was not what we expected, we could then buffer the out-of-order message and wait for the missing one to arrive. 

## d) In case messages can be delayed or lost, how does your implementation handle message loss?
The program does not directly handle the case if a message gets delayed or lost. Any message loss or delay
would in our program lead to a deadlock or incorrect behavior because the channels between client and server
assume reliable delivery and would just wait for a message.
This could be fixed by:
- Using timeouts on the channels to detect if a message have been delayed / lost.
- The client or server should resend the message or restart the handshake when a message has been lost.

## e) Why is the 3-way handshake important?
The 3-way handshake is very important both for security but also for the check of not loosing information
It helps create a connection to the server and if one package is missing the connection fails.
This will ensure that we have a secure and reliable connection to the server.
The handshake also ensures that both client and server are ready to send and recieve data.
Furthermore it makes the platform more secure since the sequence is random and is incremented from the server.
This ensures that a bad actor cant intercept.

Henceforth without the 3way handshake it would be less reliable, and the potential for data loss and duplication, or failure to detect connection issues
would be way higher. 
