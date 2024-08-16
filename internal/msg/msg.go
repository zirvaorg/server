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

	ServerRunning       = "server is running on port %s"
	PrivilegesErr       = "you must run the server as root!"
	ShutdownServer      = "shutting down the server..."
	ServerForceShutdown = "server forced to shutdown: %v"

	RegistrarOk          = "registrar has been successfully!"
	RegistrarErr         = "the server is not registered! register the server by entering this link in the portal:\nhttp://%s:%s/registrar?t=%s"
	RegistrarEnterPortal = "do not click directly on the link. just enter it in the field on the portal."
)
