### Description

This is a simple metrics agent that runs a web server and exposes the metrics for byteload.app health monitoring.

This agent uses [systeminformation](https://github.com/sebhildebrandt/systeminformation) package to collect the metrics.

### Requirements

This agent requires Node.js to be installed on the system.

Minimum supported version: 12.x

### Installation

To install the agent, first clone the repository:

```bash
git clone git@github.com:byteload/metrics-agent.git
```

Then, install the dependencies:

```bash
cd metrics-agent && npm install
```

### Configuration

The agent can be configured using the following environment variables:

- `PORT` - The port on which the web server will run. Default: `3000`
- SERVICES - A comma-separated list of services to monitor. Default: `nginx,mysql`

You can set these environment variables in a `.env` file in the root of the project or pass them directly when starting the agent.

If you want to use .env file, copy the `.env.example` file to `.env` and set the values accordingly.

### Usage

To start the agent, run:

```bash
npm run start
```

You can also use the `PORT` and `SERVICES` environment variables to override the default values:

```bash
PORT=4000 SERVICES=nginx,mysql npm run start
```

Alternatively, run the agent using the `node` command:

```bash
PORT=4000 SERVICES=nginx,mysql node src/index.js
```

After starting the agent, you can access the metrics at `http://localhost:3000`.

### Production Deployment

To deploy the agent in production, you can use a process manager like `pm2`:

```bash
npm install pm2@latest -g
pm2 start src/index.js --name metrics-agent
```

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.



