# ecs-exec

## Installation

```
go get github.com/travisjeffery/ecs-exec
```

## Usage

```
ecs-exec --cluster=CLUSTER<cmd>...
```

Here's an example showing the containers running on your container instances:


```
ecs-exec --cluster app docker ps
```

## License

MIT
