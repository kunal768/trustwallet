# TrustWallet Block Txn Parser Assignment

## How To Use

This Go program provides a simple API server for interacting with a TrustWallet block transaction parser.

### API Endpoints

**1. Get Latest Block**

```bash
curl http://localhost:8080/currentBlock
```


**2. Get Transactions for Address**

```bash
curl http://localhost:8080/transactions\?address\=0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5
```

**3. Subscribe an Address**

```bash
curl -X POST http://localhost:8080/subscribe\?address\=0xf0588C1d1BCa1caDC91dFf8788a1BA123Afe5Cb2
```
