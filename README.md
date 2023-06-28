# Remove Node Modules
This script is designed to help you remove unused node_modules folders from old projects. It saves disk space by deleting the node_modules folders that haven't been modified within a specified number of days.

## Usage
To use the script, follow the instructions below:

1. Ensure you have Golang installed on your system.

2. Clone or download the repository containing the script to your local machine.

3. Open a terminal or command prompt and navigate to the directory where you downloaded the script.

4. Execute the script using the following command:

```bash
go run cmd/main.go -path <path> -days <days>
```
Replace <path> with the root directory where you want to search for and remove node_modules folders. <days> should be the threshold in days. For example, if you specify 30 as the threshold, any project folder that hasn't been modified within the last 30 days will have its node_modules folder deleted.

#### Note: The script will search for and delete node_modules folders recursively within the specified path.

## Example
Suppose you want to remove unused node_modules folders from the /Users/username/projects directory that haven't been modified within the last 60 days. You would execute the following command:

```bash
go run cmd/main.go -path /ujjujjuj/projects -days 60
```

Make sure to adjust the <path> and <days> values according to your requirements.
