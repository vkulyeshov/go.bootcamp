docker compose down
docker rm -f $(docker ps -aq)
docker rmi $(docker images -aq)
sudo rm -R pgdata ollama
docker system prune -f
docker system df