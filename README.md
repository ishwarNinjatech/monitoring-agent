# Blockchain Monitoring Agent

## Setup Instructions

### Prerequisites

- Docker
- Docker Compose
- Go
- Anvil --- _Using `curl -L https://foundry.paradigm.xyz | bash`_

### Running the Monitoring Agent

Clone the repository:

git clone https://github.com/Ishwar-Parmar/monitoring-agent

cd monitoring-agent
Build and run the Go application:

go build -o monitor_metrics/main.go
./main

### Running the metrics locally

1. **Start local blockchain network:**
   run in a new terminal`anvil`

2. **Go to monitor_metrics dir in cli** and
   run `go run .`

### Running the metrics using docker

1. `sudo systemctl start docker`
2. `docker-compose build`
3. `docker-compose up`

## Approach Explanation

### Chosen Approach: Local Blockchain Setup

I chose a local blockchain setup using Anvil because it allows for controlled simulations of specific scenarios and attacks.
This environment provides the flexibility to test various conditions and monitor responses without the unpredictability of a live network.
It also enables detailed exploration into blockchain mechanics and targeted testing, which is crucial for developing a robust monitoring agent.
It create a local testnet node for deploying and testing smart contracts. It can also be used to fork other EVM compatible networks.

### Benefits

- **Controlled Environment:** Ability to simulate specific attack scenarios.
- **Flexibility:** Easier to manipulate and configure.
- **Isolation:** No dependency on external factors, ensuring consistent testing conditions.
- **EasierTesting:** It provides many accounts locally to test the netwrok reliably.

## METRICS

### Blockchain Network Metrics

**Get no of transactions in block (monitorTransactionThroughput)**
_Steps to Implement_

1.  Fetch recent block by number.
2.  Track the total number of transactions.

**How to be used to detect vulnerabilities**

1.  _Unusual Transaction Volumes_  
    **Vulnerability Detection:**
    _Spam Attacks:_ A sudden spike in the number of transactions can indicate a spam attack aimed at congesting the network.
    Detection Method: Set a threshold based on historical transaction volumes. Trigger an alert if the number of transactions in a block exceeds this threshold by a significant margin.
2.  _Dropped transaction volumes_
    A sudden drop in the number of transactions per block can be a sign of various issues, such as network congestion, denial-of-service attacks, or other forms of disruption.

**Monitor block time (monitorBlockTime)**
_Steps to Implement_

1. Fetch recent block number.
2. Fetch recent block by block number.
3. Parse the block to get timestamp.

**How to be used to detect vulnerabilities**

    1. _Consistency and Stability:_ Fluctuations in block time can indicate network congestion, inefficient mining processes, or potential issues with the blockchain network's infrastructure. Monitoring block time helps ensure that the network is stable and operating efficiently.

    2. _Fork Detection:_ Sudden changes in block time can be a sign of a network fork or a potential attack on the blockchain. By monitoring block time closely, you can detect any anomalies that may indicate the occurrence of a fork and take appropriate action to resolve the issue.

    3. _DDoS Attacks:_ Distributed Denial of Service (DDoS) attacks can impact a blockchain network's performance by causing delays in block creation. Monitoring block time can help in identifying abnormal delays that may be caused by DDoS attacks targeting the network.

    4. _Selfish Mining:_ Selfish mining is a strategy where miners withhold blocks to gain a competitive advantage. Monitoring block time can help in detecting instances of selfish mining by observing inconsistencies in block creation times and uncovering potential malicious behavior.

    5. _51% Attacks:_ A 51% attack occurs when an entity controls the majority of the network's mining power, allowing them to manipulate transactions. Monitoring block time can help in identifying signs of a 51% attack, such as unusually fast block creation times or irregularities in the blockchain's consensus mechanism.

**Monitor pool status (monitorTransactionPoolStatus)**
_Steps to Implement_

1. Fetch pool status using _txpool_status_ method.
2. Parse the pending and queued transaction from the pool.

**How to be used to detect vulnerabilities**

1.  Monitoring pending and queued transactions helps detect network congestion and transaction backlogs, enabling proactive optimization of network performance.

2.  Detection of transaction malleability attacks, spam attacks, and transaction reordering threats by monitoring pending and queued transactions, allowing for timely mitigation of security risks.

<!-- ----- DAPP Metrics ------ -->

**Get transaction count of contract address in each block (getTransactionCountOfContractInBlock)**
_Steps to Implement_

1. Fetch transaction count of contract address in each block using eth_getTransactionCount method passing contract address and the block number.
2. Parse the transaction from the result.

**How to be used to detect vulnerabilities**

1.  _Replay Attacks:_ An attacker could intercept and replay the eth_getTransactionCount request to trick the application into performing unintended actions based on the returned count.

2.  _Denial of Service (DoS) Attacks:_ An attacker could flood the application with a large number of requests to the eth_getTransactionCount method, causing resource exhaustion and disrupting the application's normal operation.

3.  _Privacy Concerns:_ The count of transactions from an address can reveal information about the address owner's activity, potentially compromising their privacy.

**Get transaction reciept of given transaction (getTransactionReciept)**
_Steps to Implement_

1. Fetch transaction reciept using _eth_getTransactionReceipt_ method passing transaction id.
2. Parse the transaction reciept from the result.

**How to be used to detect vulnerabilities**

1.  _Information Leakage:_ The transaction receipt contains sensitive information such as gas used, status, and logs. If this information is exposed to unauthorized parties, it could lead to privacy breaches or reveal details about the application's internal workings.

2.  _Denial of Service (DoS) Attacks:_ An attacker could flood the application with a large number of requests to the eth_getTransactionReceipt method, causing resource exhaustion and disrupting the application's normal operation.

3.  _Gas Price Manipulation:_ Attackers could analyze gas usage information in transaction receipts to manipulate gas prices or perform gas price oracle attacks.
