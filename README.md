# Go Lighthouse Tool

> A lighthouse score tracking tool in Go, SQLite and React/Svelte.

- Go REST API and cron for fetching website performance scores with the Google Lighthouse [CLI tool](https://github.com/GoogleChrome/lighthouse#using-the-node-cli 'Lighthouse CLI tool docs on GitHub'). Data saved locally in the app with SQLite.
- React & Svelte widgets for displaying the results.

**Very much in progress..**

## Endpoints

### `GET /`

Check if the API is up

### `GET /website?url=https://github.com/`

Fetch the latest report for a specific website

Params:

```js
url String (required) (url param)
```

### `GET /websites`

Fetch the latest report for all registered websites

Params:

```js
url String (required) (url param)
```

### `POST /website`

Trigger a refetch for a specific website

Params:

```js
url String (required) (url param)
```

### `POST /websites`

Trigger a refetch for all registered websites

Params:

```js
n / a
```

### `GET /view/website?url=https://github.com/`

View the latest report for a specific website

Params:

```js
url String (required) (url param)
```

### `GET /view/websites`

View the latest reports for all registered websites

Params:

```js
TBC
```
