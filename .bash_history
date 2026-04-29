sudo apt update && sudo apt upgrade -y
sudo apt update && sudo apt upgrade -y
sudo apt install docker.io docker-compose-v2 -y
cd ~/cloud-ng-detector/detector
nano config.yaml
ls -a
nano config.yaml
ls -a
sudo go run main.go baseline.go blocker.go notifier.go unbanner.go
cd ~/cloud-ng-detector/detector
sudo iptables -D DOCKER-USER -s 196.216.144.2 -j DROP
nano dashboard.go
nano main.go
sudo go run main.go baseline.go blocker.go notifier.go unbanner.go dashboard.go
sudo iptables -D DOCKER-USER -s 196.216.144.2 -j DROP
sudo go run main.go baseline.go blocker.go notifier.go unbanner.go dashboard.go
sudo apt install docker.io docker-compose-v2 -y
sudo systemctl enable --now docker
mkdir -p ~/cloud-ng-detector/nginx
cd ~/cloud-ng-detector
nano nginx/nginx.conf
nano docker-compose.yml
sudo docker compose up -d
sudo docker ps
curl -I http://localhost
sudo iptables -L INPUT -n --line-numbers
sudo iptables -D INPUT -s 196.216.144.2 -p tcp --dport 80 -j DROP
nano blocker.go
ls -a
cd cloud-ng-dectector
cd cloud-ng-detector
cd detector/
cat audit.log
sudo go run main.go baseline.go blocker.go notifier.go unbanner.go dashboard.go
sudo apt update
sudo apt install golang-go -y
cd ~/cloud-ng-detector/detector
mkdir -p ~/cloud-ng-detector/detector
cd ~/cloud-ng-detector/detector
go mod init hng-detector
nano main.go
sudo go run main.go
nano ~/cloud-ng-detector/nginx/nginx.conf
cd ~/cloud-ng-detector
sudo docker compose restart nginx
cd ~/cloud-ng-detector/detector
sudo go run main.go
nano baseline.go
nano main.go
sudo go run main.go baseline.go
nano baseline.go
sudo go run main.go baseline.go
nano blocker.go
nano main.go
sudo go run main.go baseline.go blocker.go
nano blocker.go
sudo go run main.go baseline.go blocker.go notifier.go
nano notifier.go
sudo go run main.go baseline.go blocker.go notifier.go
nano blocker.go
nano notifier.go
nano unbanner.go
nano main.go
sudo go run main.go baseline.go blocker.go notifier.go unbanner.go
sudo iptables -D DOCKER-USER -s 196.216.144.2 -j DROP
sudo go run main.go baseline.go blocker.go notifier.go unbanner.go
sudo iptables -D DOCKER-USER -s 196.216.144.2 -j DROP
nano config.yaml
