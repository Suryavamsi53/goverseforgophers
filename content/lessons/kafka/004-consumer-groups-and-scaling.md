# Consumer Groups and Scaling

If your `order_created` Topic is receiving 10,000 messages per second, a single Go Consumer application will be completely overwhelmed. The CPU will hit 100%, and the Consumer will fall further and further behind real-time (Consumer Lag).

You need to horizontally scale your Go application to 10 instances. But how do you ensure that all 10 instances don't process the *exact same messages*?

Kafka solves this with **Consumer Groups**.

## 1. The Consumer Group

When you configure your Go Kafka Consumer, you must assign it a `GroupID` (e.g., `group.id=billing-service`).

If you spin up 5 Go Pods in Kubernetes, and they all share the exact same `GroupID`, Kafka clusters them together into a single logical entity.

**The Golden Rule of Kafka Scaling:**
Kafka dynamically divides the Topic's Partitions among the active consumers in the Group. 
* 1 Partition can only be read by exactly 1 Consumer in the group at a time.

## 2. Scaling Scenarios

Imagine the `order_created` Topic has **10 Partitions**.

**Scenario A (1 Go Consumer):**
* You deploy 1 Go Pod (`billing-service`).
* That single Go Pod is assigned all 10 Partitions. It has to process 10,000 messages/sec alone.

**Scenario B (5 Go Consumers):**
* You deploy 5 Go Pods (`billing-service`).
* Kafka automatically triggers a **Rebalance**.
* Each Go Pod is assigned exactly 2 Partitions. 
* The load is perfectly distributed! Each Pod only processes 2,000 messages/sec.

**Scenario C (The Scaling Limit - 15 Go Consumers):**
* You deploy 15 Go Pods.
* Kafka assigns 1 Partition to each of the first 10 Pods.
* **The remaining 5 Pods are assigned 0 Partitions!** They sit completely idle, doing absolutely nothing!

*Enterprise Rule: You can never have more active Consumers in a Group than you have Partitions in a Topic. If you think you might need 50 Consumers in the future, you must configure the Topic to have 50 Partitions on Day 1!*

## 3. Multiple Microservices (Fan-Out)

What if the `Email Service` also needs to read the `order_created` events to send a receipt?

If the Email Service uses the *same* `GroupID` as the Billing Service, Kafka will distribute the partitions between them. The Billing Service will process half the orders, and the Email Service will process the other half! This is a catastrophe.

To implement the **Fan-Out Pattern** (Publish/Subscribe), the Email Service must use a completely different `GroupID` (e.g., `group.id=email-service`).

Because it is a different Group, Kafka treats it as a completely independent entity. The Billing Group gets 100% of the messages, and the Email Group *also* gets 100% of the exact same messages! 
Kafka maintains separate independent Offsets for every Consumer Group on the server.

## 4. Rebalancing Pauses

When a new Go Pod boots up, or an old Go Pod crashes, Kafka must redistribute the Partitions among the surviving Consumers. This is called a **Rebalance**.

During a Rebalance (which can take a few seconds), all Consumers stop processing data! If your Kubernetes cluster constantly scales Pods up and down every 30 seconds, your Kafka throughput will drop to zero due to endless Rebalancing pauses. Be careful with aggressive Autoscaling!
