[![Zero Detection Logo](https://github.com/ZERODETECTION/H9_Stage0/blob/main/logo_sm.png)](https://www.zerodetection.net/)


# Go Screenshot and Command Execution Server

A Go program that captures screenshots at regular intervals and serves them via an HTTP server. Optionally, it can execute commands received via HTTP requests.

## Features

- **Screenshot Capture**: Periodically captures screenshots of the entire screen and saves them locally.
- **HTTP Server**: Serves the saved screenshots and optionally provides an endpoint for executing system commands.
- **Command Execution**: Run commands on the host system via HTTP requests (optional and controlled by a flag).

## Installation

1. **Clone the Repository**:

   ```sh
   git clone https://github.com/yourusername/your-repo.git
   cd your-repo
   ```

2. **Install Dependencies**:

   Ensure you have Go installed. Run the following commands to set up the project:

   ```sh
   go mod tidy
   ```

## Usage

### Build and Run

To build and run the program with default settings:

```sh
go run main.go
```

Compile to executeable:
```sh
set GOOS=windows
set GOARCH=amd64
go build -o screenshot_server.exe
```sh

### Command-Line Arguments

- **`-dir`**: Directory to save screenshots. Default is `./screenshots/`.
- **`-addr`**: Address and port for the HTTP server. Default is `:8080`.
- **`-interval`**: Interval between screenshots (e.g., `2s`, `5m`). Default is `2s`.
- **`-max`**: Maximum number of screenshots to keep before deletion. Default is `10`.
- **`-enable-command`**: Enable the command execution endpoint. If this flag is provided, the endpoint `/sendcommand` will be available. 

### Examples

1. **Start the Server with Command Execution Enabled**:

   ```sh
   go run main.go -dir ./screenshots/ -addr :8080 -interval 5s -max 20 -enable-command
   ```

2. **Start the Server with Command Execution Disabled**:

   ```sh
   go run main.go -dir ./screenshots/ -addr :8080 -interval 5s -max 20
   ```

### HTTP Endpoints

- **`/`**: Serves the directory containing screenshots.
- **`/sendcommand`**: Executes a command on the host system. Requires `command` and optionally `param` query parameters.

  Example:

  - Execute Calculator:
  
    ```
    http://localhost:8080/sendcommand?command=calc
    ```

  - Open Notepad with a File:
  
    ```
    http://localhost:8080/sendcommand?command=notepad&param=example.txt
    ```

## Security Considerations

- **Command Execution**: Enabling command execution can be a security risk. Ensure you validate and sanitize input, and only enable this feature in a controlled environment.
- **Permissions**: The program should be run with appropriate permissions to prevent unauthorized access or actions.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please follow the [contribution guidelines](CONTRIBUTING.md) to submit issues or pull requests.
