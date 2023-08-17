# Go-Sockets

## Description

-   Simple Go Programs to allow users to send messages between each other in a chat room
-   Uses only standard library; just `go run ...`

## Organization

-   `Client`: User-facing program to send messages
-   `Server`: Server to facilitate said messages; Uses Go channels to sync concurrent processes
-   `Server-Lock`: Server ...; Uses Read/Write Mutex locks to sync concurrent processes
