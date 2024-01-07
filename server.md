## Microservices Architecture

### Challenges Faced

---

#### Inter-service communication between microservices

- Context:
  - I have two microservices - Inventory and Order. The Inventory Microservice manages the quantity of specific products, while the Order Microservice handles customer orders. When orders are placed, a message is sent via Apache Kafka Message Queue to the Inventory Microservice to decrease the quantity of the ordered product. Before sending the message to reduce the quantity, I need to retrieve the current product quantity.
- Initial thought process
  1. Retrieve the inventory quantity through gRPC communication. In the Inventory Microservice, implement a database transaction to fetch the inventory quantity and send this information back to the Order Microservice.
  1. In the Order Microservice, if the ordered quantity is less than the inventory quantity, we send an event via Apache Kafka to the Inventory Microservice to update the inventory count.
  - Will not work because there is a time period between checking the current inventory and actually decreasing the inventory
- Proposed Solution:
  1. In the Order Microservice, we send an event via Apache Kafka to the Inventory Microservice to update the inventory count. Additionally, we tag this order with a UUID.
  1. In the Inventory Microservice, we will recheck the product count in a database transaction. If the quantity is sufficient, we will proceed and eventually produce a message back to the Order Microservice to update the status of the order from 'Submitted' to 'Created'.
  1. If the quantity is insufficient, we will produce a message back to the Order Microservice to update the status of the order from 'Submitted' to 'Failed to Process (with reason)'.

---
