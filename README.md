# envguard

> A lightweight utility to validate and audit `.env` files against a schema before deployment.

---

## Installation

```bash
go install github.com/yourname/envguard@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envguard.git && cd envguard && go build ./...
```

---

## Usage

Define a schema file (`.env.schema`) listing required keys and optional rules:

```
DATABASE_URL=required
PORT=required,numeric
DEBUG=optional
API_KEY=required
```

Then validate your `.env` file against it:

```bash
envguard --schema .env.schema --env .env
```

**Example output:**

```
✔ DATABASE_URL   present
✔ PORT           present, valid
✘ API_KEY        missing
✘ PORT           expected numeric value

2 error(s) found. Deployment blocked.
```

Use the `--strict` flag to fail on any undeclared keys found in the `.env` file:

```bash
envguard --schema .env.schema --env .env --strict
```

---

## Flags

| Flag       | Description                              | Default      |
|------------|------------------------------------------|--------------|
| `--schema` | Path to the schema file                  | `.env.schema`|
| `--env`    | Path to the `.env` file to validate      | `.env`       |
| `--strict` | Fail on undeclared keys                  | `false`      |
| `--quiet`  | Suppress output, exit code only          | `false`      |

---

## License

[MIT](LICENSE)