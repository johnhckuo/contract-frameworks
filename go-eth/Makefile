compile:
	# solc github.com/OpenZeppelin/openzeppelin-contracts/=openzeppelin-contracts/ contracts/contract.sol
	solc --optimize --abi ./contracts/contract.sol --bin ./contracts/contract.sol -o build --overwrite
	abigen --abi=./build/Amis.abi --bin=./build/Amis.bin --pkg=test --out=./api/test/contract.go
	solc --optimize --abi ./contracts/contract.sol --bin ./contracts/contract.sol -o build --overwrite
	solc --abi ./contracts/IERC721.sol --bin ./contracts/IERC721.sol -o build --overwrite
	abigen --abi=./build/IERC721.abi --bin=./build/IERC721.bin --pkg=erc721 --out=./api/erc721/ierc721.go

node:
	docker run --detach --publish 8545:8545 trufflesuite/ganache-cli:latest --accounts 10 --debug

test:
	ginkgo -failFast -cover -v ./...