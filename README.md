# Go_Stomp

The purpose of this application is to ensure that when the live server crashes, the client automatically connects to another server or reconnect again to the live server with STOMP. 
ActiveMQ Artemis server has been used for the connection. I have two client applications written in Go. The STOMP protocol is being utilized, allowing the clients to connect to the server and communicate with each other.
To solve this problem, I used error handling. Both the master and slave brokers handle the error checking for whether it's nil or not. This is important because if the connection is not usable, an error will be returned, and I can perform a new Dial operation in that scenario.


## Test Results

1. Normal Connection: When Master server is live and the Slave is backup or Slave is down (even if the slave is not functioning, what matters is the master being live). Sender client starts to send message and receiver client get message.
   <img width="435" alt="image" src="https://github.com/nagkim/Go_Stomp/assets/65765559/7807f4c9-6b12-4d80-9b8f-f1b654bdb3f9">
   <img width="435" alt="image" src="https://github.com/nagkim/Go_Stomp/assets/65765559/2c34f099-6ff7-4819-901e-eeb4e907599a">

2. Master server is down: automatically backup server will be the current live server. Clients can know the master’s connection is lost by checking connection error nil or not. After the connection is lost, they disconnect from the queue. Therefore, they cannot continue sending messages. To proceed, they attempt to connect to a slave server, and once they successfully connect, the messaging process continues uninterrupted. However, unlike the receiver, the sender continues attempting to connect to the master even if it establishes a connection with the slave. Nevertheless, this is not an issue as it does not disrupt the flow.
   <img width="431" alt="image" src="https://github.com/nagkim/Go_Stomp/assets/65765559/61d0693b-67ae-427c-b0a4-be4f5c0fbd22">
   <img width="438" alt="image" src="https://github.com/nagkim/Go_Stomp/assets/65765559/142ba836-0991-4af1-9aba-b0f5c6f43e50">

3. Both server are down: If both servers crash, it will keep attempting to connect until the master is live again. If the backup becomes live again before the master, it won't connect; it waits for the master server to become live first.
   <img width="438" alt="image" src="https://github.com/nagkim/Go_Stomp/assets/65765559/1e28d585-7ac4-4c08-ade2-31ab01731fbe">
   <img width="438" alt="image" src="https://github.com/nagkim/Go_Stomp/assets/65765559/76f774fc-e088-4d71-9647-9dc521be32a4">






