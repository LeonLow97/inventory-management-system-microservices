## Rate Limiting

- Rate limiting is a popular distributed system pattern.
- Rate limiting controls the rate at which users or services can access a resource, like an API, service, or a network.
- It plays a critical role in protecting system resources and ensuring fair use among all users, and also maintaining system stability.
- When the rate of requests exceeds the threshold defined by the rate limiter, the requests are throttled or blocked.
  - E.g., A user can send a message no more than 2 per second. One can create a maximum of 10 accounts per day from the same IP Address.

---

### Benefits of Rate Limiting

1. **Prevent Resource Starvation:**

- Mitigates Denial of Service (DoS) attacks by restricting the number of requests allowed within a certain timeframe.
  - A Denial of Service (DoS) attack is a malicious attempt to disrupt or limit access to a network, server, or service, making it inaccessible to its intended users by overwhelming it with an excessive amount of traffic or requests.
- Large tech companies, such as Twitter and Google, utilize rate limiting to restrict the number of actions users can perform within specified intervals, preventing system overload.

2. **Reduce Cost**:

- Helps control and limit resource usage, reducing the potential for overuse and preventing excessive costs.
- Particularly beneficial for services that interact with paid third-party APIs, where limiting API calls is crucial to controlling expenses.

3. **Overload Prevention**:

- Essential in maintaining server health and performance by avoiding overloading caused by high request volumes.
- Besides countering malicious attacks, rate limiting assists in handling heavy usage scenarios that could otherwise strain server resources and degrade service quality.

---

### Applications of Rate Limiting

- User-Level Rate Limiting
  - Implemented in platforms like social media to prevent spam or abuse by limiting the number of posts, comments, or actions a **user** can perform within a specific timeframe, ensuring fair usage and deterring malicious activities or bots.
- Application-Level Rate Limiting
  - Used in scenarios like online ticketing platforms during high-demand events to limit the rate of requests or transactions, preventing the system from being overwhelmed. For instance, limiting the number of ticket purchases per minute during a concert ticket sale.
- API-Level Rate Limiting
  - Commonly employed in services offering APIs (like cloud storage) to restrict the number of API calls a user can make per unit of time. This practice ensures fair access to resources, protects against misuse, and maintains system stability.
- User Account Level Rate Limiting
  - Implement in software-as-a-service (SaaS) platforms with multiple service tiers. Each tier may have distinct usage limits. For instance, free-tier users might have lower rate limits compared to premium-tier users, encouraging upgrades while managing resource allocation. Like ChatGPT OpenAI.

---

### Core Concepts of Rate Limiting

- Limit
  - Defines the maximum allowable requests or action within a specified time interval. For instance, limiting a user to send 100 messages per hour.
- Window
  - Represents the duration during which the limit is applied. It can span different timeframes such as an hour, a day or a week. Longer windows can pose challenges like storage durability.
- Identifier
  - Uniquely identifies individual callers, distinguishing between different entities making requests. Common examples include User IDs or IP Addresses.

---

### Types of Rate Limiting Responses

- Blocking
  - Denies access to the resource for requests that exceed the limit.
  - Often expressed as an error message (E.g., **HTTP Status Code - 429 Too Many Requests**), informing users of exceeded limits.
- Throttling
  - Slows down or delays requests that surpass the limit.
  - For instance, a video streaming service might lower the quality of the stream for users who exceed their data cap.
- Shaping
  - Allows requests beyond the limit but assigns them lower priority.
  - Users exceeding their limits receive a lower priority in processing, ensuring users within limits get better service.
  - For example, in a Content Delivery Network (CDN), requests from users over the limit may be processed last compared to requests from compliant users.

---

### Rate Limiting Algorithms

- Fixed Window Counter
  - Description: Tracks the count of requests within fixed intervals.
  - Algorithm: Resets the counter at fixed intervals (e.g., every second).
  - Pros: Simple to implement
  - Cons: Susceptible to bursting at the start of a new interval if requests arrive in large numbers. Large number of requests arriving at the start of the interval leads to bursts that exceed the limit before the algorithm has a chance to react and enforce the rate limit.
- Sliding Window Log
  - Description: Uses a log to track timestamps of requests
  - Algorithm: Records timestamps of requests in a sliding time window. Counts requests within the window and compares against the time limit.
  - Pros: Tracks requests precisely within the defined window.
  - Cons: May require extensive memory usage for logging timestamps in high-volume scenarios.
- Sliding Window Counter
  - Description: Similar to the fixed window counter but dynamically slides the window.
  - Algorithm: Divides time into intervals and slides the window along with time. Maintains a counter for each interval.
  - Pros: Provides more flexibility than the fixed window, allowing a sliding time frame.
  - Cons: Can be complex to implement and may require additional computational resources.
- Token Bucket
  - Description: Uses tokens to control the rate of requests.
  - Algorithm: Initializes a bucket with tokens (representing requests). Requests can be processed as long as tokens are available. Tokens are replenished at a fixed rate.
  - Pros: Allows bursts up to the bucket size, providing flexibility.
  - Cons: Implementation complexity, especially in scenarios with varying request rates.
- Leaky Bucket
  - Description: Controls the output of data at a constant rate.
  - Algorithm: Maintains a bucket with a leak rate. Requests fill the bucket, and excess requests are discarded or delayed.
  - Pros: Smooths the rate of requests over time.
  - Cons: Complexity in managing the leak rate and handling bursts.

---
