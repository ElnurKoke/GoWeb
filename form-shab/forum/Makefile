.SILENT:

run:
	go run ./cmd/main.go

docBuild:
	docker build -t forum .

docrun: docBuild
	docker run -it -p 8080:8080 forum

docDelete:
	docker rmi forum

docClear: docDelete
	docker system prune -a

1:
	docker image build -f Dockerfile -t forum .
2:
	docker run -p 8081:8080 forum
3:
	docker stop $$(docker ps -a -q)
4:
	docker rm $$(docker ps -a -q)
5:  
	docker rmi $$(docker images -q)