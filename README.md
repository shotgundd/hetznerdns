# Hetzner DNS CLI

A command-line tool to create and modify DNS records on Hetzner DNS service.

## Features

- Configure API token for authentication
- List DNS zones
- List DNS records for a zone
- Create new DNS records (A, AAAA, CNAME, MX, TXT, etc.)
- Update existing DNS records
- Delete DNS records
- Reference zones by name or ID

## Installation

### Using Go Install (Recommended)

If you have Go installed (version 1.16+), you can install directly using:

```
go install github.com/shotgundd/hetznerdns/cmd/hetznerdns@latest
```

This will install the latest version of the CLI tool to your `$GOPATH/bin` directory, which should be in your PATH.

### From Source

1. Clone the repository:
   ```
   git clone https://github.com/shotgundd/hetznerdns.git
   cd hetznerdns
   ```

2. Install using Make:
   ```
   make install
   ```

   This will use `go install` to install the binary to your `$GOPATH/bin` directory.

3. Verify the installation:
   ```
   hetznerdns version
   ```

### Pre-built Binaries

Download the pre-built binary for your platform from the [Releases](https://github.com/shotgundd/hetznerdns/releases) page and place it somewhere in your PATH.

## Usage

### Configuration

Before using the CLI, you need to configure your Hetzner DNS API token:

```
hetznerdns config set
```

You will be prompted to enter your API token. You can get this token from the Hetzner DNS Console.

You can also set the API token directly:

```
hetznerdns config set api-token YOUR_API_TOKEN
```

To view your current configuration:

```
hetznerdns config show
```

### Managing DNS Zones

List all your DNS zones:

```
hetznerdns zone list
```

### Managing DNS Records

List all records for a zone (you can use zone name or ID):

```
hetznerdns record list --zone example.com
```

Create a new record:

```
hetznerdns record create --zone example.com --name www --type A --value 192.168.1.1 --ttl 3600
```

Update an existing record:

```
hetznerdns record update --id RECORD_ID --zone example.com --value 192.168.1.2
```

Delete a record:

```
hetznerdns record delete --id RECORD_ID
```

## Examples

### Create an A record

```
hetznerdns record create --zone example.com --name www --type A --value 203.0.113.10
```

### Create a CNAME record

```
hetznerdns record create --zone example.com --name blog --type CNAME --value example.com
```

### Create an MX record

```
hetznerdns record create --zone example.com --name @ --type MX --value "10 mail.example.com"
```

## Development

### Building from Source

To build the project:

```
make build
```

This will create the binary in the `build` directory.

### Running Tests

To run all tests:

```
make test
```

To run only unit tests:

```
make test-unit
```

To run only integration tests:

```
make test-integration
```

To generate test coverage:

```
make test-coverage
```

### Cross-Compilation

To build for multiple platforms:

```
make build-all
```

This will create binaries for Linux, Windows, and macOS in the `build` directory.

## Continuous Integration

This project uses GitHub Actions for continuous integration:

- **Tests**: All tests are run on every push to the main branch and on pull requests.
- **Releases**: When a new tag is pushed (e.g., `v0.1.0`), a release is automatically created with binaries for multiple platforms.

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request