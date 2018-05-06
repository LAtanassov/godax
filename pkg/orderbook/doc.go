// Package orderbook represents a cryptocurrency exchange orderbook.
// trade life cycle - http://www.allaboutfinancecareers.co.uk/industry/infrastructure/the-trade-life-cycle-explained
//
// event sourcing - https://martinfowler.com/eaaDev/EventSourcing.html
// synonyms from text: aggregation = application
// - every change to the state of an aggregate is captured in an event
// - events are stored in the sequence they were applied
// - two different things are persisted an aggregation state and an event log
// - all changes to the aggregate are initiated by the event objects
// - complete rebuild, temporal query, event replay are possible
// - take regular snapshots
package orderbook
