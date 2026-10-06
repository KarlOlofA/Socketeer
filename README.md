# __**Socketeer**__

Socketeer is a simple socket service made to facilitate really simple communications through a middle layer. 

You drop this on a server somewhere and you're now able to easily configure middleware and processing of events reaching it.

## Motivation

I built this to allow for simple communication for further development for me while not needing a bigger solution. Specifically I've used it to build a realtime collaborative drawing program. So essentially it will allow you to simply start networking and communicating without rebuilding the wheel or requiring a massive replacement.

## ⚙️ Install

` go get https://github.com/karlolofa/socketeer `

## 🚀 Quick Start 

Install the project and setup a quick server

```
func main() {
	godotenv.Load()

	args := os.Args

	settings := pool.TcpServerSettings{
		Host:               "0.0.0.0",
		Port:               "8080",
		Method:             "tcp",
		Key:                "1234",
		ConnectionPoolSize: 10,
	}
	if len(args) > 1 && len(args[1]) > 0 {
		settings.Host = args[1]
	}

	server := pool.NewTcpServer(settings)
	if server == nil {
		log.Fatal("Failed to construct tcp server.")
		return
	}
	defer server.Close()

	fmt.Printf("Listening to %s:%s.\n", settings.Host, settings.Port)
	server.Run()

}
  ```
## Usage

You can easily change the general settings such as host, port and tcp method with the following code block. You can also set a "auth key" that is required to be sent with each request to allow it to be processed. (This is not a security measure how ever. Just the simplest kind of block.)

```
	settings := pool.TcpServerSettings{
		Host:               "0.0.0.0",
		Port:               "8080",
		Method:             "tcp",
		Key:                "1234",
		ConnectionPoolSize: 10,
	}
```

## Contributing

For new features, please open a Discussion prior to creating a pull request. This allows the me and community to think about the change and provide feedback. When the maintainers are ready to accept new features, we will look through the "Ideas" in the project's GitHub Discussions.
