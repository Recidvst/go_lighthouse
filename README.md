# Go Lighthouse Tool

> A lighthouse score tracking tool in Go, SQLite and React/Svelte.

- Go REST API and cron for fetching website performance scores with the Google
  Lighthouse [CLI tool](https://github.com/GoogleChrome/lighthouse#using-the-node-cli 'Lighthouse CLI tool docs on GitHub')
  . Data saved locally in the app with SQLite.
- React & Svelte widgets for displaying the results.

**Very much in progress..**

## TODO:

- Test suite (Go testing package)
- Add basic templates for viewing response directly
- Fix bug where date_fetched and date_edited in DB are in different timezones (GMT / BST)
- Fix bug where temporary CLI results file is not always removed.
- Decide how to handle description field in the site table.
- Code needs a bit of a tidy in terms of structure and organised comments

## Endpoints

### `GET /`

Check if the API is up

### `GET /website?url=https://github.com/`

Fetch the latest report saved in the DB for a specific website

Params:

```shell
url String (required) (url param)
```

### `GET /websites`

Fetch the latest report saved in the DB for all registered websites

Params:

```shell
url String (required) (url param)
```

### `POST /website`

Trigger a refetch of the report for a specific website

Params:

```shell
url String (required) (url param)
```

### `POST /websites`

Trigger a refetch of the report for all registered websites

Params:

```shell
n / a
```

### `GET /view/website?url=https://github.com/`

View the latest report for a specific website

Params:

```shell
url String (required) (url param)
```

### `GET /view/websites`

View the latest reports for all registered websites

Params:

```shell
TBC
```
