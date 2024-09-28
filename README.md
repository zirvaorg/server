<div align="center"><img src="https://portal.zirva.org/resources/img/logo.svg" alt="zirva.org" width="150" /> <h1>Zirva Community Server</h1></div>

This server package is specially developed for [zirva.org](https://zirva.org). With this package, you can add a community server to Zirva and start earning points.

![goreport](https://goreportcard.com/badge/github.com/zirvaorg/server)
![license](https://badgen.net/github/license/zirvaorg/server)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/zirvaorg/server)](https://pkg.go.dev/github.com/zirvaorg/server)

## Getting Started
1. Start the server with `zirva` command. You can specify a port with the `-p PORT` argument.
2. Add the register url in the add server section in the [portal](https://portal.zirva.org).
3. Everything is ready!

## Installation
```bash
sudo wget -q -O - https://zirva.org/install.sh | sudo bash
```

## Remove Installation
```bash
sudo wget -q -O - https://zirva.org/uninstall.sh | sudo bash
```

## Requirements
| Requirement | Details                           |
|-------------|-----------------------------------|
| OS          | Linux (`x86_64`, `i386`, `arm64`) |
| curl        | **Required**                      |
| crontab     | **Required**                      |
| systemctl   | Optional                          |

## Contributing
We welcome contributions from the community! You can open PRs and issues to help us improve the project.

## License
This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.
