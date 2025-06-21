# 🚦 Go Concurrency Problem: Traffic Light Controller

You're implementing a **concurrent traffic light controller** for a 4-way intersection. Each direction (North, East, South, West) gets a green light one at a time, in a rotating manner. Each light stays green for **2 seconds**, then turns red. The others must remain red during this time.

## 🧩 Your Task

1. Rotate the green light between all 4 directions continuously.
2. Print to stdout:

    ```
    [Time: 00:00:02] North is GREEN, East/South/West are RED
    [Time: 00:00:04] East is GREEN, North/South/West are RED
    ...
    ```

3. Run this for **10 seconds**, then stop cleanly.

## 📋 Constraints

- Each light should be represented as a goroutine.
- Use channels to coordinate which light is green.
- You must **gracefully stop all goroutines** after the 10 seconds are up (no leaking goroutines).

## 🧠 What This Tests

- Goroutines and scheduling
- Channels (including coordination and communication)
- Clean shutdown using `context.Context` or signal channels
- Time management with `time.Ticker`, `time.After`, etc.
- Structuring concurrent programs cleanly

---

**Bonus Challenge:** Add a "pedestrian crossing" signal that randomly interrupts and pauses traffic for 1 second — but make sure the system still shuts down cleanly at the 10-second mark.
