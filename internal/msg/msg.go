package msg

const (
	Logo = `╔──────────────────────────────────────╗
│ ███████╗██╗██████╗ ██╗   ██╗ █████╗  │
│ ╚══███╔╝██║██╔══██╗██║   ██║██╔══██╗ │
│   ███╔╝ ██║██████╔╝██║   ██║███████║ │
│  ███╔╝  ██║██╔══██╗╚██╗ ██╔╝██╔══██║ │
│ ███████╗██║██║  ██║ ╚████╔╝ ██║  ██║ │
│ ╚══════╝╚═╝╚═╝  ╚═╝  ╚═══╝  ╚═╝  ╚═╝ │
╚──────────────────────────────────────╝
zirva server v1.0.0 (https://zirva.org)`
	PrivilegesErr = "you must run the server as root!"
	ServerRunning = "server is running on port %s"
	RegistrarErr  = "the server is not registered! register the server by entering this link in the portal:\nhttp://%s:%s/registrar?t=%s"
	RegistrarOk   = "registrar has been successfully!"
)
