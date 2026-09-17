# toDoApp Backend

This is the backend for this frontend project, visit [here](https://github.com/yogisyo16/toDoApp).

In this project i use :
- [Docker](https://www.docker.com/) = for containerization
- [Golang](https://golang.org/) = as a primaly the backend code.
- [chi](https://github.com/go-chi/chi) = a lightweight, idiomatic and composable router for building Go HTTP services.
- [MongoDB](https://www.mongodb.com/) = for data storage and retrieval.
- [Mongoose](https://mongoosejs.com/) = for MongoDB object modeling.

Purpose of this project is to create a to-do list application, which allows users to add, edit, and delete tasks, as well as mark them as completed or incomplete, and also can send the reminder into their email later on. As in future user also can collaborate with other users to share and manage tasks together. and also provide a way to track progress and completion of tasks.

## Behind the scenes
Why making a to-do list application? Simple nowadays a lot of web apps like Notion, Trello, and Asana are popular tools for managing tasks and projects. However, the popular option is also having more pricing option, which can be a barrier for some users. So here i want to create a simpler with also user friendly interface. and also provide the authentication and authorization. So that user can manage their tasks easily and securely.

## Hope and meaning 
Hopefully with this project it can grow and become more like notion, but more reliable and user friendly.

## Features
Available features:
- CRUD operations for tasks, including adding, editing, and deleting tasks. (✔️)
- Marking tasks as completed or incomplete. (✔️)
- Adding details to tasks. (✔️)

In-Progress feature:
- User authentication and authorization. (✖️)
- User profile management. (✖️)
- User collaboration and sharing. (✖️)
- User notification system. (✖️)
- User email verification. (✖️)
- User password reset. (✖️)

## Documentation
- [`docs/API.md`](docs/API.md) — full endpoint reference, request/response shapes, example `curl` calls
- [`docs/SERVICES.md`](docs/SERVICES.md) — the data-access layer, one section per service
- [`docs/HANDLERS.md`](docs/HANDLERS.md) — the HTTP layer, one section per handler
- [`docs/RUNNING.md`](docs/RUNNING.md) — local setup (macOS/Linux/Windows), with Docker for MongoDB
- [`docs/BACKLOG.md`](docs/BACKLOG.md) — the JWT auth implementation plan and other known issues found while documenting


## AI Usage this far
- For adding documentation (Learning and reminder inline code)
- For any recomendation in the code (Not yet take full control code base)
   - Why still not taking the full control? This code is for learning, the next step AI here only to help me understand and make some reminder on what i code, since i made this project for a long term learning projects. So AI for now not the center of this code.
   - Do you use AI? Of course for my job projects, since it need us to do fast and accurate, so AI mainly on my job project. But learning? Not yet, still need to learn the basic and base.
