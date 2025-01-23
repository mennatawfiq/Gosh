# Gosh - A Simple Shell in Go

**Gosh** is a lightweight shell implementation written in Go. It supports basic shell commands, built-in utilities, and output redirection. This project is a learning exercise and a foundation for building more advanced shell features.

---

## Features

- **Built-in Commands**:
  - like `exit`, `echo`, and `pwd`.

- **Output Redirection**:
  - Supports `>`, `>>`, `1>`, `2>`, and `2>>` for redirecting `stdout` and `stderr` to files.

- **External Commands**:
  - Execute external programs available in the system's `PATH`.

---

## Upcoming Features

I’m currently working on adding more features to make this shell more powerful and user-friendly. Some of the features in progress include:

**Autocompletion**, **History**, and **Pipelines**.

---

## Getting Started

### Prerequisites

- Go 1.20 or higher.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/mennatawfiq/Gosh.git
   cd Gosh
   ```
2. Build and run the project:
   ```bash
   go run .
   ```
