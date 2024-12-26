## Dashboard

Imagine having a dashboard for your personal needs, to playaround or just admire all the data you are generating at one place


### Development

#### Backend

I am using Pocketbase as a framework to extend any APIs required, and also because I liked the overall development of Pocketbase, so just to challenge myself, I have been using it. To run the project simply do:

```
go run cmd/dashboard/main.go serve
```

#### Frontend

For frontend, I am using Svelte without SSR as a static site so that it can be hosted over github, and using shadcn-svelte library for all the components required.

```
bun --bun run dev --host
```

This project uses Node v20.11.1, bun as JS Runtime

#### Application

Related application can be found at: https://github.com/shashank-sharma/metadata

Metadata application is responsible for generating data over machine and uses backend API from this project
