# Add More README Later

Note: build ui first then build go

```bash
cd ui
pnpm build # or npm build

cd ..
# execute immediately: go run . --env .env

go build -o app . # cross build: set GOOS=linux&& set GOARCH=amd64&& go build -o app .
# exec in ubuntu/debian
./app --env .env

```
