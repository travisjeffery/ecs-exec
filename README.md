# ecs-exec

Run commands on your ECS container instances.

## Installation

Download the binary for your platform from the [Releases](https://github.com/travisjeffery/ecs-exec/releases) page. Or use `go get`:

```
go get github.com/travisjeffery/ecs-exec
```

## Usage

```
ecs-exec --cluster=CLUSTER <cmd>...
```

Here's an example showing the containers running on your container instances:

```
ecs-exec --cluster=app docker ps
```

## Author

[Travis Jeffery](https://travisjeffery.com)

## License

MIT
