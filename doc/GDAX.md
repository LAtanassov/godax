# GDAX API

seperated in 
* trading and 
* feed.

trading requires authentication to be able place a order and request other account information.

## Matching Engine 

operates a continous first-come, first serve order book.
orders are executed in price-time priority.

### Example
User A places a Buy order for 1 BTC at 100 USD. User B then wishes to sell 1 BTC at 80 USD. Because User A’s order was first to the trading engine, they will have price priority and the trade will occur at 100 USD.

## Order Lifecycle

Order states are received, done, open, canceled, filled

## Trading Fees

operates a maker-taker model.

### Example
There is an existing SELL order for 5 BTC at 100 USD on the order book. You enter a BUY order for 7 BTC at 100 USD. 5 BTC of your BUY order are immediately matched and you are charged the taker fee because you are taking liquidity from the order book. The remaining 2 BTC of your order are now sitting on the BID side of the order book. A SELL order for 2 BTC at 100 USD arrives and matches against your 2 BTC BUY order. In this case you provided liquidity and are not charged any fees.

## Market Data

https://api.gdax.com
