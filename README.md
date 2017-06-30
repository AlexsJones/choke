# choke
This is a proof of concept to show:

- RESTful API for MongoDB
- Limiting mongo connections to this server side process
- Sharing mongo session across request loading
- Rate limiting request load using channels

# Run the example

```
godep restore

go run main.go &

./create_tea.sh

```
