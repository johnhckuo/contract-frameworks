compile:
	# solc github.com/OpenZeppelin/openzeppelin-contracts/=openzeppelin-contracts/ contracts/contract.sol
	solc --optimize --abi ./contracts/contract.sol -o build --overwrite
	solc --optimize --bin ./contracts/contract.sol -o build --overwrite
	abigen --abi=./build/Amis.abi --bin=./build/Amis.bin --pkg=api --out=./api/contract.go

node:
	docker run --detach --publish 8545:8545 trufflesuite/ganache-cli:latest --accounts 10 --debug

test:
	ginkgo -failFast -cover -v ./api