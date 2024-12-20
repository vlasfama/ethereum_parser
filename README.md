# Ethereum Transaction Parser

This project is an Ethereum blockchain parser that allows users to subscribe to Ethereum addresses, query transactions, and receive notifications for incoming/outgoing transactions.

---

## Features

- Subscribe to Ethereum addresses.
- Query transactions for subscribed addresses.
- Notify users about new transactions.

---

## Setup

1. Clone the repository:

   ```bash
   git clone git@github.com:vlasfama/ethereum_parser.git
   cd ethereum-parser
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the project:

   ```bash
   make build
   ```

---

## Commands

### Start the Server

Run the server to monitor blockchain transactions:

```bash
./eth-tx-parser start --rpc-url="https://ethereum-sepolia-rpc.publicnode.com" --port=8080 --log-level=info
```

### Generate a New Key Pair

Create a new Ethereum key pair:

```bash
./eth-tx-parser create_key
```

### Send a Transaction

Send Ethereum transactions:

```bash
./eth-tx-parser send --private-key="YOUR_PRIVATE_KEY" --to-address="0xRecipientAddress" --value="1000000000000000000" --rpc-url="https://ethereum-sepolia-rpc.publicnode.com"
```

---

## API Endpoints

### Subscribe to an Address

- **POST** `/subscribe`

  ```json
  {
      "address": "0xYourEthereumAddress"
  }
  ```

  Response:

  ```json
  {
      "success": true
  }
  ```

### Query Transactions

- **GET** `/transactions?address=0xYourEthereumAddress`
  Response:

  ```json
  [
      {
          "hash": "0xTransactionHash",
          "from": "0xSenderAddress",
          "to": "0xYourEthereumAddress",
          "value": "1000000000000000000",
          "blockNumber": 1234567,
          "timestamp": 1693450000
      }
  ]
  ```

### Get Current Block

- **GET** `/current-block`
  Response:

  ```json
  {
      "block": 1234567
  }
  ```

---

## License

This project is licensed under the MIT License.
