// Package messaging will have a publisher and consumer.

// publisher should
// - resubscripe to queue in case of a failure
// - have a replay (all order in state OrderCreated) functionality
// - publish all incoming messages

// consumer should
// - resubscripe to queue in case of a failure
// - consume incoming messages
package messaging
