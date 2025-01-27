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
2. Install [Go][go]. If already installed, make sure it's up-to-date.
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
