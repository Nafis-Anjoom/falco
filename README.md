# Falco
A fast chat app built on a concise, efficient, custom binary message transmission protocol over websocket

### Motivation
A chat application is a perfect experiment to learn about protocol design, data serialization, performance optimization, and networking. I wanted to explore the process
of designing a custom protocol over websocket that handles message sending/receiving, acknowledgements, and online status of users.

### Installation
1. Clone the repository and navigate to the project directory.
```bash
git clone https://github.com/Nafis-Anjoom/falco.git
cd falco
``` 
2. Install [Go](https://go.dev/doc/install). If already installed, make sure it's up-to-date.
3. Add environment variables for port, machine id, database url, and jwt secret.
```bash
export PORT=3000
export MACHINE_ID=0
export DB_URL=postgres://postgres:admin@localhost:1234/chat
export JWT_SECRET=sEcReT
```
4. Navigate to the server directory and start the backend.
```bash
go run ./api
```
6. Navigate to the ui directory and start the next app.
```bash
pnpm run dev
```
7. Open a browser, and go to hosting url for the next app.

### Reflections
This has been one of the most time-consuming, challenging projects I have undertaken. Initially, I wanted to compare the performance and productivity of JSON and custom binary protocol by creating a simple CLI chat application. Later, however, I decided to make it a functional centralized chat service that handles multiple live connections simultaneously. This system better reflected real-world chat systems, such as Facebook, WhatsApp, etc. The many hurdles throughout the development entailed considerable research and investigation into low-level topics. One of the notable examples is the nature of a WebSocket connection.\
Most of us are fairly familiar with HTTP; it is a half-duplex, transient, and stateless connection. For browsing the internet HTTP is perfect because clients (browsers) do not need a persistent connection to the server. However, for low-latency, real-time applications, HTTP falls short due to its transient property. Websockets are well-suited for real-time applications. It is fully directional, persistent, and stateful connections. It is much more flexible than HTTP. First, there are no defined HTTP methods and protocol overhead. Second, the protocol handshake only happens once during connection. Finally, it is more compact than HTTP. However, all of the benefits come at the cost of higher
complexity. Both sides of the connection are required to agree on the structure of the messages. Thus, for my application, I developed a custom binary protocol. Moreover, error handling and reconnection are much more difficult. In a chat application, for example, the developer has to decide what to do with messages sent when the server is down, how the server recovers state after restart, etc.\
These are the challenges that motivated me to pursue this project. I have some major plans for this project. First, I want to improve the error handling and message routing using a message queue. A message queue would handle push notifications and guaranteed message delivery. Second, I am in the process of implementing message seen/sent/delivered status. Finally, I want to experiment with horizontal scaling. This is perhaps the most challenging. I am still thinking about how messages would be stored, delivered, and retrieved when multiple servers are involved. Should I migrate to NoSQL for easier message retrieval and storage? How will I handle load balancing? How should I migrate the state to another server? These are the questions I will explore much more in-depth in the future.
