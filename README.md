# zoom-to-control4

Simple TCP server to connect zoom rooms to Control 4. Built for a friend.

## 🚀 Getting Started

There are two ways to run the application, either via the native binary or via docker

### 🔟 Binary

1. Download the binary of the [latest commit from github](https://github.com/mhornbacher/zoom-to-control4/releases/tag/latest) for your platform and unzip it.
2. Run the binary with the environment variables set like so
    ```sh
    PORT=$PORT TARGET=$TARGET ./zoom-to-control-4-latest-$PLATFORM-$ARCH
    ```

<!-- ### 🐳 Docker -->

## ⚙️ Configuration

| Name | Purpose | Example |
| ---- | ------- | ------- |
| PORT | the port for the server to run on | 999 |
| TARGET | the ip/port for the messages to be sent to | 192.168.17.214:12345 |


