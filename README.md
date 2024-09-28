# Zirva Backbone Server
This server package is specially developed for zirva.org. With this package you can add backbone on zirva.org and start earning points.

![goreport](https://goreportcard.com/badge/github.com/zirvaorg/server)
![license](https://badgen.net/github/license/zirvaorg/server)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/zirvaorg/server)](https://pkg.go.dev/github.com/zirvaorg/server)

## Features
- [x] ping
- [x] http
- [x] tcp/udp
- [x] traceroute

## Installation
You can install the server by running the following command.
```bash
sudo wget -q -O - https://zirva.org/install.sh | sudo bash
```

## Getting Started
1. Start the server with `zirva` command. You can specify a port with the `-p PORT` argument.
2. Add the register url in the add server section in the portal.

## Requirements

| Requirement | Details                           |
|-------------|-----------------------------------|
| OS          | Linux (`x86_64`, `i386`, `arm64`) |
| curl        | **Required**                      |
| crontab     | **Required**                      |
| systemctl   | Optional                          |

## Remove Installation
If you want to remove the server, you can run the following command.

```bash
sudo wget -q -O - https://zirva.org/uninstall.sh | sudo bash
```

## Contributing
We welcome contributions from the community! You can open PRs and issues to help us improve the project.

## License
This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.
