# TrustWallet Block Txn Parser Assignment

## How To Use

This Go program provides a simple API server for interacting with a TrustWallet block transaction parser.

### API Endpoints & Example Responses

**1. Get Latest Block**

```bash
curl http://localhost:8080/currentBlock
```

##### Response 
```bash
Current block: 19434993%
```

**2. Get Transactions for Address**

```bash
curl http://localhost:8080/transactions\?address\=0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5
```

##### Response 
```bash
{"address":"0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5","inbound":null,"outbound":[{"hash":"0x0bfb1e7deddc65ad079cd49c88e3502325c44d0ee87ccf0217b09a3ed2b226d4","from":"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5","to":"0x876528533158c07c1b87291c35f84104cd64ec01","value":"0x1b28bb3568448d2"}]}
```

**3. Subscribe an Address**

```bash
curl -X POST http://localhost:8080/subscribe\?address\=0xf0588C1d1BCa1caDC91dFf8788a1BA123Afe5Cb2
```

##### Response 
```bash
Subscribed to address: 0xf0588C1d1BCa1caDC91dFf8788a1BA123Afe5Cb2%
```
