# Redis Configuration

QKO uses a custom `redis.conf` file to set important Redis security + configuration options. This is especially important for Redis 6 and higher, which include security enhancements that require explicit configuration.

## TLDR;
- **Why:** A custom `redis.conf` file is used to enforce security settings (like password protection) on Redis 6+, which requires explicit configuration.
- **How:** The file should contain config directives (even if it’s just one line) and must end with a newline to ensure proper parsing.
- **What to Commit:** Only commit a template file (`redis.conf.example`) with dummy values. Keep your actual configuration (`redis.conf`) out of version control via `.gitignore`.

## 1. Why We Use a `redis.conf` File

- **Enhanced Security:**
  Redis 6 introduced new security features. By using a custom `redis.conf` file, we can configure settings such as the password requirement (`requirepass`) to protect our Redis instance.

- **Custom Settings:**
  The configuration file allows us to override the default settings and enforce best practices. This ensures that our Redis instance is locked down properly before exposing any endpoints.

## 2. How to Structure the `redis.conf` File

- **Single-Line Configuration:**
  In some cases, the file may only need to contain a single configuration directive (for example, setting the required password). For instance:

  ```
  requirepass your_secure_password_here
  ```

- **Importance of Ending with a Hard Return:**
  After the configuration line, it is extremely important to add a newline (a hard return). This means after writing the line, press "Enter" so that there’s an empty line at the end of the file.
  Without this final newline, Redis might fail to read the last line of the file properly.

### 3. Example and Dummy Values

- **Providing a Template:**
  To prevent accidentally committing sensitive information, do **not** include real passwords in the repository. Instead, create a file named `redis.conf.example` that includes dummy values. For example:

  ```
  # Example Redis configuration file
  requirepass example_password
  ```

Once that is complete, add `redis.conf` to your `.gitignore` so that your real configuration doesn’t get pushed to the public repo.

### 4. How It Works

- **Custom Config Injection:**
  The `docker-compose.yml` file mounts `./redis.conf` into the Redis container at `/usr/local/etc/redis/redis.conf` (using a read-only mount). This ensures that when Redis starts, it uses our custom configuration instead of the default settings.

- **Key Configuration Settings:**
  In `redis.conf`, you can adjust important parameters such as:
  - Memory limits and eviction policies.
  - Authentication (if needed, by setting a password with `requirepass`).
  - Persistence settings.

### 5. How to Modify

1. **Editing the File:**
   To update Redis settings, open `redis.conf` in a text editor and make your changes.

2. **Applying Changes:**
   After updating `redis.conf`, you must restart the Redis container. You can do this by running:
   ```bash
   docker compose down
   docker compose up --build -d
   ```
   This ensures that Redis picks up the changes in its configuration file.

3. **Best Practices:**
   - **Test Changes:** Always test any modifications in a development environment before deploying to production.
   - **Version Control:** Although `redis.conf` is part of your repository, avoid storing sensitive values (such as passwords) directly in it. Instead, consider using environment variables for sensitive data.

### Note

If you have multiple environments (development, test, production) and need different Redis configurations, consider maintaining separate configuration files like `redis.dev.conf` and `redis.prod.conf` in the repository, and adjust your `docker-compose.yml` accordingly (using the `--env-file` approach). Make sure to document any environment-specific differences clearly.
