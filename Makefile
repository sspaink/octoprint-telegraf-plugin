build:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o bin/octoprint cmd/main.go

deploy:
	ssh pi "rm ~/plugins/octoprint"
	scp bin/octoprint pi:~/plugins/octoprint
	ssh pi "chmod +x ~/plugins/octoprint"

reload:
	ssh pi "sudo service telegraf reload"